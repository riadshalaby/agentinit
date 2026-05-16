package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	agentmcp "github.com/riadshalaby/agentinit/internal/mcp"
	"github.com/spf13/cobra"
)

var profileCmd = newProfileCmd()

func newProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Show or switch the active workflow profile",
		Long: `Show or change the workflow profile stored in .ai/config.json.

Profiles control which role sessions and commands the project should use:
- full: planner + implementer + reviewer + optional PO orchestration
- lite: planner + dev session with the human approval stop at ready_for_review

Use "aide profile" to inspect the active mode, or "aide profile full|lite" to switch.`,
		Example: "aide profile\naide profile lite\naide profile full --help",
		Args:    cobra.NoArgs,
		RunE:    runProfileShow,
	}

	cmd.AddCommand(newProfileModeCmd("full"), newProfileModeCmd("lite"))
	return cmd
}

func newProfileModeCmd(profile string) *cobra.Command {
	modeDetails := map[string]struct {
		short   string
		long    string
		example string
	}{
		"full": {
			short: "Switch to the full multi-session workflow profile",
			long: `Switch the workflow profile to full mode.

Full mode uses separate role sessions:
- ` + "`aide plan`" + ` writes the plan and task board.
- ` + "`aide implement`" + ` handles implementation work and ` + "`commit_task`" + ` after review approval.
- ` + "`aide review`" + ` owns ` + "`in_review`" + ` and moves approved tasks to ` + "`ready_to_commit`" + `.
- ` + "`aide po`" + ` is available for orchestration when you want the auto loop.

Use this mode when you want distinct role ownership and optional PO automation.`,
			example: "aide profile full",
		},
		"lite": {
			short: "Switch to the lite planner + dev workflow profile",
			long: `Switch the workflow profile to lite mode.

Lite mode keeps planning separate and then uses one dev session:
- ` + "`aide plan`" + ` writes the plan and task board.
- ` + "`aide dev`" + ` implements, reviews, and prepares commits in one flow.
- ` + "`ready_for_review`" + ` is the human checkpoint before approval.
- The human approves by typing ` + "`commit_task`" + ` after the dev session halts.
- ` + "`aide implement`" + `, ` + "`aide review`" + `, and ` + "`aide po`" + ` are not used in lite mode.

Use this mode when you want fewer long-lived sessions and direct human approval at the review stop.`,
			example: "aide profile lite",
		},
	}

	details := modeDetails[profile]
	return &cobra.Command{
		Use:     profile,
		Short:   details.short,
		Long:    details.long,
		Example: details.example,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProfileSet(cmd, profile)
		},
	}
}

func runProfileShow(cmd *cobra.Command, args []string) error {
	cwd, err := getWorkingDir()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	cfg, err := loadLaunchConfig(cwd)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "Active profile: %s\nConfig: %s\n", cfg.ActiveProfile(), profileConfigPath(cwd))
	return err
}

func runProfileSet(cmd *cobra.Command, profile string) error {
	cwd, err := getWorkingDir()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	if err := agentmcp.WriteProfile(cwd, profile); err != nil {
		return err
	}

	if hasInFlight, err := hasInFlightTasks(cwd); err != nil {
		return err
	} else if hasInFlight {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Advisory: task board has in-flight work; switching profiles changes which commands and roles should continue the cycle."); err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "Updated profile to %s in %s\n", profile, profileConfigPath(cwd))
	return err
}

func hasInFlightTasks(cwd string) (bool, error) {
	data, err := os.ReadFile(filepath.Join(cwd, ".ai", "TASKS.md"))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read tasks: %w", err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		row, ok := parseTasksBoardRow(line)
		if !ok {
			continue
		}
		if row.Status != "done" {
			return true, nil
		}
	}

	return false, nil
}

func profileConfigPath(cwd string) string {
	return filepath.Join(cwd, ".ai", "config.json")
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
