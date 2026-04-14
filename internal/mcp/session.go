package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	outputIdleTimeout   = 5 * time.Second
	startupReadTimeout  = 200 * time.Millisecond
	startupQuietTimeout = 20 * time.Millisecond
	stopTermGracePeriod = 2 * time.Second
	stopKillGracePeriod = 500 * time.Millisecond
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
	codexSessionIDPattern = regexp.MustCompile(`(?m)^session id:\s+([^\s]+)$`)
)

type SessionStatus string

const (
	SessionStatusRunning SessionStatus = "running"
	SessionStatusExited  SessionStatus = "exited"
	SessionStatusStopped SessionStatus = "stopped"
)

type managedSession interface {
	info() SessionInfo
	sendCommand(ctx context.Context, command string) error
	readOutput(ctx context.Context, timeout time.Duration) (string, error)
	hasBufferedOutput() bool
	stop() error
}

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
	stopping          bool
	logger            *slog.Logger
}

type SpawnSession struct {
	ID     string        `json:"session_id"`
	Role   string        `json:"role"`
	Agent  string        `json:"agent"`
	Status SessionStatus `json:"status"`
	PID    int           `json:"pid,omitempty"`

	launch spawnLauncherFunc

	outputMu    sync.Mutex
	sessionID   string
	lastOutput  string
	outputReady bool
	currentCmd  *exec.Cmd
	currentDone chan struct{}
	waitErr     error
	stopping    bool
	logger      *slog.Logger
}

type SessionManager struct {
	mu       sync.Mutex
	sessions map[string]managedSession
	launch   launcherFunc
	spawn    spawnLauncherFunc
	logger   *slog.Logger
}

type launcherFunc func(ctx context.Context, role, agent string) (*exec.Cmd, error)

type spawnLauncherFunc func(ctx context.Context, req spawnRequest) (*exec.Cmd, error)

type spawnRequest struct {
	Role      string
	Agent     string
	Prompt    string
	SessionID string
	Resume    bool
}

type sessionConfig struct {
	Roles map[string]roleConfig `json:"roles"`
}

type roleConfig struct {
	Model string `json:"model"`
}

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
	return newSessionManager(defaultLauncher, defaultSpawnLauncher, logger)
}

func newSessionManager(launch launcherFunc, spawn spawnLauncherFunc, logger *slog.Logger) *SessionManager {
	if logger == nil {
		logger = newDiscardLogger()
	}

	return &SessionManager{
		sessions: make(map[string]managedSession),
		launch:   launch,
		spawn:    spawn,
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

func defaultSpawnLauncher(ctx context.Context, req spawnRequest) (*exec.Cmd, error) {
	if err := validateRole(req.Role); err != nil {
		return nil, err
	}
	if err := validateAgent(req.Agent); err != nil {
		return nil, err
	}
	if req.Agent != "codex" {
		return nil, fmt.Errorf("spawn sessions are unsupported for agent %q", req.Agent)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("determine working directory: %w", err)
	}

	prompt := strings.TrimSpace(req.Prompt)
	if !req.Resume {
		prompt, err = loadRolePrompt(cwd, req.Role)
		if err != nil {
			return nil, err
		}
	}

	roleModel, err := loadRoleModel(cwd, req.Role)
	if err != nil {
		return nil, err
	}

	args := []string{"exec"}
	if req.Resume {
		args = append(args, "resume")
		if req.SessionID != "" {
			args = append(args, req.SessionID)
		} else {
			args = append(args, "--last")
		}
	} else {
		args = append(args, "--sandbox", "workspace-write")
	}
	args = append(args, "-c", "sandbox_workspace_write.network_access=true")
	if roleModel != "" {
		args = append(args, "-m", roleModel)
	}
	args = append(args, "-")

	cmd := exec.Command("codex", args...)
	cmd.Dir = cwd
	cmd.Stdin = strings.NewReader(prompt + "\n")
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
		if existing.info().Status == SessionStatusRunning {
			m.mu.Unlock()
			err := fmt.Errorf("session for role %q is already running", role)
			m.logger.Error("start session failed", "role", role, "agent", agent, "error", err)
			return SessionInfo{}, err
		}
		delete(m.sessions, role)
	}
	m.mu.Unlock()

	var (
		session managedSession
		err     error
	)

	if isSpawnAgent(agent) {
		spawnSession := newSpawnSession(role, agent, m.spawn, m.logger)
		if err = spawnSession.start(ctx); err == nil {
			session = spawnSession
		}
	} else {
		session, err = m.startProcessSession(ctx, role, agent)
	}
	if err != nil {
		m.logger.Error("start session failed", "role", role, "agent", agent, "error", err)
		return SessionInfo{}, err
	}

	m.mu.Lock()
	m.sessions[role] = session
	m.mu.Unlock()

	info := session.info()
	m.logger.Info("session started", "role", role, "agent", agent, "pid", info.PID, "session_id", info.SessionID)
	return info, nil
}

func (m *SessionManager) startProcessSession(ctx context.Context, role, agent string) (managedSession, error) {
	cmd, err := m.launch(ctx, role, agent)
	if err != nil {
		return nil, err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("open stdin pipe: %w", err)
	}

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		_ = stdin.Close()
		return nil, fmt.Errorf("create stdout pipe: %w", err)
	}
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stdoutWriter

	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		_ = stdoutReader.Close()
		_ = stdoutWriter.Close()
		return nil, fmt.Errorf("start session process: %w", err)
	}
	if err := stdoutWriter.Close(); err != nil {
		_ = stdin.Close()
		_ = stdoutReader.Close()
		return nil, fmt.Errorf("close stdout writer: %w", err)
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

	go session.captureOutput()
	go session.waitForExit()
	session.captureStartupOutput()

	return session, nil
}

