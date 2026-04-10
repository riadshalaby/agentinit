package update

import (
	"fmt"
	"strings"
)

const (
	markerStart = "<!-- agentinit:managed:start -->"
	markerEnd   = "<!-- agentinit:managed:end -->"
)

func ExtractSections(content string) (before, managed, after string, err error) {
	start := strings.Index(content, markerStart)
	end := strings.Index(content, markerEnd)
	if start == -1 || end == -1 {
		return "", "", "", fmt.Errorf("managed markers not found")
	}
	if end < start {
		return "", "", "", fmt.Errorf("managed marker end appears before start")
	}

	before = content[:start]
	managedStart := start + len(markerStart)
	managed = strings.Trim(content[managedStart:end], "\n")
	after = content[end+len(markerEnd):]
	return before, managed, after, nil
}

func ReplaceManagedSection(existing, newManaged string) (string, error) {
	trimmedManaged := strings.Trim(newManaged, "\n")
	before, _, after, err := ExtractSections(existing)
	if err == nil {
		return before + markerStart + "\n" + trimmedManaged + "\n" + markerEnd + after, nil
	}

	existing = strings.TrimLeft(existing, "\n")
	if existing == "" {
		return markerStart + "\n" + trimmedManaged + "\n" + markerEnd + "\n", nil
	}

	return markerStart + "\n" + trimmedManaged + "\n" + markerEnd + "\n" + existing, nil
}
