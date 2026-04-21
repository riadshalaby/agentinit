package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/riadshalaby/agentinit/internal/overlay"
	updater "github.com/riadshalaby/agentinit/internal/update"
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

var cycleEndCmd = &cobra.Command{
	Use:   "end [version]",
	Short: "Close the current development cycle",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := ""
		if len(args) > 0 {
			version = args[0]
		}
		return runCycleEnd(cmd.Context(), version)
	},
}

type cycleCopySpec struct {
	source string
	target string
}

type prSyncOptions struct {
	BaseBranch string
	Title      string
	DryRun     bool
	SkipPush   bool
}

type commandResult struct {
	stdout string
	err    error
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
	cycleCmd.AddCommand(cycleEndCmd)
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

func runCycleEnd(ctx context.Context, version string) error {
	if err := requireCycleCommand("git"); err != nil {
		return err
	}

	repoRoot, err := getWorkingDir()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	incompleteTasks, err := cycleIncompleteTasks(repoRoot)
	if err != nil {
		return err
	}
	if len(incompleteTasks) > 0 {
		return fmt.Errorf("cannot close cycle; tasks not done: %s", strings.Join(incompleteTasks, ", "))
	}

	if err := cycleRunCommand(ctx, "git", "add", ".ai/"); err != nil {
		return fmt.Errorf("failed to stage .ai artifacts")
	}

	commitArgs := []string{"commit", "-m", "chore(ai): close cycle"}
	if version != "" {
		commitArgs = append(commitArgs, "-m", "Release-As: "+version)
	}
	if err := cycleRunCommand(ctx, "git", commitArgs...); err != nil {
		return fmt.Errorf("failed to commit cycle-close artifacts")
	}

	remoteURL, ok := cycleRemoteURL(ctx)
	if !ok || !isGitHubRemote(remoteURL) {
		_, err := fmt.Fprintln(cliOutput, "No GitHub remote detected — skipping PR.")
		return err
	}

	if err := requireCycleCommand("gh"); err != nil {
		return err
	}

	branch, err := cycleCurrentBranch(ctx)
	if err != nil {
		return err
	}
	if err := cycleRunCommand(ctx, "git", "push", "-u", "origin", branch); err != nil {
		return fmt.Errorf("failed to push branch %q to origin", branch)
	}

	return runPRSync(ctx, repoRoot, prSyncOptions{
		BaseBranch: "main",
		SkipPush:   true,
	})
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

func runPRSync(ctx context.Context, repoRoot string, opts prSyncOptions) error {
	if opts.BaseBranch == "" {
		opts.BaseBranch = "main"
	}

	if err := requireCycleCommand("git"); err != nil {
		return err
	}
	if !opts.DryRun {
		if err := requireCycleCommand("gh"); err != nil {
			return err
		}
	}

	branch, err := cycleCurrentBranch(ctx)
	if err != nil {
		return err
	}
	if branch == opts.BaseBranch {
		return fmt.Errorf("run pr on a feature branch, not %s", opts.BaseBranch)
	}

	remoteURL, hasRemote := cycleRemoteURL(ctx)
	if !opts.DryRun && (!hasRemote || !isGitHubRemote(remoteURL)) {
		_, err := fmt.Fprintln(cliOutput, "no remote configured — skipping PR")
		return err
	}
	if hasRemote {
		_ = cycleRunCommand(ctx, "git", "fetch", "origin", opts.BaseBranch)
	}

	baseRef, err := determinePRBaseRef(ctx, opts.BaseBranch)
	if err != nil {
		return err
	}
	mergeBase, err := cycleTrimmedOutput(ctx, "git", "merge-base", baseRef, "HEAD")
	if err != nil {
		return fmt.Errorf("determine merge base: %w", err)
	}
	rangeExpr := mergeBase + "..HEAD"

	commitCountText, err := cycleTrimmedOutput(ctx, "git", "rev-list", "--count", rangeExpr)
	if err != nil {
		return fmt.Errorf("count commits in PR range: %w", err)
	}
	commitCount, err := strconv.Atoi(commitCountText)
	if err != nil {
		return fmt.Errorf("parse commit count %q: %w", commitCountText, err)
	}
	if commitCount == 0 {
		return fmt.Errorf("no commits detected between %s and HEAD", baseRef)
	}

	existingPRNumber := ""
	if !opts.DryRun {
		existingPRNumber, err = findExistingPRNumber(ctx, branch, opts.BaseBranch)
		if err != nil {
			return err
		}
	}

	prTitle, err := existingOrDefaultTitle(ctx, opts.Title, existingPRNumber)
	if err != nil {
		return err
	}
	commitList, err := buildCommitListMarkdown(ctx, rangeExpr)
	if err != nil {
		return err
	}
	breakingChanges, err := buildBreakingChangesMarkdown(ctx, rangeExpr)
	if err != nil {
		return err
	}
	prBody := buildPRBody(branch, opts.BaseBranch, commitCount, breakingChanges, commitList, prTestPlanItems(repoRoot))

	if opts.DryRun {
		_, err := fmt.Fprintf(cliOutput, "Title: %s\n\n%s\n", prTitle, prBody)
		return err
	}

	if !opts.SkipPush {
		if err := cycleRunCommand(ctx, "git", "push", "-u", "origin", branch); err != nil {
			return fmt.Errorf("failed to push branch %q to origin", branch)
		}
	}

	if existingPRNumber != "" {
		if err := cycleRunCommand(ctx, "gh", "pr", "edit", existingPRNumber, "--title", prTitle, "--body", prBody); err != nil {
			return fmt.Errorf("update existing PR #%s: %w", existingPRNumber, err)
		}
		_, err := fmt.Fprintf(cliOutput, "Updated existing PR #%s for branch %s.\n", existingPRNumber, branch)
		return err
	}

	if err := cycleRunCommand(ctx, "gh", "pr", "create", "--base", opts.BaseBranch, "--head", branch, "--title", prTitle, "--body", prBody); err != nil {
		return fmt.Errorf("create PR for branch %s: %w", branch, err)
	}
	return nil
}

func cycleIncompleteTasks(repoRoot string) ([]string, error) {
	tasksPath := filepath.Join(repoRoot, ".ai", "TASKS.md")
	content, err := cycleReadFile(tasksPath)
	if err != nil {
		return nil, fmt.Errorf("read task board %q: %w", tasksPath, err)
	}

	var blocking []string
	for _, row := range strings.Split(string(content), "\n") {
		cols := parseMarkdownRow(row)
		if len(cols) < 6 || cols[0] == "Task ID" {
			continue
		}
		if cols[2] == "done" {
			continue
		}
		blocking = append(blocking, fmt.Sprintf("%s (%s)", cols[0], cols[2]))
	}
	return blocking, nil
}

func parseMarkdownRow(row string) []string {
	if !strings.HasPrefix(row, "|") || !strings.HasSuffix(row, "|") {
		return nil
	}
	parts := strings.Split(row, "|")
	if len(parts) < 3 {
		return nil
	}
	cols := make([]string, 0, len(parts)-2)
	for _, part := range parts[1 : len(parts)-1] {
		value := strings.TrimSpace(part)
		if value != "" {
			cols = append(cols, value)
			continue
		}
		cols = append(cols, "")
	}
	if len(cols) > 0 && strings.Trim(cols[0], "-") == "" {
		return nil
	}
	return cols
}

func cycleCurrentBranch(ctx context.Context) (string, error) {
	branch, err := cycleTrimmedOutput(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("determine current branch: %w", err)
	}
	return branch, nil
}

func cycleRemoteURL(ctx context.Context) (string, bool) {
	url, err := cycleTrimmedOutput(ctx, "git", "remote", "get-url", "origin")
	if err != nil || url == "" {
		return "", false
	}
	return url, true
}

func determinePRBaseRef(ctx context.Context, baseBranch string) (string, error) {
	if _, err := cycleOutputCommand(ctx, "git", "rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+baseBranch); err == nil {
		return "origin/" + baseBranch, nil
	}
	if _, err := cycleOutputCommand(ctx, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+baseBranch); err == nil {
		return baseBranch, nil
	}
	return "", fmt.Errorf("cannot determine PR base ref (expected origin/%s or %s)", baseBranch, baseBranch)
}

func findExistingPRNumber(ctx context.Context, branch, baseBranch string) (string, error) {
	output, err := cycleTrimmedOutput(ctx, "gh", "pr", "list", "--head", branch, "--base", baseBranch, "--state", "open", "--limit", "1", "--json", "number", "--jq", ".[0].number // empty")
	if err != nil {
		return "", fmt.Errorf("find existing PR for branch %s: %w", branch, err)
	}
	return output, nil
}

func existingOrDefaultTitle(ctx context.Context, explicitTitle, existingPRNumber string) (string, error) {
	if explicitTitle != "" {
		return explicitTitle, nil
	}
	if existingPRNumber != "" {
		title, err := cycleTrimmedOutput(ctx, "gh", "pr", "view", existingPRNumber, "--json", "title", "--jq", ".title")
		if err != nil {
			return "", fmt.Errorf("read existing PR title for #%s: %w", existingPRNumber, err)
		}
		return title, nil
	}
	title, err := cycleTrimmedOutput(ctx, "git", "log", "-1", "--format=%s")
	if err != nil {
		return "", fmt.Errorf("read latest commit title: %w", err)
	}
	return title, nil
}

func buildCommitListMarkdown(ctx context.Context, rangeExpr string) (string, error) {
	output, err := cycleTrimmedOutput(ctx, "git", "log", "--reverse", "--no-merges", "--format=- %h %s", rangeExpr)
	if err != nil {
		return "", fmt.Errorf("build commit list: %w", err)
	}
	if output == "" {
		return "- no commits detected in range", nil
	}
	return output, nil
}

func buildBreakingChangesMarkdown(ctx context.Context, rangeExpr string) (string, error) {
	output, err := cycleTrimmedOutput(ctx, "git", "log", "--reverse", "--no-merges", "--format=%s", rangeExpr)
	if err != nil {
		return "", fmt.Errorf("build breaking changes list: %w", err)
	}

	pattern := regexp.MustCompile(`^[a-z]+(\([^)]+\))?!: (.+)$`)
	seen := map[string]struct{}{}
	var items []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		matches := pattern.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}
		note := strings.TrimSpace(matches[2])
		if note == "" {
			continue
		}
		if _, ok := seen[note]; ok {
			continue
		}
		seen[note] = struct{}{}
		items = append(items, "- "+note)
	}
	if len(items) == 0 {
		return "<!-- None. Remove this section if not needed. -->", nil
	}
	return strings.Join(items, "\n"), nil
}

