package cmd

import (
	"bytes"
	"io/fs"
	"strings"
	"testing"

	"github.com/riadshalaby/agentinit/internal/prereq"
	updater "github.com/riadshalaby/agentinit/internal/update"
)

func TestUpdateCommandUsesFlagsAndPrintsChanges(t *testing.T) {
	originalRunUpdate := runUpdate
	originalRunToolCheck := runUpdateToolCheck
	originalOutput := cliOutput
	originalDir := updateTargetDir
	originalDryRun := updateDryRun
	originalStdinStat := stdinStat
	originalOutputStat := updateOutputStat
	t.Cleanup(func() {
		runUpdate = originalRunUpdate
		runUpdateToolCheck = originalRunToolCheck
		cliOutput = originalOutput
		updateTargetDir = originalDir
		updateDryRun = originalDryRun
		stdinStat = originalStdinStat
		updateOutputStat = originalOutputStat
	})

	updateTargetDir = "/tmp/project"
	updateDryRun = true
	stdinStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{mode: fs.ModeCharDevice}, nil
	}
	updateOutputStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{mode: fs.ModeCharDevice}, nil
	}

	var output bytes.Buffer
	cliOutput = &output
	runUpdateToolCheck = func(prereq.Commander) error { return nil }
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
	originalRunToolCheck := runUpdateToolCheck
	originalOutput := cliOutput
	originalDir := updateTargetDir
	originalDryRun := updateDryRun
	originalStdinStat := stdinStat
	originalOutputStat := updateOutputStat
	t.Cleanup(func() {
		runUpdate = originalRunUpdate
		runUpdateToolCheck = originalRunToolCheck
		cliOutput = originalOutput
		updateTargetDir = originalDir
		updateDryRun = originalDryRun
		stdinStat = originalStdinStat
		updateOutputStat = originalOutputStat
	})

	updateTargetDir = "/tmp/project"
	updateDryRun = false
	stdinStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{mode: fs.ModeCharDevice}, nil
	}
	updateOutputStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{mode: fs.ModeCharDevice}, nil
	}

	var output bytes.Buffer
	cliOutput = &output
	runUpdateToolCheck = func(prereq.Commander) error { return nil }
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

func TestUpdateCommandRunsToolCheckAfterPrintingChanges(t *testing.T) {
	originalRunUpdate := runUpdate
	originalRunToolCheck := runUpdateToolCheck
	originalOutput := cliOutput
	originalDir := updateTargetDir
	originalDryRun := updateDryRun
	originalStdinStat := stdinStat
	originalOutputStat := updateOutputStat
	t.Cleanup(func() {
		runUpdate = originalRunUpdate
		runUpdateToolCheck = originalRunToolCheck
		cliOutput = originalOutput
		updateTargetDir = originalDir
		updateDryRun = originalDryRun
		stdinStat = originalStdinStat
		updateOutputStat = originalOutputStat
	})

	updateTargetDir = "/tmp/project"
	updateDryRun = false
	stdinStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{mode: fs.ModeCharDevice}, nil
	}
	updateOutputStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{mode: fs.ModeCharDevice}, nil
	}

	var output bytes.Buffer
	cliOutput = &output
	var calls []string
	runUpdate = func(string, bool) (updater.Result, error) {
		calls = append(calls, "update")
		return updater.Result{
			Changes: []updater.Change{{Path: "AGENTS.md", Action: "update"}},
		}, nil
	}
	runUpdateToolCheck = func(prereq.Commander) error {
		calls = append(calls, "toolcheck")
		if got := output.String(); got != "Updated AGENTS.md (update)\n" {
			t.Fatalf("output before tool check = %q", got)
		}
		return nil
	}

	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	if got, want := strings.Join(calls, ","), "update,toolcheck"; got != want {
		t.Fatalf("calls = %q, want %q", got, want)
	}
}

func TestUpdateCommandSkipsToolCheckWithoutTTY(t *testing.T) {
	originalRunUpdate := runUpdate
	originalRunToolCheck := runUpdateToolCheck
	originalOutput := cliOutput
	originalDir := updateTargetDir
	originalDryRun := updateDryRun
	originalStdinStat := stdinStat
	originalOutputStat := updateOutputStat
	t.Cleanup(func() {
		runUpdate = originalRunUpdate
		runUpdateToolCheck = originalRunToolCheck
		cliOutput = originalOutput
		updateTargetDir = originalDir
		updateDryRun = originalDryRun
		stdinStat = originalStdinStat
		updateOutputStat = originalOutputStat
	})

	updateTargetDir = "/tmp/project"
	updateDryRun = false
	stdinStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{}, nil
	}
	updateOutputStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{}, nil
	}

	var output bytes.Buffer
	cliOutput = &output
	var calls []string
	runUpdate = func(string, bool) (updater.Result, error) {
		calls = append(calls, "update")
		return updater.Result{
			Changes: []updater.Change{{Path: "AGENTS.md", Action: "update"}},
		}, nil
	}
	runUpdateToolCheck = func(prereq.Commander) error {
		calls = append(calls, "toolcheck")
		return nil
	}

	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	if got, want := strings.Join(calls, ","), "update"; got != want {
		t.Fatalf("calls = %q, want %q", got, want)
	}
}
