package mcp

import "time"

type SessionStatus string

const (
	StatusIdle    SessionStatus = "idle"
	StatusRunning SessionStatus = "running"
	StatusErrored SessionStatus = "errored"
	StatusStopped SessionStatus = "stopped"
)

type ProviderState struct {
	SessionID string `json:"session_id,omitempty"`
}

type Session struct {
	Name          string        `json:"name"`
	Role          string        `json:"role"`
	Provider      string        `json:"provider"`
	Model         string        `json:"model,omitempty"`
	Status        SessionStatus `json:"status"`
	ProviderState ProviderState `json:"provider_state"`
	CreatedAt     time.Time     `json:"created_at"`
	LastActiveAt  time.Time     `json:"last_active_at"`
	RunCount      int           `json:"run_count"`
	Error         string        `json:"error,omitempty"`
}

type SessionInfo struct {
	Name     string        `json:"name"`
	Role     string        `json:"role"`
	Provider string        `json:"provider"`
	Status   SessionStatus `json:"status"`
	RunCount int           `json:"run_count"`
	Error    string        `json:"error,omitempty"`
}

func (s *Session) info() SessionInfo {
	return SessionInfo{
		Name:     s.Name,
		Role:     s.Role,
		Provider: s.Provider,
		Status:   s.Status,
		RunCount: s.RunCount,
		Error:    s.Error,
	}
}
