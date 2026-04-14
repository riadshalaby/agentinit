package mcp

import (
	"context"
	"time"
)

// StartOpts are passed to Adapter.Start.
type StartOpts struct {
	PromptFile string        // path to role prompt file (required)
	Model      string        // provider-specific model string (optional)
	Effort     string        // claude: --effort value (optional)
	Timeout    time.Duration // 0 means no timeout
}

// RunOpts are passed to Adapter.Run.
type RunOpts struct {
	Model   string
	Timeout time.Duration // 0 means no timeout
}

// Adapter handles provider-specific CLI invocation.
// Each method spawns a short-lived subprocess and returns its full output.
type Adapter interface {
	// Start runs the initial CLI invocation with the role system prompt.
	// It updates session.ProviderState in place.
	Start(ctx context.Context, session *Session, opts StartOpts) (output string, err error)

	// Run resumes the session with a command.
	// It updates session.ProviderState in place.
	Run(ctx context.Context, session *Session, command string, opts RunOpts) (output string, err error)

	// Stop kills the process identified by session.ProviderState if it is
	// currently running. No-op if nothing is running.
	Stop(ctx context.Context, session *Session) error
}
