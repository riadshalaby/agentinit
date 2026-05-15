package cmd

import "github.com/spf13/cobra"

var planCmd = &cobra.Command{
	Use:   "plan [claude|codex] [agent-options...]",
	Short: "Launch the planner role session",
	Long: `Start the planner session for roadmap refinement and task planning.

Use this at the beginning of a cycle to turn ` + "`ROADMAP.md`" + ` into
` + " `.ai/PLAN.md` and `.ai/TASKS.md`." + `
The planner is shared by both full and lite profiles.`,
	Example: "aide plan\naide plan claude --model sonnet",
	Args:    cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRoleLaunch("plan", "planner.md", "claude", args)
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
