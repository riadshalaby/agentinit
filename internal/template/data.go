package template

// ProjectData holds all values available to templates during rendering.
type ProjectData struct {
	ProjectName        string
	ProjectType        string
	Workflow           string
	ValidationCommands []ValidationCommand
	GitignoreExtra     string
	PRTestPlanItems    []string
}

const (
	WorkflowManual = "manual"
	WorkflowAuto   = "auto"
)

// NormalizeWorkflow applies the default workflow when none is specified.
func NormalizeWorkflow(workflow string) string {
	if workflow == "" {
		return WorkflowManual
	}
	return workflow
}

// ValidWorkflow reports whether the workflow is supported.
func ValidWorkflow(workflow string) bool {
	switch NormalizeWorkflow(workflow) {
	case WorkflowManual, WorkflowAuto:
		return true
	default:
		return false
	}
}

// ValidationCommand represents a single validation step shown in CLAUDE.md and PLAN.template.md.
type ValidationCommand struct {
	Label   string
	Command string
}
