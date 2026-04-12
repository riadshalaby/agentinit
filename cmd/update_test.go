package cmd

import (
	"bytes"
	"testing"

	updater "github.com/riadshalaby/agentinit/internal/update"
)

func TestUpdateCommandUsesFlagsAndPrintsChanges(t *testing.T) {
	originalRunUpdate := runUpdate
	originalOutput := cliOutput
	originalDir := updateTargetDir
	originalDryRun := updateDryRun
	t.Cleanup(func() {
		runUpdate = originalRunUpdate
		cliOutput = originalOutput
		updateTargetDir = originalDir
		updateDryRun = originalDryRun
	})

	updateTargetDir = "/tmp/project"
	updateDryRun = true

	var output bytes.Buffer
	cliOutput = &output
	runUpdate = func(dir string, dryRun bool) (updater.Result, error) {
		if dir != "/tmp/project" {
			t.Fatalf("dir = %q, want %q", dir, "/tmp/project")
		}
		if !dryRun {
			t.Fatal("dryRun = false, want true")
		}
		return updater.Result{
			Changes: []updater.Change{
				{Path: "AGENTS.md", Action: "update"},
				{Path: ".ai/.manifest.json", Action: "create"},
				{Path: ".ai/prompts/tester.md", Action: "delete"},
			},
		}, nil
	}

	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	got := output.String()
	if !bytes.Contains([]byte(got), []byte("Would update AGENTS.md (update)")) {
		t.Fatalf("output = %q", got)
	}
	if !bytes.Contains([]byte(got), []byte("Would create .ai/.manifest.json (create)")) {
		t.Fatalf("output = %q", got)
	}
	if !bytes.Contains([]byte(got), []byte("Would delete .ai/prompts/tester.md (delete)")) {
		t.Fatalf("output = %q", got)
	}
}

func TestUpdateCommandPrintsNoChangeMessage(t *testing.T) {
	originalRunUpdate := runUpdate
	originalOutput := cliOutput
	originalDir := updateTargetDir
	originalDryRun := updateDryRun
	t.Cleanup(func() {
		runUpdate = originalRunUpdate
		cliOutput = originalOutput
		updateTargetDir = originalDir
		updateDryRun = originalDryRun
	})

	updateTargetDir = "/tmp/project"
	updateDryRun = false

	var output bytes.Buffer
	cliOutput = &output
	runUpdate = func(string, bool) (updater.Result, error) {
		return updater.Result{}, nil
	}

	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	if got := output.String(); got != "No managed files changed.\n" {
		t.Fatalf("output = %q", got)
	}
}
