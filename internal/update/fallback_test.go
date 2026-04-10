package update

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInferProjectTypePrefersGoThenNodeThenJava(t *testing.T) {
	dir := t.TempDir()
	if got := InferProjectType(dir); got != "" {
		t.Fatalf("InferProjectType() = %q, want empty", got)
	}

	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := InferProjectType(dir); got != "node" {
		t.Fatalf("InferProjectType() = %q, want %q", got, "node")
	}

	if err := os.WriteFile(filepath.Join(dir, "pom.xml"), []byte("<project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := InferProjectType(dir); got != "node" {
		t.Fatalf("InferProjectType() = %q, want %q", got, "node")
	}

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module demo"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := InferProjectType(dir); got != "go" {
		t.Fatalf("InferProjectType() = %q, want %q", got, "go")
	}
}

func TestInferProjectTypeFallsBackToAgentsValidationCommands(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("Format: `go fmt ./...`\nVet: `go vet ./...`\nTest: `go test ./...`\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if got := InferProjectType(dir); got != "go" {
		t.Fatalf("InferProjectType() = %q, want %q", got, "go")
	}
}

func TestDiscoverManagedFilesFindsKnownPaths(t *testing.T) {
	dir := t.TempDir()
	for _, relPath := range []string{"AGENTS.md", ".gitignore", "scripts/ai-launch.sh"} {
		absPath := filepath.Join(dir, relPath)
		if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(absPath, []byte("content"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	manifest := DiscoverManagedFiles(dir)
	if len(manifest.Files) != 3 {
		t.Fatalf("len(Files) = %d, want 3", len(manifest.Files))
	}

	got := map[string]string{}
	for _, file := range manifest.Files {
		got[file.Path] = file.Management
	}
	if got["AGENTS.md"] != "marker" {
		t.Fatalf("AGENTS.md management = %q, want %q", got["AGENTS.md"], "marker")
	}
	if got[".gitignore"] != "full" {
		t.Fatalf(".gitignore management = %q, want %q", got[".gitignore"], "full")
	}
}
