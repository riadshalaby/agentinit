package mcp

import (
	"context"
	"log/slog"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

const serverName = "agentinit"

var serveStdio = mcpserver.ServeStdio

type Server struct {
	server  *mcpserver.MCPServer
	manager *SessionManager
	logger  *slog.Logger
}

func NewServer(version string) *Server {
	logger, err := NewFileLogger(defaultMCPLogPath)
	if err != nil {
		logger = newDiscardLogger()
	}

	return newServer(version, NewSessionManager(logger), logger)
}

func newServer(version string, manager *SessionManager, logger *slog.Logger) *Server {
	if logger == nil {
		logger = newDiscardLogger()
	}

	srv := mcpserver.NewMCPServer(
		serverName,
		version,
		mcpserver.WithToolCapabilities(false),
	)
	registerTools(srv, manager, logger)

	return &Server{
		server:  srv,
		manager: manager,
		logger:  logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	_ = ctx
	return serveStdio(s.server)
}
