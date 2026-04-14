package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

func registerTools(server *mcpserver.MCPServer, manager *SessionManager, cfg Config, logger *slog.Logger) {
	server.AddTool(
		mcpproto.NewTool(
			"session_start",
			mcpproto.WithDescription("Create and initialize a named session."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Unique session name.")),
			mcpproto.WithString("role", mcpproto.Required(), mcpproto.Description("Role to start: implement or review.")),
			mcpproto.WithString("provider", mcpproto.Description("Provider backend: claude or codex.")),
		),
		mcpproto.NewTypedToolHandler(func(ctx context.Context, _ mcpproto.CallToolRequest, args sessionStartArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_start", "args", args)
			provider := args.Provider
			if provider == "" {
				provider = cfg.ProviderForRole(args.Role)
			}
			info, output, err := manager.StartSession(ctx, args.Name, args.Role, provider)
			if err != nil {
				logger.Error("tool call failed", "tool", "session_start", "args", args, "error", err)
				return mcpproto.NewToolResultErrorf("session_start failed: %v", err), nil
			}
			logger.Info("tool call completed", "tool", "session_start", "name", info.Name, "role", info.Role, "provider", info.Provider, "status", info.Status)
			session, sessionErr := manager.store.Get(args.Name)
			if sessionErr != nil {
				logger.Error("tool call failed", "tool", "session_start", "args", args, "error", sessionErr)
				return mcpproto.NewToolResultErrorf("session_start failed: %v", sessionErr), nil
			}
			return jsonResult(struct {
				Session   SessionInfo `json:"session"`
				SessionID string      `json:"session_id,omitempty"`
				Output    string      `json:"output"`
			}{Session: info, SessionID: session.ProviderState.SessionID, Output: output}, output)
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
		mcpproto.NewTypedToolHandler(func(ctx context.Context, _ mcpproto.CallToolRequest, args sessionRunArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_run", "args", args)
			timeoutSeconds := args.TimeoutSeconds
			if timeoutSeconds <= 0 {
				timeoutSeconds = 300
			}
			info, output, err := manager.RunSession(ctx, args.Name, args.Command, time.Duration(timeoutSeconds)*time.Second)
			if err != nil {
				logger.Error("tool call failed", "tool", "session_run", "args", args, "error", err)
				return mcpproto.NewToolResultErrorf("session_run failed: %v", err), nil
			}
			logger.Info("tool call completed", "tool", "session_run", "name", info.Name, "status", info.Status, "run_count", info.RunCount)
			session, sessionErr := manager.store.Get(args.Name)
			if sessionErr != nil {
				logger.Error("tool call failed", "tool", "session_run", "args", args, "error", sessionErr)
				return mcpproto.NewToolResultErrorf("session_run failed: %v", sessionErr), nil
			}
			return jsonResult(struct {
				Session   SessionInfo `json:"session"`
				SessionID string      `json:"session_id,omitempty"`
				Output    string      `json:"output"`
			}{Session: info, SessionID: session.ProviderState.SessionID, Output: output}, output)
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
			info, err := manager.GetSession(args.Name)
			if err != nil {
				logger.Error("tool call failed", "tool", "session_status", "args", args, "error", err)
				return mcpproto.NewToolResultErrorf("session_status failed: %v", err), nil
			}
			logger.Info("tool call completed", "tool", "session_status", "name", info.Name, "status", info.Status)
			return jsonResult(info, fmt.Sprintf("%+v", info))
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"session_list",
			mcpproto.WithDescription("List all tracked sessions."),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, _ struct{}) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_list")
			list, err := manager.ListSessions()
			if err != nil {
				logger.Error("tool call failed", "tool", "session_list", "error", err)
				return mcpproto.NewToolResultErrorf("session_list failed: %v", err), nil
			}
			logger.Info("tool call completed", "tool", "session_list", "count", len(list))
			return jsonResult(list, fmt.Sprintf("%+v", list))
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
			info, err := manager.StopSession(args.Name)
			if err != nil {
				logger.Error("tool call failed", "tool", "session_stop", "args", args, "error", err)
				return mcpproto.NewToolResultErrorf("session_stop failed: %v", err), nil
			}
			logger.Info("tool call completed", "tool", "session_stop", "name", info.Name, "status", info.Status)
			return jsonResult(info, fmt.Sprintf("%+v", info))
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
			info, err := manager.ResetSession(args.Name)
			if err != nil {
				logger.Error("tool call failed", "tool", "session_reset", "args", args, "error", err)
				return mcpproto.NewToolResultErrorf("session_reset failed: %v", err), nil
			}
			logger.Info("tool call completed", "tool", "session_reset", "name", info.Name, "status", info.Status)
			return jsonResult(info, fmt.Sprintf("%+v", info))
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
			if err := manager.DeleteSession(args.Name); err != nil {
				logger.Error("tool call failed", "tool", "session_delete", "args", args, "error", err)
				return mcpproto.NewToolResultErrorf("session_delete failed: %v", err), nil
			}
			logger.Info("tool call completed", "tool", "session_delete", "name", args.Name)
			return jsonResult(map[string]any{"name": args.Name, "deleted": true}, fmt.Sprintf("deleted %s", args.Name))
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