func (m *SessionManager) StopSession(role string) (SessionInfo, error) {
	m.mu.Lock()
	session, ok := m.sessions[role]
	m.mu.Unlock()
	if !ok {
		err := fmt.Errorf("no session for role %q", role)
		m.logger.Error("stop session failed", "role", role, "error", err)
		return SessionInfo{}, err
	}

	if err := session.stop(); err != nil {
		info := session.info()
		m.logger.Error("stop session failed", "role", role, "session_id", info.SessionID, "error", err)
		return SessionInfo{}, err
	}

	m.mu.Lock()
	delete(m.sessions, role)
	m.mu.Unlock()

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

	info := session.info()
	if info.Status != SessionStatusRunning {
		err := fmt.Errorf("session for role %q is not running", role)
		m.logger.Error("send command failed", "role", role, "command", command, "session_id", info.SessionID, "error", err)
		return CommandResult{}, err
	}

	m.logger.Info("sending command", "role", role, "command", command, "session_id", info.SessionID)
	if err := session.sendCommand(ctx, command); err != nil {
		m.logger.Error("send command failed", "role", role, "command", command, "session_id", info.SessionID, "error", err)
		return CommandResult{}, err
	}

	info = session.info()
	result := CommandResult{
		Role:      role,
		Command:   command,
		SessionID: info.SessionID,
	}
	m.logger.Info("command queued", "role", role, "command", command, "session_id", info.SessionID)
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

	info := session.info()
	if info.Status != SessionStatusRunning && !session.hasBufferedOutput() {
		err := fmt.Errorf("session for role %q is not running", role)
		m.logger.Error("get output failed", "role", role, "session_id", info.SessionID, "error", err)
		return OutputResult{}, err
	}

	output, err := session.readOutput(ctx, timeout)
	if err != nil {
		m.logger.Error("get output failed", "role", role, "session_id", info.SessionID, "error", err)
		return OutputResult{}, err
	}

	info = session.info()
	result := OutputResult{
		Role:      role,
		Output:    output,
		SessionID: info.SessionID,
		Status:    info.Status,
	}
	m.logger.Info("output retrieved", "role", role, "session_id", info.SessionID, "status", result.Status, "output_bytes", len(output))
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

func (s *Session) sendCommand(_ context.Context, command string) error {
	return s.writeCommand(command)
}

func (s *Session) stop() error {
	if s.Status == SessionStatusRunning && s.cmd.Process != nil {
		s.setStopping(true)
		s.logger.Info("sending session signal", "role", s.Role, "session_id", s.ID, "signal", "SIGTERM")
		if err := s.cmd.Process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
			s.setStopping(false)
			return fmt.Errorf("stop session %q: %w", s.Role, err)
		}
		if !waitForSessionExit(s.done, stopTermGracePeriod) {
			s.logger.Warn("session still running after grace period; escalating stop", "role", s.Role, "session_id", s.ID, "signal", "SIGKILL")
			if err := s.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
				s.setStopping(false)
				return fmt.Errorf("kill session %q: %w", s.Role, err)
			}
			s.logger.Info("sending session signal", "role", s.Role, "session_id", s.ID, "signal", "SIGKILL")
			if !waitForSessionExit(s.done, stopKillGracePeriod) {
				s.setStopping(false)
				return fmt.Errorf("session %q did not exit after SIGKILL", s.Role)
			}
		}
	}

	s.outputMu.Lock()
	if s.Status == SessionStatusRunning {
		s.Status = SessionStatusStopped
	}
	s.outputMu.Unlock()
	return nil
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
		if s.stopping || err == nil {
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
			return s.outputSince(startOffset), nil
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

func (s *Session) setStopping(stopping bool) {
	s.outputMu.Lock()
	defer s.outputMu.Unlock()
	s.stopping = stopping
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

func newSpawnSession(role, agent string, launch spawnLauncherFunc, logger *slog.Logger) *SpawnSession {
	if logger == nil {
		logger = newDiscardLogger()
	}

	return &SpawnSession{
		ID:     fmt.Sprintf("%s-%s-%d", role, agent, time.Now().UTC().UnixNano()),
		Role:   role,
		Agent:  agent,
		Status: SessionStatusRunning,
		launch: launch,
		logger: logger,
	}
}

func (s *SpawnSession) start(ctx context.Context) error {
	cmd, err := s.launch(ctx, spawnRequest{
		Role:  s.Role,
		Agent: s.Agent,
	})
	if err != nil {
		return err
	}

	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)
	sessionID := extractCodexSessionID(output)

	s.outputMu.Lock()
	s.lastOutput = output
	s.outputReady = true
	s.waitErr = err
	if sessionID != "" {
		s.sessionID = sessionID
		s.ID = sessionID
	}
	if err != nil {
		s.Status = SessionStatusExited
	}
	s.outputMu.Unlock()

	if err != nil {
		return fmt.Errorf("start spawn session: %w", err)
	}
	return nil
}

func (s *SpawnSession) info() SessionInfo {
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

func (s *SpawnSession) sendCommand(ctx context.Context, command string) error {
	command = strings.TrimSpace(command)
	if command == "" {
		return fmt.Errorf("command must not be empty")
	}

	s.outputMu.Lock()
	if s.currentCmd != nil {
		s.outputMu.Unlock()
		return fmt.Errorf("session for role %q already has a command in flight", s.Role)
	}
	sessionID := s.sessionID
	s.lastOutput = ""
	s.outputReady = false
	s.waitErr = nil
	done := make(chan struct{})
	s.currentDone = done
	s.outputMu.Unlock()

	cmd, err := s.launch(ctx, spawnRequest{
		Role:      s.Role,
		Agent:     s.Agent,
		Prompt:    command,
		SessionID: sessionID,
		Resume:    true,
	})
	if err != nil {
		s.outputMu.Lock()
		s.currentDone = nil
		s.waitErr = err
		s.Status = SessionStatusExited
		s.outputMu.Unlock()
		return err
	}

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		s.outputMu.Lock()
		s.currentDone = nil
		s.waitErr = err
		s.Status = SessionStatusExited
		s.outputMu.Unlock()
		return fmt.Errorf("create spawn stdout pipe: %w", err)
	}
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stdoutWriter

	if err := cmd.Start(); err != nil {
		_ = stdoutReader.Close()
		_ = stdoutWriter.Close()
		s.outputMu.Lock()
		s.currentDone = nil
		s.waitErr = err
		s.Status = SessionStatusExited
		s.outputMu.Unlock()
		return fmt.Errorf("start spawn command: %w", err)
	}
	if err := stdoutWriter.Close(); err != nil {
		_ = cmd.Process.Kill()
		_ = stdoutReader.Close()
		s.outputMu.Lock()
		s.currentDone = nil
		s.waitErr = err
		s.Status = SessionStatusExited
		s.outputMu.Unlock()
		return fmt.Errorf("close spawn stdout writer: %w", err)
	}

	s.outputMu.Lock()
	s.currentCmd = cmd
	s.PID = cmd.Process.Pid
	s.outputMu.Unlock()

	go s.waitForCommand(cmd, stdoutReader, done)
	return nil
}

func (s *SpawnSession) waitForCommand(cmd *exec.Cmd, stdout io.ReadCloser, done chan struct{}) {
	defer close(done)
	defer stdout.Close()

	outputBytes, readErr := io.ReadAll(stdout)
	waitErr := cmd.Wait()
	output := string(outputBytes)
	sessionID := extractCodexSessionID(output)

	s.outputMu.Lock()
	defer s.outputMu.Unlock()

	if readErr != nil && waitErr == nil {
		waitErr = readErr
	}
	if sessionID != "" {
		s.sessionID = sessionID
		s.ID = sessionID
	}
	s.lastOutput = output
	s.outputReady = true
	s.waitErr = waitErr
	s.currentCmd = nil
	s.currentDone = nil
	s.PID = 0
	if s.Status == SessionStatusRunning && waitErr != nil {
		if s.stopping {
			s.Status = SessionStatusStopped
		} else {
			s.Status = SessionStatusExited
		}
	}
}

func (s *SpawnSession) readOutput(ctx context.Context, timeout time.Duration) (string, error) {
	output, ready, waitErr, done := s.outputState()
	if ready {
		if output != "" {
			return output, nil
		}
		if waitErr != nil {
			return "", waitErr
		}
		return "", nil
	}
	if done == nil {
		if waitErr != nil {
			return "", waitErr
		}
		return output, nil
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		output, ready, _, _ = s.outputState()
		if ready && output != "" {
			return output, nil
		}
		return "", ctx.Err()
	case <-timer.C:
		output, ready, waitErr, _ = s.outputState()
		if ready {
			if output != "" {
				return output, nil
			}
			if waitErr != nil {
				return "", waitErr
			}
		}
		return "", nil
	case <-done:
		output, _, waitErr, _ = s.outputState()
		if output != "" {
			return output, nil
		}
		if waitErr != nil {
			return "", waitErr
		}
		return "", nil
	}
}

func (s *SpawnSession) outputState() (string, bool, error, chan struct{}) {
	s.outputMu.Lock()
	defer s.outputMu.Unlock()
	return s.lastOutput, s.outputReady, s.waitErr, s.currentDone
}

func (s *SpawnSession) hasBufferedOutput() bool {
	s.outputMu.Lock()
	defer s.outputMu.Unlock()
	return s.outputReady && s.lastOutput != ""
}

func (s *SpawnSession) stop() error {
	s.outputMu.Lock()
	s.stopping = true
	cmd := s.currentCmd
	done := s.currentDone
	s.outputMu.Unlock()

	if cmd != nil && cmd.Process != nil {
		s.logger.Info("sending session signal", "role", s.Role, "session_id", s.ID, "signal", "SIGTERM")
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
			s.outputMu.Lock()
			s.stopping = false
			s.outputMu.Unlock()
			return fmt.Errorf("stop session %q: %w", s.Role, err)
		}
		if !waitForSessionExit(done, stopTermGracePeriod) {
			s.logger.Warn("session still running after grace period; escalating stop", "role", s.Role, "session_id", s.ID, "signal", "SIGKILL")
			if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
				s.outputMu.Lock()
				s.stopping = false
				s.outputMu.Unlock()
				return fmt.Errorf("kill session %q: %w", s.Role, err)
			}
			s.logger.Info("sending session signal", "role", s.Role, "session_id", s.ID, "signal", "SIGKILL")
			if !waitForSessionExit(done, stopKillGracePeriod) {
				s.outputMu.Lock()
				s.stopping = false
				s.outputMu.Unlock()
				return fmt.Errorf("session %q did not exit after SIGKILL", s.Role)
			}
		}
	}

	s.outputMu.Lock()
	s.Status = SessionStatusStopped
	s.PID = 0
	s.currentCmd = nil
	s.currentDone = nil
	s.outputMu.Unlock()
	return nil
}

