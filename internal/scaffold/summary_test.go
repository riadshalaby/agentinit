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
		Profile:           "full",
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
	if len(model.NextSteps) != 5 {
		t.Fatalf("NextSteps len = %d, want 5 without validation step", len(model.NextSteps))
	}
}

func TestBuildSummaryIncludesValidationCommandsForTypedProject(t *testing.T) {
	result := buildResult("demo", "go", "full", "/tmp/demo", true, []template.ValidationCommand{
		{Label: "fmt", Command: "go fmt ./..."},
		{Label: "vet", Command: "go vet ./..."},
		{Label: "test", Command: "go test ./..."},
	})

	model := BuildSummary(result)

	if len(model.NextSteps) != 6 {
		t.Fatalf("NextSteps len = %d, want 6 with validation step", len(model.NextSteps))
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
	model := BuildSummary(buildResult("demo", "go", "full", "/tmp/demo", true, []template.ValidationCommand{
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
	if !strings.Contains(summary, "6. Validate the project:") {
		t.Fatalf("summary = %q", summary)
	}
}

func TestFormatWizardSummaryRendersSameContent(t *testing.T) {
	model := BuildSummary(buildResult("demo", "", "full", filepath.Join("/tmp", "demo"), false, nil))

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
	if !strings.Contains(body, "4. Run the planner: aide plan") {
		t.Fatalf("body = %q", body)
	}
	if !strings.Contains(body, "5. Start the implementer, reviewer, and PO sessions: aide implement, aide review, and aide po") {
		t.Fatalf("body = %q", body)
	}
}

func TestBuildSummaryUsesLiteProfileNextSteps(t *testing.T) {
	model := BuildSummary(buildResult("demo", "", "lite", "/tmp/demo", false, nil))

	if len(model.NextSteps) != 5 {
		t.Fatalf("NextSteps len = %d, want 5 for lite profile", len(model.NextSteps))
	}
	if got := model.NextSteps[3]; got != "Run the planner: aide plan" {
		t.Fatalf("planner step = %q", got)
	}
	if got := model.NextSteps[4]; got != "Start the dev session: aide dev\nSee aide profile lite --help for the lite workflow details." {
		t.Fatalf("dev step = %q", got)
	}
}
