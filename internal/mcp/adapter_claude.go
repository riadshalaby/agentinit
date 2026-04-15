package mcp

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type claudeExecFunc func(ctx context.Context, args []string, w io.Writer) error

type ClaudeAdapter struct {
	cwd            string
	permissionMode string
	exec           claudeExecFunc
}

func NewClaudeAdapter(cwd string, defaults ClaudeDefaults) *ClaudeAdapter {
	permissionMode := defaults.PermissionMode
	if permissionMode == "" {
		permissionMode = "acceptEdits"
	}

	a := &ClaudeAdapter{cwd: cwd, permissionMode: permissionMode}
	a.exec = a.defaultExec
	return a
}

func (a *ClaudeAdapter) Start(ctx context.Context, session *Session, opts StartOpts) (string, error) {
	if session.ProviderState.SessionID == "" {
		return "", fmt.Errorf("session %q has no session ID; caller must set one before Start", session.Name)
	}

	args := []string{
		"-p",
		"--session-id", session.ProviderState.SessionID,
		"--permission-mode", a.permissionMode,
	}
	if opts.PromptFile != "" {
		args = append(args, "--system-prompt-file", opts.PromptFile)
	}
	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}
	if opts.Effort != "" {
		args = append(args, "--effort", opts.Effort)
	}
	args = append(args, "You are now in WAIT_FOR_USER_START state.")

	var sb strings.Builder
	err := a.exec(ctx, args, &sb)
	return sb.String(), err
}

func (a *ClaudeAdapter) RunStream(ctx context.Context, session *Session, command string, opts RunOpts, w io.Writer) error {
	if session.ProviderState.SessionID == "" {
		return fmt.Errorf("session %q has no provider session ID; call Start first", session.Name)
	}

	args := []string{
		"-p",
		"--session-id", session.ProviderState.SessionID,
		"--permission-mode", a.permissionMode,
	}
	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}
	args = append(args, command)

	return a.exec(ctx, args, w)
}

func (a *ClaudeAdapter) Stop(_ context.Context, _ *Session) error {
	return nil
}

func (a *ClaudeAdapter) defaultExec(ctx context.Context, args []string, w io.Writer) error {
	cmd := exec.CommandContext(ctx, "claude", args...)
	cmd.Dir = a.cwd
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}
