package mcp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	outputIdleTimeout   = 5 * time.Second
	startupReadTimeout  = 200 * time.Millisecond
	startupQuietTimeout = 20 * time.Millisecond
)

var (
	validRoles = map[string]struct{}{
		"plan":      {},
		"implement": {},
		"review":    {},
	}
	validAgents = map[string]struct{}{
		"claude": {},
		"codex":  {},
	}
)

type SessionStatus string

const (
	SessionStatusRunning SessionStatus = "running"
	SessionStatusExited  SessionStatus = "exited"
	SessionStatusStopped SessionStatus = "stopped"
)

type Session struct {
	ID     string        `json:"session_id"`
	Role   string        `json:"role"`
	Agent  string        `json:"agent"`
	Status SessionStatus `json:"status"`
	PID    int           `json:"pid,omitempty"`

	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser

	outputMu          sync.Mutex
	output            bytes.Buffer
	lastCommandOffset int
	updateCh          chan struct{}
	done              chan struct{}
	waitErr           error
	logger            *slog.Logger
}

type SessionManager struct {
	mu       sync.Mutex
	sessions map[string]*Session
	launch   launcherFunc
	logger   *slog.Logger
}

type launcherFunc func(ctx context.Context, role, agent string) (*exec.Cmd, error)

type SessionList struct {
	Sessions []SessionInfo `json:"sessions"`
}

type SessionInfo struct {
	SessionID string        `json:"session_id"`
	Role      string        `json:"role"`
	Agent     string        `json:"agent"`
	Status    SessionStatus `json:"status"`
	PID       int           `json:"pid,omitempty"`
}

type CommandResult struct {
	Role      string `json:"role"`
	Command   string `json:"command"`
	SessionID string `json:"session_id"`
}

type OutputResult struct {
	Role      string        `json:"role"`
	Output    string        `json:"output"`
	SessionID string        `json:"session_id"`
	Status    SessionStatus `json:"status"`
}

func NewSessionManager(logger *slog.Logger) *SessionManager {
	return newSessionManager(defaultLauncher, logger)
}

func newSessionManager(launch launcherFunc, logger *slog.Logger) *SessionManager {
	if logger == nil {
		logger = newDiscardLogger()
	}

	return &SessionManager{
		sessions: make(map[string]*Session),
		launch:   launch,
		logger:   logger,
	}
}

func defaultLauncher(ctx context.Context, role, agent string) (*exec.Cmd, error) {
	if err := validateRole(role); err != nil {
		return nil, err
	}
	if err := validateAgent(agent); err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("determine working directory: %w", err)
	}

	scriptPath := filepath.Join(cwd, "scripts", "ai-launch.sh")
	if _, err := os.Stat(scriptPath); err != nil {
		return nil, fmt.Errorf("locate launcher script: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	cmd := exec.Command(scriptPath, role, agent)
	cmd.Dir = cwd
	return cmd, nil
}

func (m *SessionManager) StartSession(ctx context.Context, role, agent string) (SessionInfo, error) {
	if err := validateRole(role); err != nil {
		m.logger.Error("start session failed", "role", role, "agent", agent, "error", err)
		return SessionInfo{}, err
	}
	if err := validateAgent(agent); err != nil {
		m.logger.Error("start session failed", "role", role, "agent", agent, "error", err)
		return SessionInfo{}, err
	}

	m.mu.Lock()
	if existing, ok := m.sessions[role]; ok {
		if existing.Status == SessionStatusRunning {
			m.mu.Unlock()
			err := fmt.Errorf("session for role %q is already running", role)
			m.logger.Error("start session failed", "role", role, "agent", agent, "error", err)
			return SessionInfo{}, err
		}
		delete(m.sessions, role)
	}
	m.mu.Unlock()

	cmd, err := m.launch(ctx, role, agent)
	if err != nil {
		m.logger.Error("start session failed", "role", role, "agent", agent, "error", err)
		return SessionInfo{}, err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		m.logger.Error("start session failed", "role", role, "agent", agent, "error", err)
		return SessionInfo{}, fmt.Errorf("open stdin pipe: %w", err)
	}

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		m.logger.Error("start session failed", "role", role, "agent", agent, "error", err)
		return SessionInfo{}, fmt.Errorf("create stdout pipe: %w", err)
	}
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stdoutWriter

	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdoutReader.Close()
		stdoutWriter.Close()
		m.logger.Error("start session failed", "role", role, "agent", agent, "error", err)
		return SessionInfo{}, fmt.Errorf("start session process: %w", err)
	}
	if err := stdoutWriter.Close(); err != nil {
		stdin.Close()
		stdoutReader.Close()
		m.logger.Error("start session failed", "role", role, "agent", agent, "error", err)
		return SessionInfo{}, fmt.Errorf("close stdout writer: %w", err)
	}

	session := &Session{
		ID:       fmt.Sprintf("%s-%s-%d", role, agent, time.Now().UTC().UnixNano()),
		Role:     role,
		Agent:    agent,
		Status:   SessionStatusRunning,
		PID:      cmd.Process.Pid,
		cmd:      cmd,
		stdin:    stdin,
		stdout:   stdoutReader,
		updateCh: make(chan struct{}, 1),
		done:     make(chan struct{}),
		logger:   m.logger,
	}

	m.mu.Lock()
	m.sessions[role] = session
	m.mu.Unlock()

	go session.captureOutput()
	go session.waitForExit()
	session.captureStartupOutput()

	m.logger.Info("session started", "role", role, "agent", agent, "pid", session.PID, "session_id", session.ID)

	return session.info(), nil
}

