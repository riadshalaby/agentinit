package mcp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const runResultExitSummaryLimit = 500
const sessionWaitPollInterval = 25 * time.Millisecond

// SessionManager owns the named session registry.
// It is the single entry point for all session lifecycle operations.
type SessionManager struct {
	ctx      context.Context
	store    *Store
	adapters map[string]Adapter
	config   Config
	cwd      string
	mu       sync.Mutex
	running  map[string]context.CancelFunc
	outputs  map[string]*outputBuffer
	logger   *slog.Logger
}

func NewSessionManager(ctx context.Context, store *Store, adapters map[string]Adapter, config Config, cwd string, logger *slog.Logger) *SessionManager {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		logger = newDiscardLogger()
	}
	if store == nil {
		store = NewStore("")
	}

	m := &SessionManager{
		ctx:      ctx,
		store:    store,
		adapters: adapters,
		config:   config,
		cwd:      cwd,
		running:  make(map[string]context.CancelFunc),
		outputs:  make(map[string]*outputBuffer),
		logger:   logger,
	}
	m.recoverStaleRunning()
	return m
}

// recoverStaleRunning marks any session persisted as StatusRunning as StatusErrored
// on startup because the previous MCP server died mid-run.
func (m *SessionManager) recoverStaleRunning() {
	sessions, err := m.store.List()
	if err != nil {
		m.logger.Error("recover stale sessions failed", "error", err)
		return
	}
	for _, session := range sessions {
		if session.Status != StatusRunning {
			continue
		}
		session.Status = StatusErrored
		session.Error = "recovered after interrupted run"
		if err := m.store.Put(session); err != nil {
			m.logger.Error("persist recovered session failed", "name", session.Name, "error", err)
		}
	}
}

// StartSession creates a new named session, runs the provider CLI with the role
// prompt, and persists the session. Returns an error if the name is already in use.
func (m *SessionManager) StartSession(ctx context.Context, name, role, provider string) (SessionInfo, string, error) {
	if name == "" {
		return SessionInfo{}, "", fmt.Errorf("session name must not be empty")
	}
	if err := validateRole(role); err != nil {
		return SessionInfo{}, "", err
	}
	if err := validateProvider(provider); err != nil {
		return SessionInfo{}, "", err
	}

	if _, err := m.store.Get(name); err == nil {
		return SessionInfo{}, "", fmt.Errorf("session %q already exists", name)
	} else if !strings.Contains(err.Error(), "not found") {
		return SessionInfo{}, "", err
	}

	adapter, ok := m.adapters[provider]
	if !ok {
		return SessionInfo{}, "", fmt.Errorf("adapter %q is not configured", provider)
	}

	promptFile, err := promptFileForRole(m.cwd, role)
	if err != nil {
		return SessionInfo{}, "", err
	}

	now := time.Now().UTC()
	model := m.config.ModelForRoleAndProvider(role, provider)
	effort := m.config.EffortForRoleAndProvider(role, provider)
	session := &Session{
		Name:         name,
		Role:         role,
		Provider:     provider,
		Model:        model,
		Status:       StatusIdle,
		CreatedAt:    now,
		LastActiveAt: now,
	}
	if provider == "claude" {
		session.ProviderState.SessionID = uuid.NewString()
	}

	output, err := adapter.Start(ctx, session, StartOpts{
		PromptFile: promptFile,
		Model:      model,
		Effort:     effort,
	})
	if err != nil {
		session.Status = StatusErrored
		session.Error = err.Error()
		if putErr := m.store.Put(session); putErr != nil {
			return SessionInfo{}, "", fmt.Errorf("start session failed: %w (persist error: %v)", err, putErr)
		}
		return session.info(), output, err
	}

	session.Status = StatusIdle
	session.Error = ""
	if err := m.store.Put(session); err != nil {
		return SessionInfo{}, "", err
	}
	return session.info(), output, nil
}

// RunSession sends a command to an existing session asynchronously.
// Returns an error if the session is already running.
func (m *SessionManager) RunSession(ctx context.Context, name, command string) (SessionInfo, error) {
	session, err := m.store.Get(name)
	if err != nil {
		return SessionInfo{}, err
	}
	adapter, ok := m.adapters[session.Provider]
	if !ok {
		return SessionInfo{}, fmt.Errorf("adapter %q is not configured", session.Provider)
	}

	runCtx, cancel := context.WithCancel(m.ctx)
	m.mu.Lock()
	if _, exists := m.running[name]; exists {
		m.mu.Unlock()
		cancel()
		return SessionInfo{}, fmt.Errorf("session %q is already running", name)
	}
	m.running[name] = cancel
	m.mu.Unlock()

	buf := &outputBuffer{}
	m.mu.Lock()
	m.outputs[name] = buf
	m.mu.Unlock()

	session.Result = nil
	session.Status = StatusRunning
	session.Error = ""
	if err := m.store.Put(session); err != nil {
		cancel()
		m.mu.Lock()
		delete(m.running, name)
		delete(m.outputs, name)
		m.mu.Unlock()
		return SessionInfo{}, err
	}

	runStartedAt := time.Now().UTC()
	go func() {
		defer func() {
			cancel()
			m.mu.Lock()
			delete(m.running, name)
			m.mu.Unlock()
		}()

		runErr := adapter.RunStream(runCtx, session, command, RunOpts{Model: session.Model}, buf)
		current, err := m.store.Get(name)
		if err != nil {
			m.logger.Error("load session after run failed", "name", name, "error", err)
			return
		}

		finishedAt := time.Now().UTC()
		current.ProviderState = session.ProviderState
		current.LastActiveAt = finishedAt
		if runErr != nil {
			if errors.Is(runErr, context.Canceled) || errors.Is(runErr, context.DeadlineExceeded) {
				current.Status = StatusStopped
				current.Error = ""
			} else {
				current.Status = StatusErrored
				current.Error = runErr.Error()
			}
		} else {
			current.Status = StatusIdle
			current.RunCount++
			current.Error = ""
		}
		current.Result = &RunResult{
			Status:       current.Status,
			Error:        current.Error,
			ExitSummary:  buf.Tail(runResultExitSummaryLimit),
			DurationSecs: finishedAt.Sub(runStartedAt).Seconds(),
		}

		if err := m.store.Put(current); err != nil {
			m.logger.Error("persist session after run failed", "name", name, "error", err)
		}
	}()

	return session.info(), nil
}

