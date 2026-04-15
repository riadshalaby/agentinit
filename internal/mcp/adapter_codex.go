package mcp

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var codexSessionIDPattern = regexp.MustCompile(`(?m)^session id:\s+(\S+)$`)

type codexExecFunc func(ctx context.Context, args []string, stdin string, w io.Writer) error

type CodexAdapter struct {
	cwd     string
	sandbox string
	network bool
	exec    codexExecFunc
}

func NewCodexAdapter(cwd string, defaults CodexDefaults) *CodexAdapter {
	sandbox := defaults.Sandbox
	if sandbox == "" {
		sandbox = "workspace-write"
	}
	a := &CodexAdapter{cwd: cwd, sandbox: sandbox, network: defaults.NetworkAccess}
	a.exec = a.defaultExec
	return a
}

func (a *CodexAdapter) Start(ctx context.Context, session *Session, opts StartOpts) (string, error) {
	prompt, err := readPromptFile(opts.PromptFile)
	if err != nil {
		return "", err
	}

	args := []string{"exec", "--sandbox", a.sandbox}
	if a.network {
		args = append(args, "-c", fmt.Sprintf("sandbox_%s.network_access=true", strings.ReplaceAll(a.sandbox, "-", "_")))
	}
	if opts.Model != "" {
		args = append(args, "-m", opts.Model)
	}
	args = append(args, "-")

	var sb strings.Builder
	err = a.exec(ctx, args, prompt, &sb)
	output := sb.String()
	if sessionID := extractCodexSessionID(output); sessionID != "" {
		session.ProviderState.SessionID = sessionID
	}
	return output, err
}

func (a *CodexAdapter) RunStream(ctx context.Context, session *Session, command string, opts RunOpts, w io.Writer) error {
	if session.ProviderState.SessionID == "" {
		return fmt.Errorf("session %q has no provider session ID; call Start first", session.Name)
	}

	args := []string{"exec", "resume", session.ProviderState.SessionID}
	if a.network {
		args = append(args, "-c", fmt.Sprintf("sandbox_%s.network_access=true", strings.ReplaceAll(a.sandbox, "-", "_")))
	}
	if opts.Model != "" {
		args = append(args, "-m", opts.Model)
	}
	args = append(args, "-")

	var sb strings.Builder
	mw := io.MultiWriter(w, &sb)
	err := a.exec(ctx, args, command, mw)
	output := sb.String()
	if sessionID := extractCodexSessionID(output); sessionID != "" {
		session.ProviderState.SessionID = sessionID
	}
	return err
}

func (a *CodexAdapter) Stop(_ context.Context, _ *Session) error {
	return nil
}

func (a *CodexAdapter) defaultExec(ctx context.Context, args []string, stdin string, w io.Writer) error {
	cmd := exec.CommandContext(ctx, "codex", args...)
	cmd.Dir = a.cwd
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

func extractCodexSessionID(output string) string {
	matches := codexSessionIDPattern.FindStringSubmatch(output)
	if len(matches) != 2 {
		return ""
	}
	return matches[1]
}

func promptFileForRole(cwd, role string) (string, error) {
	var name string
	switch role {
	case "implement":
		name = "implementer.md"
	case "review":
		name = "reviewer.md"
	default:
		return "", fmt.Errorf("no prompt file for role %q", role)
	}

	path := filepath.Join(cwd, ".ai", "prompts", name)
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("locate prompt file for role %q: %w", role, err)
	}
	return path, nil
}

func readPromptFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read prompt file %q: %w", path, err)
	}
	return string(data), nil
}
