package overlay

import "github.com/riadshalaby/agentinit/internal/template"

func init() {
	register(Overlay{
		Name: "java",
		ToolPermissions: []string{
			"mvn",
			"gradle",
			"javac",
			"java",
		},
		ValidationCommands: []template.ValidationCommand{
			{Label: "Format", Command: "mvn -q spotless:apply"},
			{Label: "Compile", Command: "mvn -q -DskipTests test-compile"},
			{Label: "Test", Command: "mvn -T 1C -q test"},
		},
		PRTestPlanItems: []string{"mvn test", "spotless:check"},
	})
}
