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

const (
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
)

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
		ToolPermissions:    ov.ToolPermissions,
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

	paths := managedPaths(currentByPath, desiredByPath)
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

	deletionChanges, err := deleteRemovedManagedFiles(targetDir, currentByPath, desiredByPath, dryRun)
	if err != nil {
		return Result{}, err
	}
	changes = append(changes, deletionChanges...)

	migrationChanges, err := migrateExcludedFiles(targetDir, dryRun)
	if err != nil {
		return Result{}, err
	}
	changes = append(changes, migrationChanges...)

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

func managedPaths(currentByPath, desiredByPath map[string]string) []string {
	pathSet := make(map[string]struct{}, len(currentByPath)+len(desiredByPath))
	for path := range currentByPath {
		pathSet[path] = struct{}{}
	}
	for path := range desiredByPath {
		pathSet[path] = struct{}{}
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
		return renderedContent, actionCreate, true, nil
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
		return updated, actionUpdate, true, nil
	default:
		if renderedContent == existing {
			return existing, "", false, nil
		}
		return renderedContent, actionUpdate, true, nil
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
			return true, actionCreate, nil
		}
		return false, "", fmt.Errorf("read %s: %w", path, err)
	}

	if bytes.Equal(existing, desired) {
		return false, "", nil
	}
	return true, actionUpdate, nil
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

func deleteRemovedManagedFiles(targetDir string, currentByPath, desiredByPath map[string]string, dryRun bool) ([]Change, error) {
	paths := make([]string, 0, len(currentByPath))
	for relPath := range currentByPath {
		if _, ok := desiredByPath[relPath]; ok {
			continue
		}
		paths = append(paths, relPath)
	}
	sort.Strings(paths)

	changes := make([]Change, 0, len(paths))
	for _, relPath := range paths {
		absPath := filepath.Join(targetDir, relPath)
		if !fileExists(absPath) {
			continue
		}

		changes = append(changes, Change{Path: relPath, Action: actionDelete})
		if dryRun {
			continue
		}
		if err := os.Remove(absPath); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("delete %s: %w", absPath, err)
		}
	}

	return changes, nil
}

func migrateExcludedFiles(targetDir string, dryRun bool) ([]Change, error) {
	var changes []Change

	tasksTemplateChanges, err := migrateTasksTemplate(targetDir, dryRun)
	if err != nil {
		return nil, err
	}
	changes = append(changes, tasksTemplateChanges...)

	configChanges, err := migrateConfig(targetDir, dryRun)
	if err != nil {
		return nil, err
	}
	changes = append(changes, configChanges...)

	testReportChanges, err := deleteIfExists(targetDir, ".ai/TEST_REPORT.template.md", dryRun)
	if err != nil {
		return nil, err
	}
	changes = append(changes, testReportChanges...)

	scriptChanges, err := migrateScripts(targetDir, dryRun)
	if err != nil {
		return nil, err
	}
	changes = append(changes, scriptChanges...)

	return changes, nil
}

func migrateScripts(targetDir string, dryRun bool) ([]Change, error) {
	scriptPaths := []string{
		"scripts/ai-implement.sh",
		"scripts/ai-launch.sh",
		"scripts/ai-plan.sh",
		"scripts/ai-po.sh",
		"scripts/ai-pr.sh",
		"scripts/ai-review.sh",
		"scripts/ai-start-cycle.sh",
	}

	changes := make([]Change, 0, len(scriptPaths)+1)
	for _, relPath := range scriptPaths {
		fileChanges, err := deleteIfExists(targetDir, relPath, dryRun)
		if err != nil {
			return nil, err
		}
		changes = append(changes, fileChanges...)
	}

	scriptsDir := filepath.Join(targetDir, "scripts")
	info, err := os.Stat(scriptsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return changes, nil
		}
		return nil, fmt.Errorf("stat %s: %w", scriptsDir, err)
	}
	if !info.IsDir() {
		return changes, nil
	}

	entries, err := os.ReadDir(scriptsDir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", scriptsDir, err)
	}
	if len(entries) != 0 {
		return changes, nil
	}

	changes = append(changes, Change{Path: "scripts", Action: actionDelete})
	if dryRun {
		return changes, nil
	}
	if err := os.Remove(scriptsDir); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("delete %s: %w", scriptsDir, err)
	}
	return changes, nil
}

