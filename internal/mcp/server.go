package mcp

import (
	"context"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

const serverName = "agentinit"

var serveStdio = mcpserver.ServeStdio

type Server struct {
	server  *mcpserver.MCPServer
	manager *SessionManager
}

func NewServer(version string) *Server {
	return newServer(version, NewSessionManager())
}

func newServer(version string, manager *SessionManager) *Server {
	srv := mcpserver.NewMCPServer(
		serverName,
		version,
		mcpserver.WithToolCapabilities(false),
	)
	registerTools(srv, manager)

	return &Server{
		server:  srv,
		manager: manager,
	}
}

func (s *Server) Run(ctx context.Context) error {
	_ = ctx
	return serveStdio(s.server)
}
