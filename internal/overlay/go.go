package overlay

import "github.com/riadshalaby/agentinit/internal/template"

func init() {
	register(Overlay{
		Name: "go",
		ToolPermissions: []string{
			"go fmt",
			"go vet",
			"go test",
			"go build",
			"go run",
			"go mod",
		},
		ValidationCommands: []template.ValidationCommand{
			{Label: "Format", Command: "go fmt ./..."},
			{Label: "Vet", Command: "go vet ./..."},
			{Label: "Test", Command: "go test ./..."},
		},
		PRTestPlanItems: []string{"go test", "go vet"},
	})
}
