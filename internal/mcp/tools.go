package mcp

import (
	"context"
	"log/slog"

	mcpproto "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type sessionStartArgs struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Provider string `json:"provider"`
}

type sessionRunArgs struct {
	Name           string `json:"name"`
	Command        string `json:"command"`
	TimeoutSeconds int    `json:"timeout_seconds"`
}

type sessionNameArgs struct {
	Name string `json:"name"`
}

func registerTools(server *mcpserver.MCPServer, logger *slog.Logger) {
	server.AddTool(
		mcpproto.NewTool(
			"session_start",
			mcpproto.WithDescription("Create and initialize a named session."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Unique session name.")),
			mcpproto.WithString("role", mcpproto.Required(), mcpproto.Description("Role to start: implement or review.")),
			mcpproto.WithString("provider", mcpproto.Description("Provider backend: claude or codex.")),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, args sessionStartArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_start", "args", args)
			return mcpproto.NewToolResultErrorf("not implemented"), nil
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"session_run",
			mcpproto.WithDescription("Send a command to a named session and return full output."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Session name.")),
			mcpproto.WithString("command", mcpproto.Required(), mcpproto.Description("Command to execute.")),
			mcpproto.WithNumber("timeout_seconds", mcpproto.Description("Timeout in seconds. Defaults to 300.")),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, args sessionRunArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_run", "args", args)
			return mcpproto.NewToolResultErrorf("not implemented"), nil
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"session_status",
			mcpproto.WithDescription("Get the current status for a named session."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Session name.")),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, args sessionNameArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_status", "args", args)
			return mcpproto.NewToolResultErrorf("not implemented"), nil
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"session_list",
			mcpproto.WithDescription("List all tracked sessions."),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, _ struct{}) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_list")
			return mcpproto.NewToolResultErrorf("not implemented"), nil
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"session_stop",
			mcpproto.WithDescription("Stop an in-flight run for a named session."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Session name.")),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, args sessionNameArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_stop", "args", args)
			return mcpproto.NewToolResultErrorf("not implemented"), nil
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"session_reset",
			mcpproto.WithDescription("Reset provider-specific state for a named session."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Session name.")),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, args sessionNameArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_reset", "args", args)
			return mcpproto.NewToolResultErrorf("not implemented"), nil
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"session_delete",
			mcpproto.WithDescription("Delete a named session."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Session name.")),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, args sessionNameArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_delete", "args", args)
			return mcpproto.NewToolResultErrorf("not implemented"), nil
		}),
	)
}

func jsonResult(data any, fallbackText string) (*mcpproto.CallToolResult, error) {
	result, err := mcpproto.NewToolResultJSON(data)
	if err != nil {
		return nil, err
	}
	result.Content = append([]mcpproto.Content{mcpproto.NewTextContent(fallbackText)}, result.Content...)
	return result, nil
}
