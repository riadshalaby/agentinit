package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"time"
)

const manifestRelPath = ".ai/.manifest.json"

const (
	managementMarker = "marker"
	managementFull   = "full"
)

var manifestNow = func() time.Time {
	return time.Now().UTC()
}

var readManifestBuildInfo = debug.ReadBuildInfo

type Manifest struct {
	Version     string         `json:"version"`
	GeneratedAt string         `json:"generated_at"`
	Files       []ManifestFile `json:"files"`
}

type ManifestFile struct {
	Path       string `json:"path"`
	Management string `json:"management"`
}

func GenerateManifest(files map[string]string, version string) Manifest {
	paths := make([]string, 0, len(files))
	for path := range files {
		if shouldIncludeInManifest(path) {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)

	entries := make([]ManifestFile, 0, len(paths))
	for _, path := range paths {
		entries = append(entries, ManifestFile{
			Path:       path,
			Management: managementTypeForPath(path),
		})
	}

	return Manifest{
		Version:     normalizeManifestVersion(version),
		GeneratedAt: manifestNow().Format(time.RFC3339),
		Files:       entries,
	}
}

func WriteManifest(targetDir string, manifest Manifest) error {
	path := filepath.Join(targetDir, manifestRelPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}

	content, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}
	content = append(content, '\n')

	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write manifest %s: %w", path, err)
	}
	return nil
}

func ReadManifest(targetDir string) (Manifest, error) {
	path := filepath.Join(targetDir, manifestRelPath)
	content, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest %s: %w", path, err)
	}

	var manifest Manifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("unmarshal manifest %s: %w", path, err)
	}
	return manifest, nil
}

func currentVersion() string {
	info, ok := readManifestBuildInfo()
	if !ok || info == nil || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "(dev)"
	}
	return info.Main.Version
}

func shouldIncludeInManifest(path string) bool {
	if path == manifestRelPath {
		return false
	}
	_, excluded := manifestExcludedPaths[path]
	return !excluded
}

func managementTypeForPath(path string) string {
	if path == "AGENTS.md" {
		return managementMarker
	}
	return managementFull
}

func normalizeManifestVersion(version string) string {
	if version == "" {
		return "(dev)"
	}
	return version
}

var manifestExcludedPaths = map[string]struct{}{
	"ROADMAP.md":                  {},
	"CLAUDE.md":                   {},
	"README.md":                   {},
	".ai/TASKS.template.md":       {},
	".ai/PLAN.template.md":        {},
	".ai/REVIEW.template.md":      {},
	".ai/TEST_REPORT.template.md": {},
	".ai/HANDOFF.template.md":     {},
}
