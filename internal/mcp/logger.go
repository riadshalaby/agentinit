package mcp

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

const defaultMCPLogPath = ".ai/mcp-server.log"

func NewFileLogger(path string) (*slog.Logger, error) {
	if path == "" {
		path = defaultMCPLogPath
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create log directory: %w", err)
		}
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	return slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug})), nil
}

func newDiscardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
