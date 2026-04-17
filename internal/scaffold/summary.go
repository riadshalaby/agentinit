package scaffold

import (
	"fmt"
	"strings"
)

type SummaryModel struct {
	Heading           string
	DocumentationPath string
	Rows              []SummaryRow
	NextSteps         []string
}

type SummaryRow struct {
	Label string
	Value string
}

func BuildSummary(result Result) SummaryModel {
	rows := []SummaryRow{
		{Label: "Name", Value: result.ProjectName},
	}
	if result.ProjectType != "" {
		rows = append(rows, SummaryRow{Label: "Type", Value: result.ProjectType})
	}
	rows = append(rows,
		SummaryRow{Label: "Path", Value: result.TargetDir},
		SummaryRow{Label: "Git", Value: gitStatusLabel(result.GitInitDone)},
		SummaryRow{Label: "Documentation", Value: result.DocumentationPath},
	)
	for _, keyPath := range result.KeyPaths {
		rows = append(rows, SummaryRow{
			Label: keyPath.Path,
			Value: keyPath.Description,
		})
	}

	nextSteps := []string{
		fmt.Sprintf("cd %s", result.TargetDir),
		"Edit ROADMAP.md with your project goals",
		"Start a development cycle: agentinit cycle start feature/<scope>",
		"Run the planner: scripts/ai-plan.sh",
	}
	if len(result.ValidationCommands) > 0 {
		lines := []string{"Validate the project:"}
		for _, cmd := range result.ValidationCommands {
			lines = append(lines, cmd.Command)
		}
		nextSteps = append(nextSteps, strings.Join(lines, "\n"))
	}

	return SummaryModel{
		Heading:           "Project scaffold complete!",
		DocumentationPath: result.DocumentationPath,
		Rows:              rows,
		NextSteps:         nextSteps,
	}
}

func FormatCLISummary(model SummaryModel) string {
	var b strings.Builder
	b.WriteString(model.Heading)
	b.WriteString("\n\n")
	b.WriteString("Documentation: ")
	b.WriteString(model.DocumentationPath)
	b.WriteString("\n\n")
	writeRows(&b, model.Rows, func(label, value string, width int) string {
		return fmt.Sprintf("%s%-*s %s\n", "  ", width, label, value)
	})
	writeNextSteps(&b, model.NextSteps, "Next steps:")
	return strings.TrimRight(b.String(), "\n")
}

func FormatWizardSummary(model SummaryModel) (title string, body string) {
	var b strings.Builder
	b.WriteString("Documentation: ")
	b.WriteString(model.DocumentationPath)
	b.WriteString("\n\n")
	writeRows(&b, model.Rows, func(label, value string, width int) string {
		return fmt.Sprintf("%-*s %s\n", width, label, value)
	})
	writeNextSteps(&b, model.NextSteps, "Next steps:")
	return model.Heading, strings.TrimRight(b.String(), "\n")
}

func writeRows(b *strings.Builder, rows []SummaryRow, format func(label, value string, width int) string) {
	width := 0
	for _, row := range rows {
		if len(row.Label) > width {
			width = len(row.Label)
		}
	}
	for _, row := range rows {
		b.WriteString(format(row.Label, row.Value, width))
	}
	b.WriteString("\n")
}

func writeNextSteps(b *strings.Builder, steps []string, heading string) {
	if len(steps) == 0 {
		return
	}
	b.WriteString(heading)
	b.WriteString("\n")
	for i, step := range steps {
		lines := strings.Split(step, "\n")
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, lines[0]))
		for _, line := range lines[1:] {
			b.WriteString("   ")
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
}

func gitStatusLabel(initGit bool) string {
	if initGit {
		return "initialized"
	}
	return "not initialized"
}