func buildPRBody(branch, baseBranch string, commitCount int, breakingChanges, commitList string, testPlanItems []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## Summary\n\n- source branch: %s\n- base branch: %s\n- commits in PR: %d\n\n", branch, baseBranch, commitCount)
	b.WriteString("## Breaking Changes\n\n")
	b.WriteString(breakingChanges)
	b.WriteString("\n\n## Included Commits\n\n")
	b.WriteString(commitList)
	b.WriteString("\n\n## Test Plan\n")
	for _, item := range testPlanItems {
		fmt.Fprintf(&b, "- [ ] %s\n", item)
	}
	return strings.TrimRight(b.String(), "\n")
}

func prTestPlanItems(repoRoot string) []string {
	projectType := updater.InferProjectType(repoRoot)
	ov, err := overlay.Get(projectType)
	if err != nil || len(ov.PRTestPlanItems) == 0 {
		return []string{"All validations pass"}
	}
	return ov.PRTestPlanItems
}

func cycleTrimmedOutput(ctx context.Context, name string, args ...string) (string, error) {
	output, err := cycleOutputCommand(ctx, name, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func isGitHubRemote(remoteURL string) bool {
	return strings.Contains(strings.ToLower(remoteURL), "github.com")
}

func commandError(text string) error {
	if text == "" {
		return nil
	}
	return errors.New(text)
}
