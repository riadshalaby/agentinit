package mcp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	outputIdleTimeout     = 150 * time.Millisecond
	outputResponseTimeout = 2 * time.Second
)

var (
	errSessionOutputTimeout = errors.New("timed out waiting for session output")
	validRoles              = map[string]struct{}{
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

	outputMu sync.Mutex
	output   bytes.Buffer
	updateCh chan struct{}
	done     chan struct{}
	waitErr  error
}

type SessionManager struct {
	mu       sync.Mutex
	sessions map[string]*Session
	launch   launcherFunc
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
	Output    string `json:"output"`
	SessionID string `json:"session_id"`
}

func NewSessionManager() *SessionManager {
	return newSessionManager(defaultLauncher)
}

func newSessionManager(launch launcherFunc) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		launch:   launch,
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

	cmd := exec.CommandContext(ctx, scriptPath, role, agent)
	cmd.Dir = cwd
	return cmd, nil
}

func (m *SessionManager) StartSession(_ context.Context, role, agent string) (SessionInfo, error) {
	if err := validateRole(role); err != nil {
		return SessionInfo{}, err
	}
	if err := validateAgent(agent); err != nil {
		return SessionInfo{}, err
	}

	m.mu.Lock()
	if existing, ok := m.sessions[role]; ok {
		if existing.Status == SessionStatusRunning {
			m.mu.Unlock()
			return SessionInfo{}, fmt.Errorf("session for role %q is already running", role)
		}
		delete(m.sessions, role)
	}
	m.mu.Unlock()

	cmd, err := m.launch(context.Background(), role, agent)
	if err != nil {
		return SessionInfo{}, err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return SessionInfo{}, fmt.Errorf("open stdin pipe: %w", err)
	}

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		return SessionInfo{}, fmt.Errorf("create stdout pipe: %w", err)
	}
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stdoutWriter

	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdoutReader.Close()
		stdoutWriter.Close()
		return SessionInfo{}, fmt.Errorf("start session process: %w", err)
	}
	if err := stdoutWriter.Close(); err != nil {
		stdin.Close()
		stdoutReader.Close()
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
	}

	m.mu.Lock()
	m.sessions[role] = session
	m.mu.Unlock()

	go session.captureOutput()
	go session.waitForExit()

	return session.info(), nil
}

func (m *SessionManager) StopSession(role string) (SessionInfo, error) {
	m.mu.Lock()
	session, ok := m.sessions[role]
	if !ok {
		m.mu.Unlock()
		return SessionInfo{}, fmt.Errorf("no session for role %q", role)
	}
	delete(m.sessions, role)
	m.mu.Unlock()

	if session.Status == SessionStatusRunning && session.cmd.Process != nil {
		if err := session.cmd.Process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return SessionInfo{}, fmt.Errorf("stop session %q: %w", role, err)
		}
	}

	select {
	case <-session.done:
	case <-time.After(outputResponseTimeout):
	}

	session.outputMu.Lock()
	if session.Status == SessionStatusRunning {
		session.Status = SessionStatusStopped
	}
	session.outputMu.Unlock()

	return session.info(), nil
}

func (m *SessionManager) SendCommand(ctx context.Context, role, command string) (CommandResult, error) {
	m.mu.Lock()
	session, ok := m.sessions[role]
	m.mu.Unlock()
	if !ok {
		return CommandResult{}, fmt.Errorf("no session for role %q", role)
	}
	if session.Status != SessionStatusRunning {
		return CommandResult{}, fmt.Errorf("session for role %q is not running", role)
	}

	output, err := session.send(ctx, command)
	if err != nil {
		return CommandResult{}, err
	}

	return CommandResult{
		Role:      role,
		Command:   command,
		Output:    output,
		SessionID: session.ID,
	}, nil
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
	s.outputMu.Unlock()

	_ = s.stdin.Close()
	close(s.done)
}

func (s *Session) send(ctx context.Context, command string) (string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return "", fmt.Errorf("command must not be empty")
	}

	startOffset := s.outputLen()
	if _, err := io.WriteString(s.stdin, command+"\n"); err != nil {
		return "", fmt.Errorf("write command: %w", err)
	}

	quietTimer := time.NewTimer(time.Hour)
	if !quietTimer.Stop() {
		<-quietTimer.C
	}
	defer quietTimer.Stop()

	responseTimer := time.NewTimer(outputResponseTimeout)
	defer responseTimer.Stop()

	sawOutput := false

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
			if output != "" {
				return output, nil
			}
			return "", errSessionOutputTimeout
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
