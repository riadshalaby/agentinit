package mcp

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

const serverName = "agentinit"

var serveStdio = mcpserver.ServeStdio

type Server struct {
	server  *mcpserver.MCPServer
	manager *SessionManager
	config  Config
	logger  *slog.Logger
}

func NewServer(version string) *Server {
	logger, err := NewFileLogger(defaultMCPLogPath)
	if err != nil {
		logger = newDiscardLogger()
	}

	cwd, _ := os.Getwd()
	cfg, _ := LoadConfig(cwd)
	store := NewStore(filepath.Join(cwd, defaultSessionsPath))
	adapters := map[string]Adapter{
		"claude": NewClaudeAdapter(cwd, cfg.Defaults.Claude),
		"codex":  NewCodexAdapter(cwd, cfg.Defaults.Codex),
	}
	manager := NewSessionManager(store, adapters, cfg, cwd, logger)

	return newServer(version, manager, cfg, logger)
}

func newServer(version string, manager *SessionManager, cfg Config, logger *slog.Logger) *Server {
	if logger == nil {
		logger = newDiscardLogger()
	}

	srv := mcpserver.NewMCPServer(
		serverName,
		version,
		mcpserver.WithToolCapabilities(false),
	)
	registerTools(srv, manager, cfg, logger)
	return &Server{server: srv, manager: manager, config: cfg, logger: logger}
}

func (s *Server) Run(ctx context.Context) error {
	_ = ctx
	return serveStdio(s.server)
}
