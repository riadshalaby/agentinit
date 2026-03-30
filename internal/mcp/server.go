package mcp

import (
	"context"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

const serverName = "agentinit"

var serveStdio = mcpserver.ServeStdio

type Server struct {
	server *mcpserver.MCPServer
}

func NewServer(version string) *Server {
	return &Server{
		server: mcpserver.NewMCPServer(
			serverName,
			version,
			mcpserver.WithToolCapabilities(false),
		),
	}
}

func (s *Server) Run(ctx context.Context) error {
	_ = ctx
	return serveStdio(s.server)
}