func migrateTasksTemplate(targetDir string, dryRun bool) ([]Change, error) {
	const relPath = ".ai/TASKS.template.md"

	_, changed, err := rewriteFileIfNeeded(targetDir, relPath, func(existing string) (string, bool, error) {
		updated := existing
		replacements := []struct {
			old string
			new string
		}{
			{"- `ready_for_test`\n", ""},
			{"- `in_testing`\n", ""},
			{"- `test_failed`\n", ""},
			{"- implementer moves tasks into `in_implementation`, `ready_for_review`, and `done`, and resumes work from `changes_requested`, `test_failed`, and `ready_to_commit`\n", "- implementer moves tasks into `in_implementation`, `ready_for_review`, and `done`, and resumes work from `changes_requested` and `ready_to_commit`\n"},
			{"- reviewer moves tasks into `in_review`, `ready_for_test`, or `changes_requested`\n", "- reviewer moves tasks into `in_review`, `ready_to_commit`, or `changes_requested`\n"},
			{"- tester moves tasks into `in_testing`, `ready_to_commit`, or `test_failed`\n", ""},
		}
		for _, replacement := range replacements {
			updated = strings.ReplaceAll(updated, replacement.old, replacement.new)
		}
		return updated, updated != existing, nil
	}, dryRun)
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, nil
	}
	return []Change{{Path: relPath, Action: actionUpdate}}, nil
}

func migrateConfig(targetDir string, dryRun bool) ([]Change, error) {
	const relPath = ".ai/config.json"

	_, changed, err := rewriteFileIfNeeded(targetDir, relPath, func(existing string) (string, bool, error) {
		var doc map[string]json.RawMessage
		if err := json.Unmarshal([]byte(existing), &doc); err != nil {
			return "", false, fmt.Errorf("parse %s: %w", relPath, err)
		}

		rolesRaw, ok := doc["roles"]
		if !ok {
			return existing, false, nil
		}

		var roles map[string]json.RawMessage
		if err := json.Unmarshal(rolesRaw, &roles); err != nil {
			return "", false, fmt.Errorf("parse %s roles: %w", relPath, err)
		}
		if _, ok := roles["test"]; !ok {
			return existing, false, nil
		}

		delete(roles, "test")
		updatedRoles, err := json.MarshalIndent(roles, "  ", "  ")
		if err != nil {
			return "", false, fmt.Errorf("marshal %s roles: %w", relPath, err)
		}
		doc["roles"] = updatedRoles

		updatedDoc, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return "", false, fmt.Errorf("marshal %s: %w", relPath, err)
		}
		return string(append(updatedDoc, '\n')), true, nil
	}, dryRun)
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, nil
	}
	return []Change{{Path: relPath, Action: actionUpdate}}, nil
}

func rewriteFileIfNeeded(targetDir, relPath string, rewrite func(existing string) (string, bool, error), dryRun bool) (string, bool, error) {
	absPath := filepath.Join(targetDir, relPath)
	existingBytes, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("read %s: %w", absPath, err)
	}

	updated, changed, err := rewrite(string(existingBytes))
	if err != nil {
		return "", false, err
	}
	if !changed {
		return updated, false, nil
	}
	if dryRun {
		return updated, true, nil
	}
	if err := writeFile(targetDir, relPath, updated); err != nil {
		return "", false, err
	}
	return updated, true, nil
}

func deleteIfExists(targetDir, relPath string, dryRun bool) ([]Change, error) {
	absPath := filepath.Join(targetDir, relPath)
	if !fileExists(absPath) {
		return nil, nil
	}
	if !dryRun {
		if err := os.Remove(absPath); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("delete %s: %w", absPath, err)
		}
	}
	return []Change{{Path: relPath, Action: actionDelete}}, nil
}
