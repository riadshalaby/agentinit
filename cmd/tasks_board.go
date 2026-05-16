package cmd

import "strings"

type tasksBoardRow struct {
	TaskID             string
	Scope              string
	Status             string
	AcceptanceCriteria string
	Evidence           string
	NextRole           string
}

func parseTasksBoardRow(line string) (tasksBoardRow, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || !strings.HasPrefix(trimmed, "|") {
		return tasksBoardRow{}, false
	}

	columns := splitMarkdownTableRow(trimmed)
	if len(columns) != 6 {
		return tasksBoardRow{}, false
	}
	if isTasksBoardHeader(columns) || isTasksBoardSeparator(columns) {
		return tasksBoardRow{}, false
	}

	return tasksBoardRow{
		TaskID:             columns[0],
		Scope:              columns[1],
		Status:             columns[2],
		AcceptanceCriteria: columns[3],
		Evidence:           columns[4],
		NextRole:           columns[5],
	}, true
}

func splitMarkdownTableRow(line string) []string {
	var columns []string
	var current strings.Builder
	escaped := false
	for _, r := range line {
		switch {
		case escaped:
			current.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
			current.WriteRune(r)
		case r == '|':
			columns = append(columns, strings.TrimSpace(current.String()))
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}
	columns = append(columns, strings.TrimSpace(current.String()))

	if len(columns) >= 2 && columns[0] == "" && columns[len(columns)-1] == "" {
		return columns[1 : len(columns)-1]
	}
	return columns
}

func isTasksBoardHeader(columns []string) bool {
	return len(columns) == 6 &&
		columns[0] == "Task ID" &&
		columns[1] == "Scope" &&
		columns[2] == "Status" &&
		columns[3] == "Acceptance Criteria" &&
		columns[4] == "Evidence" &&
		columns[5] == "Next Role"
}

func isTasksBoardSeparator(columns []string) bool {
	for _, column := range columns {
		trimmed := strings.TrimSpace(column)
		if trimmed == "" {
			return false
		}
		if !strings.Contains(trimmed, "-") {
			return false
		}
		withoutMarkers := strings.Trim(trimmed, "-: ")
		if withoutMarkers != "" {
			return false
		}
	}
	return true
}
