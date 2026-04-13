package mcp

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
)

func TestNewServerRespondsToInitialize(t *testing.T) {
	t.Chdir(t.TempDir())

	srv := NewServer("1.2.3-test")

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

	srv := NewServer("1.2.3-test")

	if got := len(srv.server.ListTools()); got != 5 {
		t.Fatalf("registered tools = %d, want 5", got)
	}
	if _, err := os.Stat(filepath.Join(tempDir, defaultMCPLogPath)); err != nil {
		t.Fatalf("expected log file %q to exist: %v", defaultMCPLogPath, err)
	}
}

func TestServerSessionToolsLifecycle(t *testing.T) {
	srv := newServer("1.2.3-test", newSessionManager(testLauncher(t), testLogger()), testLogger())

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

	tools, err := mcpClient.ListTools(ctx, mcpproto.ListToolsRequest{})
	if err != nil {
		t.Fatalf("ListTools() error = %v", err)
	}
	if len(tools.Tools) != 5 {
		t.Fatalf("ListTools() count = %d, want 5", len(tools.Tools))
	}

	startResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "start_session",
			Arguments: map[string]any{
				"role":  "implement",
				"agent": "codex",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(start_session) error = %v", err)
	}
	if startResult.IsError {
		t.Fatalf("CallTool(start_session) result = %+v", startResult)
	}

	sendResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "send_command",
			Arguments: map[string]any{
				"role":    "implement",
				"command": "next_task T-003",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(send_command) error = %v", err)
	}
	if sendResult.IsError {
		t.Fatalf("CallTool(send_command) result = %+v", sendResult)
	}
	output := mcpproto.GetTextFromContent(sendResult.Content[0])
	if !containsAll(output, "sent command to implement") {
		t.Fatalf("CallTool(send_command) output = %q", output)
	}

	getOutputResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "get_output",
			Arguments: map[string]any{
				"role":            "implement",
				"timeout_seconds": 1,
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(get_output) error = %v", err)
	}
	if getOutputResult.IsError {
		t.Fatalf("CallTool(get_output) result = %+v", getOutputResult)
	}
	getOutputText := mcpproto.GetTextFromContent(getOutputResult.Content[0])
	if !containsAll(getOutputText, "response:next_task T-003") {
		t.Fatalf("CallTool(get_output) output = %q", getOutputText)
	}

	listResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{Name: "list_sessions"},
	})
	if err != nil {
		t.Fatalf("CallTool(list_sessions) error = %v", err)
	}
	if listResult.IsError {
		t.Fatalf("CallTool(list_sessions) result = %+v", listResult)
	}
	if got := mcpproto.GetTextFromContent(listResult.Content[0]); !containsAll(got, "implement", "codex") {
		t.Fatalf("CallTool(list_sessions) output = %q", got)
	}

	duplicateResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "start_session",
			Arguments: map[string]any{
				"role":  "implement",
				"agent": "claude",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(duplicate start_session) error = %v", err)
	}
	if !duplicateResult.IsError {
		t.Fatalf("duplicate start_session should return tool error: %+v", duplicateResult)
	}

	stopResult, err := mcpClient.CallTool(ctx, mcpproto.CallToolRequest{
		Params: mcpproto.CallToolParams{
			Name: "stop_session",
			Arguments: map[string]any{
				"role": "implement",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool(stop_session) error = %v", err)
	}
	if stopResult.IsError {
		t.Fatalf("CallTool(stop_session) result = %+v", stopResult)
	}
}

func containsAll(text string, substrings ...string) bool {
	for _, substring := range substrings {
		if !strings.Contains(text, substring) {
			return false
		}
	}
	return true
}