func (m *SessionManager) StopSession(role string) (SessionInfo, error) {
	m.mu.Lock()
	session, ok := m.sessions[role]
	if !ok {
		m.mu.Unlock()
		err := fmt.Errorf("no session for role %q", role)
		m.logger.Error("stop session failed", "role", role, "error", err)
		return SessionInfo{}, err
	}
	delete(m.sessions, role)
	m.mu.Unlock()

	if session.Status == SessionStatusRunning && session.cmd.Process != nil {
		m.logger.Info("sending session signal", "role", role, "session_id", session.ID, "signal", "SIGTERM")
		if err := session.cmd.Process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
			m.logger.Error("stop session failed", "role", role, "session_id", session.ID, "error", err)
			return SessionInfo{}, fmt.Errorf("stop session %q: %w", role, err)
		}
	}

	select {
	case <-session.done:
	case <-time.After(2 * time.Second):
	}

	session.outputMu.Lock()
	if session.Status == SessionStatusRunning {
		session.Status = SessionStatusStopped
	}
	session.outputMu.Unlock()

	info := session.info()
	m.logger.Info("session stopped", "role", role, "session_id", info.SessionID, "status", info.Status)
	return info, nil
}

func (m *SessionManager) SendCommand(ctx context.Context, role, command string) (CommandResult, error) {
	m.mu.Lock()
	session, ok := m.sessions[role]
	m.mu.Unlock()
	if !ok {
		err := fmt.Errorf("no session for role %q", role)
		m.logger.Error("send command failed", "role", role, "command", command, "error", err)
		return CommandResult{}, err
	}
	if session.Status != SessionStatusRunning {
		err := fmt.Errorf("session for role %q is not running", role)
		m.logger.Error("send command failed", "role", role, "command", command, "session_id", session.ID, "error", err)
		return CommandResult{}, err
	}

	m.logger.Info("sending command", "role", role, "command", command, "session_id", session.ID)
	if err := session.writeCommand(command); err != nil {
		m.logger.Error("send command failed", "role", role, "command", command, "session_id", session.ID, "error", err)
		return CommandResult{}, err
	}

	result := CommandResult{
		Role:      role,
		Command:   command,
		SessionID: session.ID,
	}
	m.logger.Info("command queued", "role", role, "command", command, "session_id", session.ID)
	return result, nil
}

func (m *SessionManager) GetOutput(ctx context.Context, role string, timeout time.Duration) (OutputResult, error) {
	m.mu.Lock()
	session, ok := m.sessions[role]
	m.mu.Unlock()
	if !ok {
		err := fmt.Errorf("no session for role %q", role)
		m.logger.Error("get output failed", "role", role, "error", err)
		return OutputResult{}, err
	}
	if session.Status != SessionStatusRunning && !session.hasBufferedOutput() {
		err := fmt.Errorf("session for role %q is not running", role)
		m.logger.Error("get output failed", "role", role, "session_id", session.ID, "error", err)
		return OutputResult{}, err
	}

	output, err := session.readOutput(ctx, timeout)
	if err != nil {
		m.logger.Error("get output failed", "role", role, "session_id", session.ID, "error", err)
		return OutputResult{}, err
	}

	info := session.info()
	result := OutputResult{
		Role:      role,
		Output:    output,
		SessionID: session.ID,
		Status:    info.Status,
	}
	m.logger.Info("output retrieved", "role", role, "session_id", session.ID, "status", result.Status, "output_bytes", len(output))
	return result, nil
}

func (m *SessionManager) ListSessions() SessionList {
	m.mu.Lock()
	defer m.mu.Unlock()

	sessions := make([]SessionInfo, 0, len(m.sessions))
	for _, session := range m.sessions {
		sessions = append(sessions, session.info())
	}
	return SessionList{Sessions: sessions}
}

func (s *Session) info() SessionInfo {
	s.outputMu.Lock()
	defer s.outputMu.Unlock()

	return SessionInfo{
		SessionID: s.ID,
		Role:      s.Role,
		Agent:     s.Agent,
		Status:    s.Status,
		PID:       s.PID,
	}
}

