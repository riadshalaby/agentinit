package overlay

import (
	"fmt"

	"github.com/riadshalaby/agentinit/internal/template"
)

// Overlay defines type-specific validation commands and PR test plan items.
type Overlay struct {
	Name               string
	ToolPermissions    []string
	ValidationCommands []template.ValidationCommand
	PRTestPlanItems    []string
}

var registry = map[string]Overlay{}

func register(o Overlay) {
	registry[o.Name] = o
}

// Get returns the overlay for the given name, or the base overlay if name is empty.
// Returns an error if the name is non-empty and not found.
func Get(name string) (Overlay, error) {
	if name == "" {
		return cloneOverlay(baseOverlay), nil
	}
	o, ok := registry[name]
	if !ok {
		return Overlay{}, fmt.Errorf("unknown project type %q (available: go, java, node)", name)
	}
	result := cloneOverlay(o)
	result.ToolPermissions = append(cloneStrings(baseOverlay.ToolPermissions), result.ToolPermissions...)
	return result, nil
}

func cloneOverlay(o Overlay) Overlay {
	o.ToolPermissions = cloneStrings(o.ToolPermissions)
	o.ValidationCommands = append([]template.ValidationCommand(nil), o.ValidationCommands...)
	o.PRTestPlanItems = cloneStrings(o.PRTestPlanItems)
	return o
}

func cloneStrings(values []string) []string {
	return append([]string(nil), values...)
}
