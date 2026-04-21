package mcp

import (
	"context"
	"io"
	"time"
)

// StartOpts are passed to Adapter.Start.
type StartOpts struct {
	PromptFile string        // path to role prompt file (required)
	Model      string        // provider-specific model string (optional)
	Effort     string        // provider-specific reasoning-effort value (optional)
	Timeout    time.Duration // 0 means no timeout
}

// RunOpts are passed to Adapter.RunStream.
type RunOpts struct {
	Model string
}

// Adapter handles provider-specific CLI invocation.
// Each method spawns a short-lived provider subprocess.
type Adapter interface {
	// Start runs the initial CLI invocation with the role system prompt.
	// It updates session.ProviderState in place.
	Start(ctx context.Context, session *Session, opts StartOpts) (output string, err error)

	// RunStream resumes the session with a command and streams stdout+stderr to w.
	// It updates session.ProviderState in place.
	RunStream(ctx context.Context, session *Session, command string, opts RunOpts, w io.Writer) error

	// Stop kills the process identified by session.ProviderState if it is
	// currently running. No-op if nothing is running.
	Stop(ctx context.Context, session *Session) error
}
