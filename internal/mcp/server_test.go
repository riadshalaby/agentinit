package mcp

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	mcpproto "github.com/mark3labs/mcp-go/mcp"
)

func TestNewServerRespondsToInitialize(t *testing.T) {
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

func TestNewServerStartsWithNoToolsRegistered(t *testing.T) {
	srv := NewServer("1.2.3-test")

	if got := len(srv.server.ListTools()); got != 0 {
		t.Fatalf("registered tools = %d, want 0", got)
	}
}
