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
	Name    string `json:"name"`
	Command string `json:"command"`
}

type sessionWaitArgs struct {
	Name           string `json:"name"`
	TimeoutSeconds int    `json:"timeout_seconds"`
}

type sessionGetOutputArgs struct {
	Name   string `json:"name"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type sessionNameArgs struct {
	Name string `json:"name"`
}

func registerTools(server *mcpserver.MCPServer, manager *SessionManager, cfg Config, logger *slog.Logger) {
	const defaultSessionOutputLimit = 20000

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
			mcpproto.WithDescription("Send a command to a named session. Returns immediately; use session_wait for the structured completion result."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Session name.")),
			mcpproto.WithString("command", mcpproto.Required(), mcpproto.Description("Command to execute.")),
		),
		mcpproto.NewTypedToolHandler(func(ctx context.Context, _ mcpproto.CallToolRequest, args sessionRunArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_run", "args", args)
			info, err := manager.RunSession(ctx, args.Name, args.Command)
			if err != nil {
				logger.Error("tool call failed", "tool", "session_run", "args", args, "error", err)
				return mcpproto.NewToolResultErrorf("session_run failed: %v", err), nil
			}
			logger.Info("tool call completed", "tool", "session_run", "name", info.Name, "status", info.Status, "run_count", info.RunCount)
			return jsonResult(struct {
				Session SessionInfo `json:"session"`
				Message string      `json:"message"`
			}{Session: info, Message: "run started"}, "run started")
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"session_wait",
			mcpproto.WithDescription("Wait for a named session run to finish and return the structured outcome."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Session name.")),
			mcpproto.WithNumber("timeout_seconds", mcpproto.Description("Optional wait timeout in seconds.")),
		),
		mcpproto.NewTypedToolHandler(func(ctx context.Context, _ mcpproto.CallToolRequest, args sessionWaitArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_wait", "args", args)
			waitCtx := ctx
			cancel := func() {}
			if args.TimeoutSeconds > 0 {
				waitCtx, cancel = context.WithTimeout(ctx, time.Duration(args.TimeoutSeconds)*time.Second)
			}
			defer cancel()

			info, result, err := manager.WaitSession(waitCtx, args.Name)
			response := WaitResult{Session: info, Result: result}
			if err != nil {
				if info.Name == "" {
					logger.Error("tool call failed", "tool", "session_wait", "args", args, "error", err)
					return mcpproto.NewToolResultErrorf("session_wait failed: %v", err), nil
				}
				response.Error = err.Error()
				logger.Warn("tool call completed with wait error", "tool", "session_wait", "name", info.Name, "status", info.Status, "error", err)
				return jsonResult(response, err.Error())
			}
			logger.Info("tool call completed", "tool", "session_wait", "name", info.Name, "status", info.Status)
			return jsonResult(response, fmt.Sprintf("%+v", response))
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"session_get_output",
			mcpproto.WithDescription("Poll output from a running or completed session. Pass offset=0 to read from the start, or offset=total_bytes from the previous call to read only new output. Responses are capped by the optional limit parameter."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Session name.")),
			mcpproto.WithNumber("offset", mcpproto.Description("Byte offset to start reading from.")),
			mcpproto.WithNumber("limit", mcpproto.Description("Maximum bytes to return. Default: 20000. Omit or pass 0 to use the default.")),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, args sessionGetOutputArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_get_output", "args", args)
			limit := args.Limit
			if limit == 0 {
				limit = defaultSessionOutputLimit
			}
			chunk, totalBytes, running, err := manager.GetOutput(args.Name, args.Offset, limit)
			if err != nil {
				logger.Error("tool call failed", "tool", "session_get_output", "args", args, "error", err)
				return mcpproto.NewToolResultErrorf("session_get_output failed: %v", err), nil
			}
			info, err := manager.GetSession(args.Name)
			if err != nil {
				logger.Error("tool call failed", "tool", "session_get_output", "args", args, "error", err)
				return mcpproto.NewToolResultErrorf("session_get_output failed: %v", err), nil
			}
			logger.Info("tool call completed", "tool", "session_get_output", "name", info.Name, "status", info.Status, "running", running, "total_bytes", totalBytes)
			return jsonResult(struct {
				Chunk      string        `json:"chunk"`
				TotalBytes int           `json:"total_bytes"`
				Running    bool          `json:"running"`
				Status     SessionStatus `json:"status"`
			}{Chunk: chunk, TotalBytes: totalBytes, Running: running, Status: info.Status}, chunk)
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"session_get_result",
			mcpproto.WithDescription("Get the structured result for the most recent completed run of a named session."),
			mcpproto.WithString("name", mcpproto.Required(), mcpproto.Description("Session name.")),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, args sessionNameArgs) (*mcpproto.CallToolResult, error) {
			logger.Info("tool call started", "tool", "session_get_result", "args", args)
			result, err := manager.GetResult(args.Name)
			if err != nil {
				logger.Error("tool call failed", "tool", "session_get_result", "args", args, "error", err)
				return mcpproto.NewToolResultErrorf("session_get_result failed: %v", err), nil
			}
			if result == nil {
				message := "no completed result yet"
				logger.Info("tool call completed", "tool", "session_get_result", "name", args.Name, "message", message)
				return jsonResult(map[string]any{"message": message}, message)
			}
			logger.Info("tool call completed", "tool", "session_get_result", "name", args.Name, "status", result.Status)
			return jsonResult(result, fmt.Sprintf("%+v", *result))
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
