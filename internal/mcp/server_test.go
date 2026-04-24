package mcp

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
)

func TestNewServerRespondsToInitialize(t *testing.T) {
	t.Chdir(t.TempDir())

	srv := NewServer(context.Background(), "1.2.3-test")

	mcpClient, err := client.NewInProcessClient(srv.server)
	if err != nil {
		t.Fatalf("NewInProcessClient() error = %v", err)
	}
	t.Cleanup(func() {
		if closeErr := mcpClient.Close(); closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
	})

	ctx := context.Background()
	if err := mcpClient.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	result, err := mcpClient.Initialize(ctx, mcpproto.InitializeRequest{
		Params: mcpproto.InitializeParams{
			ProtocolVersion: mcpproto.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcpproto.Implementation{
				Name:    "agentinit-test-client",
				Version: "0.0.1",
			},
		},
	})
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if result.ServerInfo.Name != serverName {
		t.Fatalf("server name = %q, want %q", result.ServerInfo.Name, serverName)
	}
	if result.ServerInfo.Version != "1.2.3-test" {
		t.Fatalf("server version = %q, want %q", result.ServerInfo.Version, "1.2.3-test")
	}
}

func TestNewServerRegistersSessionTools(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	srv := NewServer(context.Background(), "1.2.3-test")

	if got := len(srv.server.ListTools()); got != 10 {
		t.Fatalf("registered tools = %d, want 10", got)
	}
	if _, err := os.Stat(filepath.Join(tempDir, defaultMCPLogPath)); err != nil {
		t.Fatalf("expected log file %q to exist: %v", defaultMCPLogPath, err)
	}
}

func TestServerSessionToolsLifecycle(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSessionManager(
		context.Background(),
		NewStore(filepath.Join(tempDir, "sessions.json")),
		map[string]Adapter{
			"codex":  testToolAdapter{},
			"claude": testToolAdapter{},
		},
		Config{},
		filepath.Clean(filepath.Join("..", "..")),
		testLogger(),
	)
	srv := newServer(context.Background(), "1.2.3-test", manager, Config{}, testLogger())

	mcpClient, err := client.NewInProcessClient(srv.server)
	if err != nil {
		t.Fatalf("NewInProcessClient() error = %v", err)
	}
	t.Cleanup(func() {
		if closeErr := mcpClient.Close(); closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
	})

	ctx := context.Background()
	if err := mcpClient.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if _, err := mcpClient.Initialize(ctx, mcpproto.InitializeRequest{
		Params: mcpproto.InitializeParams{
			ProtocolVersion: mcpproto.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcpproto.Implementation{
				Name:    "agentinit-test-client",
				Version: "0.0.1",
			},
		},
	}); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	startResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_start",
			Arguments: map[string]any{
				"name":     "implementer",
				"role":     "implement",
				"provider": "codex",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_start) error = %v", err)
	}
	if startResult.IsError {
		t.Fatalf("CallTool(session_start) result = %+v", startResult)
	}
	assertStructuredToolResult(t, startResult, "WAIT_FOR_USER_START", `"name":"implementer"`, `"provider":"codex"`, `"session_id":"test-session-id"`)

	runResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_run",
			Arguments: map[string]any{
				"name":    "implementer",
				"command": "next_task T-001",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_run) error = %v", err)
	}
	if runResult.IsError {
		t.Fatalf("CallTool(session_run) result = %+v", runResult)
	}
	assertStructuredToolResult(t, runResult, "run started", `"message":"run started"`, `"status":"running"`)

	output := pollToolOutput(t, ctx, mcpClient, "implementer", 20000)
	if !strings.Contains(output, "response: next_task T-001") {
		t.Fatalf("session_get_output accumulated output = %q", output)
	}

	statusResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_status",
			Arguments: map[string]any{
				"name": "implementer",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_status) error = %v", err)
	}
	if statusResult.IsError {
		t.Fatalf("CallTool(session_status) result = %+v", statusResult)
	}
	assertStructuredToolResult(t, statusResult, "implementer", `"run_count":1`, `"status":"idle"`)

	listResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{Name: "session_list"},
	})
	if err != nil {
		t.Fatalf("CallTool(session_list) error = %v", err)
	}
	if listResult.IsError {
		t.Fatalf("CallTool(session_list) result = %+v", listResult)
	}
	assertStructuredToolResult(t, listResult, "implementer", `"name":"implementer"`)

	duplicateStartResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_start",
			Arguments: map[string]any{
				"name":     "implementer",
				"role":     "implement",
				"provider": "codex",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(duplicate session_start) error = %v", err)
	}
	if !duplicateStartResult.IsError {
		t.Fatalf("duplicate session_start should return tool error: %+v", duplicateStartResult)
	}

	resetResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_reset",
			Arguments: map[string]any{
				"name": "implementer",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_reset) error = %v", err)
	}
	if resetResult.IsError {
		t.Fatalf("CallTool(session_reset) result = %+v", resetResult)
	}

	deleteResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_delete",
			Arguments: map[string]any{
				"name": "implementer",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_delete) error = %v", err)
	}
	if deleteResult.IsError {
		t.Fatalf("CallTool(session_delete) result = %+v", deleteResult)
	}

	missingStatusResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_status",
			Arguments: map[string]any{
				"name": "implementer",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(missing session_status) error = %v", err)
	}
	if !missingStatusResult.IsError {
		t.Fatalf("missing session_status should return tool error: %+v", missingStatusResult)
	}
}

