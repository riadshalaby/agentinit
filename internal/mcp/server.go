package mcp

import (
	"context"
	"log/slog"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

const serverName = "agentinit"

var serveStdio = mcpserver.ServeStdio

type Server struct {
	server *mcpserver.MCPServer
	logger *slog.Logger
}

func NewServer(version string) *Server {
	logger, err := NewFileLogger(defaultMCPLogPath)
	if err != nil {
		logger = newDiscardLogger()
	}
	return newServer(version, logger)
}

func newServer(version string, logger *slog.Logger) *Server {
	if logger == nil {
		logger = newDiscardLogger()
	}

	srv := mcpserver.NewMCPServer(
		serverName,
		version,
		mcpserver.WithToolCapabilities(false),
	)
	registerTools(srv, logger)
	return &Server{server: srv, logger: logger}
}

func (s *Server) Run(ctx context.Context) error {
	_ = ctx
	return serveStdio(s.server)
}
