package mcp

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

const serverName = "aide"

var serveStdio = mcpserver.ServeStdio

type Server struct {
	ctx     context.Context
	server  *mcpserver.MCPServer
	manager *SessionManager
	config  Config
	logger  *slog.Logger
}

func NewServer(ctx context.Context, version string) *Server {
	if ctx == nil {
		ctx = context.Background()
	}
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
	manager := NewSessionManager(ctx, store, adapters, cfg, cwd, logger)

	return newServer(ctx, version, manager, cfg, logger)
}

func newServer(ctx context.Context, version string, manager *SessionManager, cfg Config, logger *slog.Logger) *Server {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		logger = newDiscardLogger()
	}

	srv := mcpserver.NewMCPServer(
		serverName,
		version,
		mcpserver.WithToolCapabilities(false),
	)
	registerTools(srv, manager, cfg, logger)
	return &Server{ctx: ctx, server: srv, manager: manager, config: cfg, logger: logger}
}

func (s *Server) Run(ctx context.Context) error {
	if ctx != nil {
		s.ctx = ctx
		if s.manager != nil {
			s.manager.ctx = ctx
		}
	}
	return serveStdio(s.server)
}