func TestServerSessionGetResultTool(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSessionManager(
		context.Background(),
		NewStore(filepath.Join(tempDir, "sessions.json")),
		map[string]Adapter{
			"codex":  testToolAdapter{},
			"claude": testToolAdapter{},
		},
		Config{},
		filepath.Clean(filepath.Join("..", "..")),
		testLogger(),
	)
	srv := newServer(context.Background(), "1.2.3-test", manager, Config{}, testLogger())

	mcpClient, err := client.NewInProcessClient(srv.server)
	if err != nil {
		t.Fatalf("NewInProcessClient() error = %v", err)
	}
	t.Cleanup(func() {
		if closeErr := mcpClient.Close(); closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
	})

	ctx := context.Background()
	if err := mcpClient.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if _, err := mcpClient.Initialize(ctx, mcpproto.InitializeRequest{
		Params: mcpproto.InitializeParams{
			ProtocolVersion: mcpproto.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcpproto.Implementation{
				Name:    "agentinit-test-client",
				Version: "0.0.1",
			},
		},
	}); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	_, err = mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_start",
			Arguments: map[string]any{
				"name":     "implementer",
				"role":     "implement",
				"provider": "codex",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_start) error = %v", err)
	}

	emptyResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_get_result",
			Arguments: map[string]any{
				"name": "implementer",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_get_result before run) error = %v", err)
	}
	if emptyResult.IsError {
		t.Fatalf("CallTool(session_get_result before run) result = %+v", emptyResult)
	}
	assertStructuredToolResult(t, emptyResult, "no completed result yet", `"message":"no completed result yet"`)

	_, err = mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_run",
			Arguments: map[string]any{
				"name":    "implementer",
				"command": "next_task T-003",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_run) error = %v", err)
	}

	waitForToolSessionStatus(t, ctx, mcpClient, "implementer", StatusIdle)

	result, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_get_result",
			Arguments: map[string]any{
				"name": "implementer",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_get_result) error = %v", err)
	}
	if result.IsError {
		t.Fatalf("CallTool(session_get_result) result = %+v", result)
	}
	assertStructuredToolResult(t, result, "Status:idle", `"status":"idle"`, `"duration_secs":`, `"exit_summary":"response: next_task T-003"`)

	resetResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_reset",
			Arguments: map[string]any{
				"name": "implementer",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_reset) error = %v", err)
	}
	if resetResult.IsError {
		t.Fatalf("CallTool(session_reset) result = %+v", resetResult)
	}

	clearedResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "session_get_result",
			Arguments: map[string]any{
				"name": "implementer",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(session_get_result after reset) error = %v", err)
	}
	if clearedResult.IsError {
		t.Fatalf("CallTool(session_get_result after reset) result = %+v", clearedResult)
	}
	assertStructuredToolResult(t, clearedResult, "no completed result yet", `"message":"no completed result yet"`)
}

