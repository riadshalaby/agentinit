package mcp

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestStoreLoadMissingFile(t *testing.T) {
	t.Parallel()

	store := NewStore(filepath.Join(t.TempDir(), "sessions.json"))
	sessions, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("Load() len = %d, want 0", len(sessions))
	}
}

func TestStorePutGetRoundTrip(t *testing.T) {
	t.Parallel()

	store := NewStore(filepath.Join(t.TempDir(), "sessions.json"))
	want := testSession("alpha")

	if err := store.Put(want); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	got, err := store.Get("alpha")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Name != want.Name || got.Provider != want.Provider || got.RunCount != want.RunCount {
		t.Fatalf("Get() = %+v, want %+v", got, want)
	}
}

func TestStorePutListContainsSession(t *testing.T) {
	t.Parallel()

	store := NewStore(filepath.Join(t.TempDir(), "sessions.json"))
	if err := store.Put(testSession("alpha")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	sessions, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(sessions) != 1 || sessions[0].Name != "alpha" {
		t.Fatalf("List() = %+v, want one alpha session", sessions)
	}
}

func TestStoreDeleteRemovesSession(t *testing.T) {
	t.Parallel()

	store := NewStore(filepath.Join(t.TempDir(), "sessions.json"))
	if err := store.Put(testSession("alpha")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if err := store.Delete("alpha"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := store.Get("alpha")
	if err == nil {
		t.Fatal("Get() error = nil, want not found")
	}
}

func TestStoreListMultipleSessions(t *testing.T) {
	t.Parallel()

	store := NewStore(filepath.Join(t.TempDir(), "sessions.json"))
	if err := store.Put(testSession("alpha")); err != nil {
		t.Fatalf("Put(alpha) error = %v", err)
	}
	if err := store.Put(testSession("beta")); err != nil {
		t.Fatalf("Put(beta) error = %v", err)
	}

	sessions, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(sessions) != 2 {
		t.Fatalf("List() len = %d, want 2", len(sessions))
	}
}

func TestStoreLoadCorruptJSON(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "sessions.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := NewStore(path)
	_, err := store.Load()
	if err == nil {
		t.Fatal("Load() error = nil, want parse error")
	}
}

func TestStorePutCreatesParentDirectory(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "x", "y", "sessions.json")
	store := NewStore(path)
	if err := store.Put(testSession("alpha")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
}

func TestStoreGetMissingReturnsError(t *testing.T) {
	t.Parallel()

	store := NewStore(filepath.Join(t.TempDir(), "sessions.json"))
	_, err := store.Get("missing")
	if err == nil {
		t.Fatal("Get() error = nil, want not found")
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Get() error = %v, want store not-found error", err)
	}
}

func testSession(name string) *Session {
	return &Session{
		Name:     name,
		Role:     "implement",
		Provider: "codex",
		Status:   StatusIdle,
		RunCount: 1,
	}
}
