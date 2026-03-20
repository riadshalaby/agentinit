package template

// ProjectData holds all values available to templates during rendering.
type ProjectData struct {
	ProjectName        string
	ProjectType        string
	ValidationCommands []ValidationCommand
	GitignoreExtra     string
	PRTestPlanItems    []string
}

// ValidationCommand represents a single validation step shown in CLAUDE.md and PLAN.template.md.
type ValidationCommand struct {
	Label   string
	Command string
}
