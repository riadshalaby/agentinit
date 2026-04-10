package scaffold

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/riadshalaby/agentinit/internal/template"
)

func TestBuildSummaryIncludesDocumentationAndKeyPaths(t *testing.T) {
	result := Result{
		ProjectName:       "demo",
		TargetDir:         "/tmp/demo",
		GitInitDone:       false,
		DocumentationPath: "/tmp/demo/README.md",
		KeyPaths:          defaultKeyPaths(),
	}

	model := BuildSummary(result)

	if model.DocumentationPath != "/tmp/demo/README.md" {
		t.Fatalf("DocumentationPath = %q", model.DocumentationPath)
	}
	if len(model.Rows) < 8 {
		t.Fatalf("Rows len = %d, want key summary rows", len(model.Rows))
	}
	if model.Rows[0] != (SummaryRow{Label: "Name", Value: "demo"}) {
		t.Fatalf("first row = %+v", model.Rows[0])
	}
	foundDocs := false
	foundAI := false
	for _, row := range model.Rows {
		if row.Label == "Documentation" && row.Value == "/tmp/demo/README.md" {
			foundDocs = true
		}
		if row.Label == ".ai/" {
			foundAI = true
		}
	}
	if !foundDocs {
		t.Fatal("expected documentation row")
	}
	if !foundAI {
		t.Fatal("expected .ai/ key path row")
	}
	if len(model.NextSteps) != 4 {
		t.Fatalf("NextSteps len = %d, want 4 without validation step", len(model.NextSteps))
	}
}

func TestBuildSummaryIncludesValidationCommandsForTypedProject(t *testing.T) {
	result := buildResult("demo", "go", "/tmp/demo", true, []template.ValidationCommand{
		{Label: "fmt", Command: "go fmt ./..."},
		{Label: "vet", Command: "go vet ./..."},
		{Label: "test", Command: "go test ./..."},
	})

	model := BuildSummary(result)

	if len(model.NextSteps) != 5 {
		t.Fatalf("NextSteps len = %d, want 5 with validation step", len(model.NextSteps))
	}
	last := model.NextSteps[len(model.NextSteps)-1]
	if !strings.Contains(last, "Validate the project:") {
		t.Fatalf("validation step = %q", last)
	}
	if !strings.Contains(last, "go test ./...") {
		t.Fatalf("validation step = %q", last)
	}
}

func TestFormatCLISummaryRendersAlignedSummary(t *testing.T) {
	model := BuildSummary(buildResult("demo", "go", "/tmp/demo", true, []template.ValidationCommand{
		{Label: "test", Command: "go test ./..."},
	}))

	summary := FormatCLISummary(model)

	if !strings.Contains(summary, "Project scaffold complete!") {
		t.Fatalf("summary = %q", summary)
	}
	if !strings.Contains(summary, "Documentation: /tmp/demo/README.md") {
		t.Fatalf("summary = %q", summary)
	}
	if !strings.Contains(summary, "  Name") {
		t.Fatalf("summary = %q", summary)
	}
	if !strings.Contains(summary, "1. cd /tmp/demo") {
		t.Fatalf("summary = %q", summary)
	}
	if !strings.Contains(summary, "5. Validate the project:") {
		t.Fatalf("summary = %q", summary)
	}
}

func TestFormatWizardSummaryRendersSameContent(t *testing.T) {
	model := BuildSummary(buildResult("demo", "", filepath.Join("/tmp", "demo"), false, nil))

	title, body := FormatWizardSummary(model)

	if title != "Project scaffold complete!" {
		t.Fatalf("title = %q", title)
	}
	if !strings.Contains(body, "Documentation: /tmp/demo/README.md") {
		t.Fatalf("body = %q", body)
	}
	if !strings.Contains(body, "Git           not initialized") {
		t.Fatalf("body = %q", body)
	}
	if !strings.Contains(body, "README.md     project overview and setup") {
		t.Fatalf("body = %q", body)
	}
	if !strings.Contains(body, "AGENTS.md     project-specific and workflow-managed agent rules") {
		t.Fatalf("body = %q", body)
	}
	if !strings.Contains(body, "4. Run the planner: scripts/ai-plan.sh") {
		t.Fatalf("body = %q", body)
	}
}
