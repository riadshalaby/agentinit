package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cycleLookPath   = exec.LookPath
	cycleReadFile   = os.ReadFile
	cycleWriteFile  = os.WriteFile
	cycleMkdirAll   = os.MkdirAll
	cycleStat       = os.Stat
	cycleRunCommand = func(ctx context.Context, name string, args ...string) error {
		command := exec.CommandContext(ctx, name, args...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		return command.Run()
	}
	cycleOutputCommand = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		command := exec.CommandContext(ctx, name, args...)
		var stdout bytes.Buffer
		command.Stdout = &stdout
		command.Stderr = os.Stderr
		err := command.Run()
		return stdout.Bytes(), err
	}
)

var cycleCmd = &cobra.Command{
	Use:   "cycle",
	Short: "Manage development cycle bootstrap and close-out",
}

var cycleStartCmd = &cobra.Command{
	Use:   "start <branch>",
	Short: "Start a new development cycle",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCycleStart(cmd.Context(), args[0])
	},
}

type cycleCopySpec struct {
	source string
	target string
}

var cycleBootstrapFiles = []cycleCopySpec{
	{source: ".ai/PLAN.template.md", target: ".ai/PLAN.md"},
	{source: ".ai/REVIEW.template.md", target: ".ai/REVIEW.md"},
	{source: ".ai/TASKS.template.md", target: ".ai/TASKS.md"},
	{source: ".ai/HANDOFF.template.md", target: ".ai/HANDOFF.md"},
	{source: "ROADMAP.template.md", target: "ROADMAP.md"},
}

func init() {
	cycleCmd.AddCommand(cycleStartCmd)
	rootCmd.AddCommand(cycleCmd)
}

func runCycleStart(ctx context.Context, branchName string) error {
	if err := requireCycleCommand("git"); err != nil {
		return err
	}
	if err := validateCycleBranchName(ctx, branchName); err != nil {
		return err
	}

	repoRoot, err := getWorkingDir()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	if err := ensureCycleWorkingTreeClean(ctx); err != nil {
		return err
	}
	if err := ensureCycleBranchAvailable(ctx, branchName); err != nil {
		return err
	}
	if err := checkoutCycleBranch(ctx, branchName); err != nil {
		return err
	}
	if err := copyCycleBootstrapFiles(repoRoot); err != nil {
		return err
	}
	if err := commitCycleBootstrap(ctx, branchName); err != nil {
		return err
	}

	_, err = fmt.Fprintf(cliOutput, "Started new cycle on branch '%s'.\n", branchName)
	return err
}

func requireCycleCommand(name string) error {
	if _, err := cycleLookPath(name); err != nil {
		return fmt.Errorf("missing required command: %s", name)
	}
	return nil
}

func validateCycleBranchName(ctx context.Context, branchName string) error {
	if branchName == "" {
		return fmt.Errorf("branch name is required")
	}
	switch branchName {
	case "feature/", "fix/", "chore/":
		return fmt.Errorf("branch name must include a suffix after the prefix")
	}
	if !strings.HasPrefix(branchName, "feature/") && !strings.HasPrefix(branchName, "fix/") && !strings.HasPrefix(branchName, "chore/") {
		return fmt.Errorf("branch name must start with feature/, fix/, or chore/")
	}
	if _, err := cycleOutputCommand(ctx, "git", "check-ref-format", "--branch", branchName); err != nil {
		return fmt.Errorf("branch name is not a valid git branch name")
	}
	return nil
}

func ensureCycleWorkingTreeClean(ctx context.Context) error {
	if err := cycleRunCommand(ctx, "git", "diff", "--quiet"); err != nil {
		return fmt.Errorf("working tree is dirty — commit or stash changes before starting a new cycle")
	}
	if err := cycleRunCommand(ctx, "git", "diff", "--cached", "--quiet"); err != nil {
		return fmt.Errorf("working tree is dirty — commit or stash changes before starting a new cycle")
	}
	output, err := cycleOutputCommand(ctx, "git", "ls-files", "--others", "--exclude-standard")
	if err != nil {
		return fmt.Errorf("failed to inspect untracked files: %w", err)
	}
	if len(bytes.TrimSpace(output)) > 0 {
		return fmt.Errorf("untracked files present — commit, stash, or gitignore them before starting a new cycle")
	}
	return nil
}

func ensureCycleBranchAvailable(ctx context.Context, branchName string) error {
	if _, err := cycleOutputCommand(ctx, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+branchName); err == nil {
		return fmt.Errorf("branch %q already exists locally", branchName)
	}
	if _, err := cycleOutputCommand(ctx, "git", "ls-remote", "--exit-code", "--heads", "origin", branchName); err == nil {
		return fmt.Errorf("branch %q already exists on origin", branchName)
	}
	return nil
}

func checkoutCycleBranch(ctx context.Context, branchName string) error {
	if err := cycleRunCommand(ctx, "git", "checkout", "main"); err != nil {
		return fmt.Errorf("failed to checkout main")
	}
	if err := cycleRunCommand(ctx, "git", "pull", "--ff-only", "origin", "main"); err != nil {
		return fmt.Errorf("failed to fast-forward local main from origin/main")
	}
	if err := cycleRunCommand(ctx, "git", "checkout", "-b", branchName); err != nil {
		return fmt.Errorf("failed to create branch %q", branchName)
	}
	return nil
}

func copyCycleBootstrapFiles(repoRoot string) error {
	for _, spec := range cycleBootstrapFiles {
		sourcePath := filepath.Join(repoRoot, spec.source)
		targetPath := filepath.Join(repoRoot, spec.target)

		data, err := cycleReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("read template %q: %w", spec.source, err)
		}
		info, err := cycleStat(sourcePath)
		if err != nil {
			return fmt.Errorf("stat template %q: %w", spec.source, err)
		}
		if err := cycleMkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return fmt.Errorf("create parent directory for %q: %w", spec.target, err)
		}
		if err := cycleWriteFile(targetPath, data, info.Mode().Perm()); err != nil {
			return fmt.Errorf("write bootstrap file %q: %w", spec.target, err)
		}
	}
	return nil
}

func commitCycleBootstrap(ctx context.Context, branchName string) error {
	addArgs := []string{"add"}
	for _, spec := range cycleBootstrapFiles {
		addArgs = append(addArgs, spec.target)
	}
	if err := cycleRunCommand(ctx, "git", addArgs...); err != nil {
		return fmt.Errorf("failed to stage cycle bootstrap files")
	}
	if err := cycleRunCommand(ctx, "git", "commit", "-m", "chore: start cycle "+filepath.Base(branchName)); err != nil {
		return fmt.Errorf("failed to commit cycle bootstrap files")
	}
	if err := cycleRunCommand(ctx, "git", "push", "-u", "origin", branchName); err != nil {
		return fmt.Errorf("failed to push branch %q to origin", branchName)
	}
	return nil
}