func (m *SessionManager) GetOutput(name string, offset, limit int) (chunk string, totalBytes int, running bool, err error) {
	session, err := m.store.Get(name)
	if err != nil {
		return "", 0, false, err
	}
	m.mu.Lock()
	buf := m.outputs[name]
	m.mu.Unlock()
	if buf == nil {
		return "", 0, session.Status == StatusRunning, nil
	}
	chunk, total := buf.StringFromLimit(offset, limit)
	return chunk, total, session.Status == StatusRunning, nil
}

func (m *SessionManager) GetResult(name string) (*RunResult, error) {
	session, err := m.store.Get(name)
	if err != nil {
		return nil, err
	}
	return session.Result, nil
}

// WaitSession blocks until the named session run has fully settled or the
// caller context is canceled. It returns the latest session info and any
// structured run result without exposing raw output.
func (m *SessionManager) WaitSession(ctx context.Context, name string) (SessionInfo, *RunResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	for {
		session, err := m.store.Get(name)
		if err != nil {
			return SessionInfo{}, nil, err
		}
		info := session.info()
		if session.Status != StatusRunning && !m.isSessionRunActive(name) {
			return info, session.Result, nil
		}

		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return info, session.Result, fmt.Errorf("session %q wait timed out", name)
			}
			return info, session.Result, ctx.Err()
		case <-time.After(sessionWaitPollInterval):
		}
	}
}

// StopSession cancels any in-flight RunSession for the named session.
func (m *SessionManager) StopSession(name string) (SessionInfo, error) {
	session, err := m.store.Get(name)
	if err != nil {
		return SessionInfo{}, err
	}
	adapter, ok := m.adapters[session.Provider]
	if !ok {
		return SessionInfo{}, fmt.Errorf("adapter %q is not configured", session.Provider)
	}

	m.mu.Lock()
	cancel := m.running[name]
	m.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	if err := adapter.Stop(context.Background(), session); err != nil {
		return SessionInfo{}, err
	}

	session.Status = StatusStopped
	session.Error = ""
	session.LastActiveAt = time.Now().UTC()
	if err := m.store.Put(session); err != nil {
		return SessionInfo{}, err
	}
	return session.info(), nil
}

// ResetSession clears provider state so the next Run starts a fresh conversation.
func (m *SessionManager) ResetSession(name string) (SessionInfo, error) {
	session, err := m.store.Get(name)
	if err != nil {
		return SessionInfo{}, err
	}
	m.mu.Lock()
	delete(m.outputs, name)
	m.mu.Unlock()

	session.ProviderState = ProviderState{}
	session.Result = nil
	session.Status = StatusIdle
	session.Error = ""
	session.LastActiveAt = time.Now().UTC()
	if err := m.store.Put(session); err != nil {
		return SessionInfo{}, err
	}
	return session.info(), nil
}

// DeleteSession removes the session entirely.
func (m *SessionManager) DeleteSession(name string) error {
	m.mu.Lock()
	if cancel := m.running[name]; cancel != nil {
		cancel()
		delete(m.running, name)
	}
	delete(m.outputs, name)
	m.mu.Unlock()
	return m.store.Delete(name)
}

// GetSession returns the current SessionInfo for a named session.
func (m *SessionManager) GetSession(name string) (SessionInfo, error) {
	session, err := m.store.Get(name)
	if err != nil {
		return SessionInfo{}, err
	}
	return session.info(), nil
}

// ListSessions returns info for all tracked sessions.
func (m *SessionManager) ListSessions() ([]SessionInfo, error) {
	sessions, err := m.store.List()
	if err != nil {
		return nil, err
	}
	infos := make([]SessionInfo, 0, len(sessions))
	for _, session := range sessions {
		infos = append(infos, session.info())
	}
	slices.SortFunc(infos, func(a, b SessionInfo) int {
		switch {
		case a.Name < b.Name:
			return -1
		case a.Name > b.Name:
			return 1
		default:
			return 0
		}
	})
	return infos, nil
}

func (m *SessionManager) isSessionRunActive(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.running[name]
	return ok
}

func validateProvider(provider string) error {
	if _, ok := validProviders[provider]; !ok {
		return fmt.Errorf("unsupported provider %q", provider)
	}
	return nil
}

func validateRole(role string) error {
	if _, ok := validRoles[role]; !ok {
		return fmt.Errorf("unsupported role %q: must be one of: implement, po, review", role)
	}
	return nil
}