type testToolAdapter struct{}

func (testToolAdapter) Start(_ context.Context, session *Session, _ StartOpts) (string, error) {
	session.ProviderState.SessionID = "test-session-id"
	return "WAIT_FOR_USER_START", nil
}

func (testToolAdapter) RunStream(_ context.Context, _ *Session, command string, _ RunOpts, w io.Writer) error {
	_, err := fmt.Fprintf(w, "response: %s", command)
	return err
}

func (testToolAdapter) Stop(_ context.Context, _ *Session) error {
	return nil
}

func testLogger() *slog.Logger {
	return newDiscardLogger()
}

func containsAll(text string, substrings ...string) bool {
	for _, substring := range substrings {
		if !strings.Contains(text, substring) {
			return false
		}
	}
	return true
}

func assertStructuredToolResult(t *testing.T, result *mcpproto.CallToolResult, textSubstring string, jsonSubstrings ...string) {
	t.Helper()

	if result.StructuredContent == nil {
		t.Fatal("tool result missing structured content")
	}
	if len(result.Content) < 2 {
		t.Fatalf("tool result content length = %d, want at least 2", len(result.Content))
	}

	text := mcpproto.GetTextFromContent(result.Content[0])
	if !containsAll(text, textSubstring) {
		t.Fatalf("tool result text = %q", text)
	}

	jsonFallback := mcpproto.GetTextFromContent(result.Content[1])
	if !containsAll(jsonFallback, jsonSubstrings...) {
		t.Fatalf("tool result JSON fallback = %q", jsonFallback)
	}
}

func pollToolOutput(t *testing.T, ctx context.Context, mcpClient *client.Client, name string, limit int) string {
	t.Helper()
	var output strings.Builder
	offset := 0
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		result, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
			Params: mcpproto.CallToolParams{
				Name: "session_get_output",
				Arguments: map[string]any{
					"name":   name,
					"offset": offset,
					"limit":  limit,
				},
			},
		})
		if err != nil {
			t.Fatalf("CallTool(session_get_output) error = %v", err)
		}
		if result.IsError {
			t.Fatalf("CallTool(session_get_output) result = %+v", result)
		}
		if result.StructuredContent == nil {
			t.Fatal("session_get_output missing structured content")
		}
		text := mcpproto.GetTextFromContent(result.Content[0])
		output.WriteString(text)
		jsonFallback := mcpproto.GetTextFromContent(result.Content[1])
		if !strings.Contains(jsonFallback, `"total_bytes":`) || !strings.Contains(jsonFallback, `"running":`) {
			t.Fatalf("session_get_output JSON fallback = %q", jsonFallback)
		}
		offset = extractTotalBytes(t, jsonFallback)
		if strings.Contains(jsonFallback, `"running":false`) {
			return output.String()
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("timed out polling session_get_output")
	return ""
}

func extractTotalBytes(t *testing.T, jsonFallback string) int {
	t.Helper()
	marker := `"total_bytes":`
	idx := strings.Index(jsonFallback, marker)
	if idx == -1 {
		t.Fatalf("total_bytes missing from %q", jsonFallback)
	}
	rest := jsonFallback[idx+len(marker):]
	end := strings.IndexAny(rest, ",}")
	if end == -1 {
		t.Fatalf("total_bytes value not terminated in %q", jsonFallback)
	}
	var total int
	if _, err := fmt.Sscanf(rest[:end], "%d", &total); err != nil {
		t.Fatalf("parse total_bytes from %q: %v", jsonFallback, err)
	}
	return total
}

func waitForToolSessionStatus(t *testing.T, ctx context.Context, mcpClient *client.Client, name string, want SessionStatus) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		result, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
			Params: mcpproto.CallToolParams{
				Name: "session_status",
				Arguments: map[string]any{
					"name": name,
				},
			},
		})
		if err != nil {
			t.Fatalf("CallTool(session_status) error = %v", err)
		}
		if result.IsError {
			t.Fatalf("CallTool(session_status) result = %+v", result)
		}
		jsonFallback := mcpproto.GetTextFromContent(result.Content[1])
		if strings.Contains(jsonFallback, fmt.Sprintf(`"status":"%s"`, want)) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for session_status %q", want)
}
