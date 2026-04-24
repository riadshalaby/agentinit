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

type RunResult struct {
	Status       SessionStatus `json:"status"`
	Error        string        `json:"error,omitempty"`
	ExitSummary  string        `json:"exit_summary,omitempty"`
	DurationSecs float64       `json:"duration_secs"`
}

// WaitResult returns the latest structured run outcome together with the
// current session metadata for callers that block on completion.
type WaitResult struct {
	Session SessionInfo `json:"session"`
	Result  *RunResult  `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type Session struct {
	Name          string        `json:"name"`
	Role          string        `json:"role"`
	Provider      string        `json:"provider"`
	Model         string        `json:"model,omitempty"`
	Status        SessionStatus `json:"status"`
	ProviderState ProviderState `json:"provider_state"`
	Result        *RunResult    `json:"result,omitempty"`
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
