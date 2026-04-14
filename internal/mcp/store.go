package mcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const defaultSessionsPath = ".ai/sessions.json"

// Store persists session metadata to a JSON file.
// All methods are safe for concurrent use.
type Store struct {
	path string
	mu   sync.Mutex
}

func NewStore(path string) *Store {
	if path == "" {
		path = defaultSessionsPath
	}
	return &Store{path: path}
}

// Load reads the store file and returns all sessions.
// A missing file returns an empty map without error.
// A corrupt file returns an error.
func (s *Store) Load() (map[string]*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.load()
}

// Put writes or replaces a session by name.
func (s *Store) Put(session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessions, err := s.load()
	if err != nil {
		return err
	}
	sessions[session.Name] = session
	return s.save(sessions)
}

// Get returns a session by name, or an error if not found.
func (s *Store) Get(name string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessions, err := s.load()
	if err != nil {
		return nil, err
	}
	session, ok := sessions[name]
	if !ok {
		return nil, fmt.Errorf("session %q not found", name)
	}
	return session, nil
}

// Delete removes a session by name. No-op if not found.
func (s *Store) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessions, err := s.load()
	if err != nil {
		return err
	}
	delete(sessions, name)
	return s.save(sessions)
}

// List returns all sessions in undefined order.
func (s *Store) List() ([]*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessions, err := s.load()
	if err != nil {
		return nil, err
	}
	out := make([]*Session, 0, len(sessions))
	for _, session := range sessions {
		out = append(out, session)
	}
	return out, nil
}

func (s *Store) load() (map[string]*Session, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return make(map[string]*Session), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read sessions file: %w", err)
	}

	var sessions map[string]*Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, fmt.Errorf("parse sessions file: %w", err)
	}
	if sessions == nil {
		sessions = make(map[string]*Session)
	}
	return sessions, nil
}

func (s *Store) save(sessions map[string]*Session) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create sessions directory: %w", err)
	}
	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal sessions: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("write sessions file: %w", err)
	}
	return nil
}
