package launcher

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type RoleLaunchOpts struct {
	Role       string
	Agent      string
	Model      string
	Effort     string
	PromptFile string
	RepoRoot   string
	ExtraArgs  []string
}

type execRunner func(name string, args []string, dir string, stdin io.Reader, stdout, stderr io.Writer) error

var runProcess execRunner = defaultRunProcess
var readFile = os.ReadFile

func Launch(opts RoleLaunchOpts) error {
	switch opts.Agent {
	case "claude":
		args := []string{
			"--permission-mode", "acceptEdits",
			"--add-dir", opts.RepoRoot,
		}
		if opts.Model != "" {
			args = append(args, "--model", opts.Model)
		}
		if opts.Effort != "" {
			args = append(args, "--effort", opts.Effort)
		}
		args = append(args, opts.ExtraArgs...)
		args = append(args, "--system-prompt-file", opts.PromptFile)
		return runProcess("claude", args, opts.RepoRoot, os.Stdin, os.Stdout, os.Stderr)
	case "codex":
		prompt, err := readFile(opts.PromptFile)
		if err != nil {
			return fmt.Errorf("read prompt file %q: %w", opts.PromptFile, err)
		}

		args := []string{
			"--sandbox", "workspace-write",
			"-c", "sandbox_workspace_write.network_access=true",
		}
		if opts.Model != "" {
			args = append(args, "-m", opts.Model)
		}
		args = append(args, opts.ExtraArgs...)
		args = append(args, string(prompt))
		return runProcess("codex", args, opts.RepoRoot, os.Stdin, os.Stdout, os.Stderr)
	default:
		return fmt.Errorf("unsupported agent %q", opts.Agent)
	}
}

func defaultRunProcess(name string, args []string, dir string, stdin io.Reader, stdout, stderr io.Writer) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
