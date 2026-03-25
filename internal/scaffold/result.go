package scaffold

import (
	"path/filepath"

	"github.com/riadshalaby/agentinit/internal/template"
)

type Result struct {
	ProjectName        string
	ProjectType        string
	TargetDir          string
	GitInitDone        bool
	DocumentationPath  string
	KeyPaths           []KeyPath
	ValidationCommands []template.ValidationCommand
}

type KeyPath struct {
	Path        string
	Description string
}

func buildResult(name, projectType, targetDir string, initGit bool, validationCommands []template.ValidationCommand) Result {
	return Result{
		ProjectName:        name,
		ProjectType:        projectType,
		TargetDir:          targetDir,
		GitInitDone:        initGit,
		DocumentationPath:  filepath.Join(targetDir, "README.md"),
		KeyPaths:           defaultKeyPaths(),
		ValidationCommands: append([]template.ValidationCommand(nil), validationCommands...),
	}
}

func defaultKeyPaths() []KeyPath {
	return []KeyPath{
		{Path: "README.md", Description: "project overview and setup"},
		{Path: "CLAUDE.md", Description: "project rules and agent workflow"},
		{Path: "ROADMAP.md", Description: "project goals to edit first"},
		{Path: ".ai/", Description: "planning, review, and handoff templates"},
		{Path: "scripts/", Description: "launchers for plan, implement, review, and PR sync"},
	}
}