func waitForSessionExit(done <-chan struct{}, timeout time.Duration) bool {
	if done == nil {
		return true
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-done:
		return true
	case <-timer.C:
		return false
	}
}

func isSpawnAgent(agent string) bool {
	return agent == "codex"
}

func loadRolePrompt(cwd, role string) (string, error) {
	promptFile, err := promptFileForRole(cwd, role)
	if err != nil {
		return "", err
	}

	promptBytes, err := os.ReadFile(promptFile)
	if err != nil {
		return "", fmt.Errorf("read prompt file %q: %w", promptFile, err)
	}
	return string(promptBytes), nil
}

func promptFileForRole(cwd, role string) (string, error) {
	var promptFile string
	switch role {
	case "plan":
		promptFile = ".ai/prompts/planner.md"
	case "implement":
		promptFile = ".ai/prompts/implementer.md"
	case "review":
		promptFile = ".ai/prompts/reviewer.md"
	default:
		return "", fmt.Errorf("unsupported role %q", role)
	}

	path := filepath.Join(cwd, promptFile)
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("locate prompt file: %w", err)
	}
	return path, nil
}

func loadRoleModel(cwd, role string) (string, error) {
	configPath := filepath.Join(cwd, ".ai", "config.json")
	configBytes, err := os.ReadFile(configPath)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read config file %q: %w", configPath, err)
	}

	var config sessionConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return "", fmt.Errorf("parse config file %q: %w", configPath, err)
	}
	return config.Roles[role].Model, nil
}

func extractCodexSessionID(output string) string {
	matches := codexSessionIDPattern.FindStringSubmatch(output)
	if len(matches) != 2 {
		return ""
	}
	return matches[1]
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
