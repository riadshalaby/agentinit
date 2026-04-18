//go:build e2e

package e2e_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/riadshalaby/agentinit/internal/mcp"
)

const stubPrompt = "You are a test agent. Respond concisely to the user's message."

// TestMCPSessionLifecycle exercises SessionManager end-to-end using real CLI
// adapters. It skips cleanly when claude or codex are absent from PATH.
func TestMCPSessionLifecycle(t *testing.T) {
	if _, err := exec.LookPath("claude"); err != nil {
		t.Skip("claude not found in PATH; skipping real-agent E2E test")
	}
	if _, err := exec.LookPath("codex"); err != nil {
		t.Skip("codex not found in PATH; skipping real-agent E2E test")
	}

	// Set up a temp directory with stub prompt files.
	// Codex requires a trusted (git) directory; initialize a bare repo so it
	// passes the git-repo check it performs before accepting any command.
	tmpDir := t.TempDir()
	if out, err := exec.Command("git", "init", tmpDir).CombinedOutput(); err != nil {
		t.Fatalf("git init temp dir: %v\n%s", err, out)
	}

	promptsDir := filepath.Join(tmpDir, ".ai", "prompts")
	if err := os.MkdirAll(promptsDir, 0o755); err != nil {
		t.Fatalf("create prompts dir: %v", err)
	}
	for _, name := range []string{"implementer.md", "reviewer.md"} {
		if err := os.WriteFile(filepath.Join(promptsDir, name), []byte(stubPrompt), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	store := mcp.NewStore(filepath.Join(tmpDir, ".ai", "sessions.json"))
	adapters := map[string]mcp.Adapter{
		"claude": mcp.NewClaudeAdapter(tmpDir, mcp.ClaudeDefaults{PermissionMode: "acceptEdits"}),
		"codex":  mcp.NewCodexAdapter(tmpDir, mcp.CodexDefaults{Sandbox: "workspace-write"}),
	}
	mgr := mcp.NewSessionManager(context.Background(), store, adapters, mcp.Config{}, tmpDir, nil)

	t.Run("codex implementer session", func(t *testing.T) {
		ctx := context.Background()

		info, _, err := mgr.StartSession(ctx, "implementer", "implement", "codex")
		if err != nil {
			t.Fatalf("StartSession: %v", err)
		}
		if info.Status != mcp.StatusIdle {
			t.Fatalf("status after Start = %q, want %q", info.Status, mcp.StatusIdle)
		}

		if _, err := mgr.RunSession(ctx, "implementer", "List your commands"); err != nil {
			t.Fatalf("RunSession: %v", err)
		}

		output := pollOutput(t, mgr, "implementer", 2*time.Minute)
		if output == "" {
			t.Error("expected non-empty output from codex implementer session")
		}
	})

	t.Run("claude reviewer session", func(t *testing.T) {
		ctx := context.Background()

		info, _, err := mgr.StartSession(ctx, "reviewer", "review", "claude")
		if err != nil {
			t.Fatalf("StartSession: %v", err)
		}
		if info.Status != mcp.StatusIdle {
			t.Fatalf("status after Start = %q, want %q", info.Status, mcp.StatusIdle)
		}

		if _, err := mgr.RunSession(ctx, "reviewer", "what is 1+1?"); err != nil {
			t.Fatalf("RunSession: %v", err)
		}

		output := pollOutput(t, mgr, "reviewer", 2*time.Minute)
		if output == "" {
			t.Error("expected non-empty output from claude reviewer session")
		}
	})
}

// pollOutput polls GetOutput until the session stops running or the deadline is
// exceeded. It returns the full output collected.
func pollOutput(t *testing.T, mgr *mcp.SessionManager, name string, timeout time.Duration) string {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		chunk, total, running, err := mgr.GetOutput(name, 0)
		if err != nil {
			t.Fatalf("GetOutput(%q): %v", name, err)
		}
		if !running {
			_ = total
			return chunk
		}
		time.Sleep(2 * time.Second)
	}
	t.Fatalf("timed out after %v waiting for session %q to finish", timeout, name)
	return ""
}
