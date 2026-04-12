package update

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/riadshalaby/agentinit/internal/scaffold"
)

var fallbackKnownPaths = []string{
	".ai/prompts/implementer.md",
	".ai/prompts/planner.md",
	".ai/prompts/po.md",
	".ai/prompts/reviewer.md",
	".gitattributes",
	".gitignore",
	"AGENTS.md",
	"ROADMAP.template.md",
	"scripts/ai-implement.sh",
	"scripts/ai-launch.sh",
	"scripts/ai-plan.sh",
	"scripts/ai-po.sh",
	"scripts/ai-pr.sh",
	"scripts/ai-review.sh",
	"scripts/ai-start-cycle.sh",
}

func DiscoverManagedFiles(targetDir string) scaffold.Manifest {
	files := make([]scaffold.ManifestFile, 0, len(fallbackKnownPaths))
	for _, relPath := range fallbackKnownPaths {
		if _, err := os.Stat(filepath.Join(targetDir, relPath)); err != nil {
			continue
		}

		management := "full"
		if relPath == "AGENTS.md" {
			management = "marker"
		}

		files = append(files, scaffold.ManifestFile{
			Path:       relPath,
			Management: management,
		})
	}

	return scaffold.Manifest{Files: files}
}

func InferProjectType(targetDir string) string {
	if fileExists(filepath.Join(targetDir, "go.mod")) {
		return "go"
	}
	if fileExists(filepath.Join(targetDir, "package.json")) {
		return "node"
	}
	if fileExists(filepath.Join(targetDir, "pom.xml")) {
		return "java"
	}
	if agentsType := inferProjectTypeFromAgents(filepath.Join(targetDir, "AGENTS.md")); agentsType != "" {
		return agentsType
	}
	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func inferProjectTypeFromAgents(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	text := string(content)
	switch {
	case strings.Contains(text, "go fmt ./...") || strings.Contains(text, "go vet ./...") || strings.Contains(text, "go test ./..."):
		return "go"
	case strings.Contains(text, "npm run lint") || strings.Contains(text, "npm run build") || strings.Contains(text, "npm test"):
		return "node"
	case strings.Contains(text, "mvn -q spotless:apply") || strings.Contains(text, "mvn -T 1C -q test"):
		return "java"
	default:
		return ""
	}
}