func (s *Session) captureOutput() {
	defer s.stdout.Close()

	buf := make([]byte, 4096)
	for {
		n, err := s.stdout.Read(buf)
		if n > 0 {
			s.outputMu.Lock()
			s.output.Write(buf[:n])
			s.outputMu.Unlock()
			s.logger.Debug("session output received", "role", s.Role, "session_id", s.ID, "bytes", n)

			select {
			case s.updateCh <- struct{}{}:
			default:
			}
		}
		if err != nil {
			return
		}
	}
}

func (s *Session) waitForExit() {
	err := s.cmd.Wait()

	s.outputMu.Lock()
	if s.waitErr == nil {
		s.waitErr = err
	}
	if s.Status == SessionStatusRunning {
		if err == nil {
			s.Status = SessionStatusStopped
		} else {
			s.Status = SessionStatusExited
		}
	}
	status := s.Status
	s.outputMu.Unlock()

	_ = s.stdin.Close()
	if err != nil {
		s.logger.Error("session exited", "role", s.Role, "session_id", s.ID, "status", status, "error", err)
	} else {
		s.logger.Info("session exited", "role", s.Role, "session_id", s.ID, "status", status)
	}
	close(s.done)
}

func (s *Session) captureStartupOutput() {
	timer := time.NewTimer(startupReadTimeout)
	defer timer.Stop()

	quietTimer := time.NewTimer(time.Hour)
	if !quietTimer.Stop() {
		<-quietTimer.C
	}
	defer quietTimer.Stop()

	for {
		select {
		case <-s.updateCh:
			if !quietTimer.Stop() {
				select {
				case <-quietTimer.C:
				default:
				}
			}
			quietTimer.Reset(startupQuietTimeout)
		case <-quietTimer.C:
			s.outputMu.Lock()
			s.lastCommandOffset = s.output.Len()
			s.outputMu.Unlock()
			return
		case <-timer.C:
			s.outputMu.Lock()
			s.lastCommandOffset = s.output.Len()
			s.outputMu.Unlock()
			return
		}
	}
}

func (s *Session) writeCommand(command string) error {
	command = strings.TrimSpace(command)
	if command == "" {
		return fmt.Errorf("command must not be empty")
	}

	s.outputMu.Lock()
	s.lastCommandOffset = s.output.Len()
	s.outputMu.Unlock()

	if _, err := io.WriteString(s.stdin, command+"\n"); err != nil {
		s.outputMu.Lock()
		s.Status = SessionStatusExited
		s.outputMu.Unlock()
		return fmt.Errorf("write command: %w", err)
	}
	return nil
}

func (s *Session) readOutput(ctx context.Context, timeout time.Duration) (string, error) {
	startOffset := s.commandOffset()

	quietTimer := time.NewTimer(time.Hour)
	if !quietTimer.Stop() {
		<-quietTimer.C
	}
	defer quietTimer.Stop()

	responseTimer := time.NewTimer(timeout)
	defer responseTimer.Stop()

	sawOutput := s.outputSince(startOffset) != ""
	if sawOutput {
		quietTimer.Reset(outputIdleTimeout)
	}

	for {
		select {
		case <-ctx.Done():
			output := s.outputSince(startOffset)
			if output != "" {
				return output, nil
			}
			return "", ctx.Err()
		case <-s.updateCh:
			if output := s.outputSince(startOffset); output != "" {
				sawOutput = true
				if !quietTimer.Stop() {
					select {
					case <-quietTimer.C:
					default:
					}
				}
				quietTimer.Reset(outputIdleTimeout)
			}
		case <-quietTimer.C:
			if sawOutput {
				return s.outputSince(startOffset), nil
			}
		case <-responseTimer.C:
			output := s.outputSince(startOffset)
			return output, nil
		case <-s.done:
			output := s.outputSince(startOffset)
			if output != "" {
				return output, nil
			}
			if s.waitErr != nil {
				return "", s.waitErr
			}
			return "", nil
		}
	}
}

func (s *Session) outputLen() int {
	s.outputMu.Lock()
	defer s.outputMu.Unlock()
	return s.output.Len()
}

func (s *Session) commandOffset() int {
	s.outputMu.Lock()
	defer s.outputMu.Unlock()
	return s.lastCommandOffset
}

func (s *Session) hasBufferedOutput() bool {
	s.outputMu.Lock()
	defer s.outputMu.Unlock()
	return s.lastCommandOffset < s.output.Len()
}

func (s *Session) outputSince(offset int) string {
	s.outputMu.Lock()
	defer s.outputMu.Unlock()

	data := s.output.Bytes()
	if offset >= len(data) {
		return ""
	}
	return string(data[offset:])
}

func validateRole(role string) error {
	if _, ok := validRoles[role]; !ok {
		return fmt.Errorf("unsupported role %q", role)
	}
	return nil
}

func validateAgent(agent string) error {
	if _, ok := validAgents[agent]; !ok {
		return fmt.Errorf("unsupported agent %q", agent)
	}
	return nil
}
