package overlay

import "github.com/riadshalaby/agentinit/internal/template"

func init() {
	register(Overlay{
		Name: "node",
		ToolPermissions: []string{
			"npm",
			"npx",
			"node",
			"eslint",
			"prettier",
		},
		ValidationCommands: []template.ValidationCommand{
			{Label: "Lint", Command: "npm run lint"},
			{Label: "Build", Command: "npm run build"},
			{Label: "Test", Command: "npm test"},
		},
		PRTestPlanItems: []string{"npm test", "npm run lint"},
	})
}
