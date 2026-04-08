package cmd

import (
	"context"
	"testing"
)

func TestMCPCommandIsRegistered(t *testing.T) {
	for _, command := range rootCmd.Commands() {
		if command == mcpCmd {
			return
		}
	}

	t.Fatal("expected mcp command to be registered on root command")
}

func TestMCPCommandRunsServerWithRootVersion(t *testing.T) {
	originalRunMCPServer := runMCPServer
	originalVersion := rootCmd.Version
	originalContext := mcpCmd.Context()
	t.Cleanup(func() {
		runMCPServer = originalRunMCPServer
		rootCmd.Version = originalVersion
		mcpCmd.SetContext(originalContext)
	})

	rootCmd.Version = "9.9.9-test"
	mcpCmd.SetContext(context.Background())

	called := false
	runMCPServer = func(ctx context.Context, serverVersion string) error {
		called = true
		if ctx == nil {
			t.Fatal("expected command context")
		}
		if serverVersion != rootCmd.Version {
			t.Fatalf("server version = %q, want %q", serverVersion, rootCmd.Version)
		}
		return nil
	}

	if err := mcpCmd.RunE(mcpCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}
	if !called {
		t.Fatal("expected MCP server runner to be called")
	}
}
