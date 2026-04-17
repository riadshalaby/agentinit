package scaffold

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestGenerateManifestClassifiesManagedFiles(t *testing.T) {
	originalNow := manifestNow
	manifestNow = func() time.Time {
		return time.Date(2026, 4, 10, 12, 34, 56, 0, time.UTC)
	}
	defer func() {
		manifestNow = originalNow
	}()

	manifest := GenerateManifest(map[string]string{
		"AGENTS.md":                  "agents",
		"README.md":                  "readme",
		"CLAUDE.md":                  "claude",
		"ROADMAP.md":                 "roadmap",
		".ai/config.json":            "config",
		"ROADMAP.template.md":        "template",
		".ai/PLAN.template.md":       "plan template",
		".ai/prompts/implementer.md": "prompt",
		".gitignore":                 "ignore",
	}, "v1.2.3")

	if manifest.Version != "v1.2.3" {
		t.Fatalf("Version = %q, want %q", manifest.Version, "v1.2.3")
	}
	if manifest.GeneratedAt != "2026-04-10T12:34:56Z" {
		t.Fatalf("GeneratedAt = %q, want %q", manifest.GeneratedAt, "2026-04-10T12:34:56Z")
	}

	want := []ManifestFile{
		{Path: ".ai/prompts/implementer.md", Management: managementFull},
		{Path: ".gitignore", Management: managementFull},
		{Path: "AGENTS.md", Management: managementMarker},
		{Path: "ROADMAP.template.md", Management: managementFull},
	}
	if !reflect.DeepEqual(manifest.Files, want) {
		t.Fatalf("Files = %#v, want %#v", manifest.Files, want)
	}
}

func TestWriteReadManifestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	manifest := Manifest{
		Version:     "v9.9.9",
		GeneratedAt: "2026-04-10T00:00:00Z",
		Files: []ManifestFile{
			{Path: "AGENTS.md", Management: managementMarker},
		},
	}

	if err := WriteManifest(dir, manifest); err != nil {
		t.Fatalf("WriteManifest() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, manifestRelPath)); err != nil {
		t.Fatalf("manifest file missing: %v", err)
	}

	got, err := ReadManifest(dir)
	if err != nil {
		t.Fatalf("ReadManifest() error: %v", err)
	}

	if !reflect.DeepEqual(got, manifest) {
		t.Fatalf("ReadManifest() = %#v, want %#v", got, manifest)
	}
}
