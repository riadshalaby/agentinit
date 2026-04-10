package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/riadshalaby/agentinit/internal/overlay"
	"github.com/riadshalaby/agentinit/internal/scaffold"
	"github.com/riadshalaby/agentinit/internal/template"
)

const manifestPath = ".ai/.manifest.json"

var readBuildInfo = debug.ReadBuildInfo

type Result struct {
	ProjectType  string
	UsedFallback bool
	Changes      []Change
}

type Change struct {
	Path   string
	Action string
}

func Run(targetDir string, dryRun bool) (Result, error) {
	if targetDir == "" {
		return Result{}, fmt.Errorf("target directory is required")
	}
	if info, err := os.Stat(targetDir); err != nil || !info.IsDir() {
		if err != nil {
			return Result{}, fmt.Errorf("stat target directory %s: %w", targetDir, err)
		}
		return Result{}, fmt.Errorf("%s is not a directory", targetDir)
	}

	projectType := InferProjectType(targetDir)
	ov, err := overlay.Get(projectType)
	if err != nil {
		return Result{}, err
	}

	files, err := template.RenderAll(&template.ProjectData{
		ProjectName:        filepath.Base(targetDir),
		ProjectType:        projectType,
		ValidationCommands: ov.ValidationCommands,
		PRTestPlanItems:    ov.PRTestPlanItems,
	})
	if err != nil {
		return Result{}, fmt.Errorf("render templates: %w", err)
	}

	currentManifest, usedFallback, err := loadManifest(targetDir)
	if err != nil {
		return Result{}, err
	}
	desiredManifest := scaffold.GenerateManifest(files, currentVersion())

	currentByPath := make(map[string]string, len(currentManifest.Files))
	for _, file := range currentManifest.Files {
		currentByPath[file.Path] = file.Management
	}
	desiredByPath := make(map[string]string, len(desiredManifest.Files))
	for _, file := range desiredManifest.Files {
		desiredByPath[file.Path] = file.Management
	}

	paths := managedPaths(targetDir, currentByPath, desiredByPath)
	changes := make([]Change, 0, len(paths)+1)

	for _, relPath := range paths {
		renderedContent, ok := files[relPath]
		if !ok {
			continue
		}

		management := currentByPath[relPath]
		if management == "" {
			management = desiredByPath[relPath]
		}

		desiredContent, action, changed, err := reconcileFile(targetDir, relPath, management, renderedContent)
		if err != nil {
			return Result{}, err
		}
		if !changed {
			continue
		}

		changes = append(changes, Change{Path: relPath, Action: action})
		if dryRun {
			continue
		}
		if err := writeFile(targetDir, relPath, desiredContent); err != nil {
			return Result{}, err
		}
	}

	manifestChanged, manifestAction, err := manifestNeedsWrite(targetDir, desiredManifest)
	if err != nil {
		return Result{}, err
	}
	if manifestChanged {
		changes = append(changes, Change{Path: manifestPath, Action: manifestAction})
		if !dryRun {
			if err := scaffold.WriteManifest(targetDir, desiredManifest); err != nil {
				return Result{}, err
			}
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Path == changes[j].Path {
			return changes[i].Action < changes[j].Action
		}
		return changes[i].Path < changes[j].Path
	})

	return Result{
		ProjectType:  projectType,
		UsedFallback: usedFallback,
		Changes:      changes,
	}, nil
}

func loadManifest(targetDir string) (scaffold.Manifest, bool, error) {
	if fileExists(filepath.Join(targetDir, manifestPath)) {
		manifest, err := scaffold.ReadManifest(targetDir)
		if err != nil {
			return scaffold.Manifest{}, false, err
		}
		return manifest, false, nil
	}
	return DiscoverManagedFiles(targetDir), true, nil
}

func managedPaths(targetDir string, currentByPath, desiredByPath map[string]string) []string {
	pathSet := make(map[string]struct{}, len(currentByPath)+len(desiredByPath))
	for path := range currentByPath {
		pathSet[path] = struct{}{}
	}
	for path := range desiredByPath {
		if _, ok := currentByPath[path]; ok || !fileExists(filepath.Join(targetDir, path)) {
			pathSet[path] = struct{}{}
		}
	}

	paths := make([]string, 0, len(pathSet))
	for path := range pathSet {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func reconcileFile(targetDir, relPath, management, renderedContent string) (string, string, bool, error) {
	absPath := filepath.Join(targetDir, relPath)
	existingBytes, err := os.ReadFile(absPath)
	if err != nil && !os.IsNotExist(err) {
		return "", "", false, fmt.Errorf("read %s: %w", absPath, err)
	}

	exists := err == nil
	if !exists {
		return renderedContent, "create", true, nil
	}

	existing := string(existingBytes)
	switch management {
	case "marker":
		_, managed, _, err := ExtractSections(renderedContent)
		if err != nil {
			return "", "", false, fmt.Errorf("extract managed section from rendered %s: %w", relPath, err)
		}
		updated, err := ReplaceManagedSection(existing, managed)
		if err != nil {
			return "", "", false, fmt.Errorf("replace managed section in %s: %w", relPath, err)
		}
		if updated == existing {
			return updated, "", false, nil
		}
		return updated, "update", true, nil
	default:
		if renderedContent == existing {
			return existing, "", false, nil
		}
		return renderedContent, "update", true, nil
	}
}

func writeFile(targetDir, relPath, content string) error {
	absPath := filepath.Join(targetDir, relPath)
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(absPath), err)
	}

	perm := os.FileMode(0o644)
	if strings.HasSuffix(relPath, ".sh") {
		perm = 0o755
	}
	if err := os.WriteFile(absPath, []byte(content), perm); err != nil {
		return fmt.Errorf("write %s: %w", absPath, err)
	}
	return nil
}

func manifestNeedsWrite(targetDir string, manifest scaffold.Manifest) (bool, string, error) {
	desired, err := marshalManifest(manifest)
	if err != nil {
		return false, "", err
	}

	path := filepath.Join(targetDir, manifestPath)
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, "create", nil
		}
		return false, "", fmt.Errorf("read %s: %w", path, err)
	}

	if bytes.Equal(existing, desired) {
		return false, "", nil
	}
	return true, "update", nil
}

func marshalManifest(manifest scaffold.Manifest) ([]byte, error) {
	content, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal manifest: %w", err)
	}
	return append(content, '\n'), nil
}

func currentVersion() string {
	info, ok := readBuildInfo()
	if !ok || info == nil || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "(dev)"
	}
	return info.Main.Version
}
