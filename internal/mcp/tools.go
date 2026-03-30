package mcp

import (
	"context"
	"fmt"

	mcpproto "github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type startSessionArgs struct {
	Role  string `json:"role"`
	Agent string `json:"agent"`
}

type stopSessionArgs struct {
	Role string `json:"role"`
}

type sendCommandArgs struct {
	Role    string `json:"role"`
	Command string `json:"command"`
}

type listSessionsArgs struct{}

func registerTools(server *mcpserver.MCPServer, manager *SessionManager) {
	server.AddTool(
		mcpproto.NewTool(
			"start_session",
			mcpproto.WithDescription("Start an agent session for a workflow role."),
			mcpproto.WithString("role", mcpproto.Required(), mcpproto.Description("Role to start: plan, implement, or review.")),
			mcpproto.WithString("agent", mcpproto.Required(), mcpproto.Description("Agent backend: claude or codex.")),
		),
		mcpproto.NewTypedToolHandler(func(ctx context.Context, _ mcpproto.CallToolRequest, args startSessionArgs) (*mcpproto.CallToolResult, error) {
			info, err := manager.StartSession(ctx, args.Role, args.Agent)
			if err != nil {
				return mcpproto.NewToolResultErrorf("start_session failed: %v", err), nil
			}
			return jsonResult(info, fmt.Sprintf("started %s session for %s", args.Agent, args.Role))
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"stop_session",
			mcpproto.WithDescription("Stop a running agent session for a workflow role."),
			mcpproto.WithString("role", mcpproto.Required(), mcpproto.Description("Role to stop: plan, implement, or review.")),
		),
		mcpproto.NewTypedToolHandler(func(_ context.Context, _ mcpproto.CallToolRequest, args stopSessionArgs) (*mcpproto.CallToolResult, error) {
			info, err := manager.StopSession(args.Role)
			if err != nil {
				return mcpproto.NewToolResultErrorf("stop_session failed: %v", err), nil
			}
			return jsonResult(info, fmt.Sprintf("stopped session for %s", args.Role))
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"send_command",
			mcpproto.WithDescription("Send a command to a running agent session and return the session output."),
			mcpproto.WithString("role", mcpproto.Required(), mcpproto.Description("Role to target: plan, implement, or review.")),
			mcpproto.WithString("command", mcpproto.Required(), mcpproto.Description("Command text to write to the session stdin.")),
		),
		mcpproto.NewTypedToolHandler(func(ctx context.Context, _ mcpproto.CallToolRequest, args sendCommandArgs) (*mcpproto.CallToolResult, error) {
			result, err := manager.SendCommand(ctx, args.Role, args.Command)
			if err != nil {
				return mcpproto.NewToolResultErrorf("send_command failed: %v", err), nil
			}
			return jsonResult(result, result.Output)
		}),
	)

	server.AddTool(
		mcpproto.NewTool(
			"list_sessions",
			mcpproto.WithDescription("List all tracked agent sessions."),
		),
		mcpproto.NewTypedToolHandler(func(ctx context.Context, _ mcpproto.CallToolRequest, args listSessionsArgs) (*mcpproto.CallToolResult, error) {
			_ = ctx
			list := manager.ListSessions()
			return jsonResult(list, fmt.Sprintf("%+v", list))
		}),
	)
}

func jsonResult(data any, fallbackText string) (*mcpproto.CallToolResult, error) {
	result, err := mcpproto.NewToolResultJSON(data)
	if err != nil {
		return nil, err
	}
	result.Content = []mcpproto.Content{mcpproto.NewTextContent(fallbackText)}
	return result, nil
}
