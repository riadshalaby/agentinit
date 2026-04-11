package overlay

var baseOverlay = Overlay{
	Name:               "",
	ToolPermissions:    []string{"gh", "rg", "fd", "bat", "jq", "sg", "fzf", "tree-sitter"},
	ValidationCommands: nil,
	PRTestPlanItems:    []string{"All validations pass"},
}
