package cmd

import "testing"

func TestParseTasksBoardRowHandlesEscapedPipes(t *testing.T) {
	row, ok := parseTasksBoardRow("| T-005 | `--profile lite\\|full` flag | done | ok | pass | none |")
	if !ok {
		t.Fatal("parseTasksBoardRow() = false, want true")
	}
	if row.TaskID != "T-005" {
		t.Fatalf("TaskID = %q, want %q", row.TaskID, "T-005")
	}
	if row.Scope != "`--profile lite\\|full` flag" {
		t.Fatalf("Scope = %q", row.Scope)
	}
	if row.Status != "done" {
		t.Fatalf("Status = %q, want done", row.Status)
	}
}

func TestParseTasksBoardRowHandlesPipeAdjacentText(t *testing.T) {
	row, ok := parseTasksBoardRow("| T-002 | mid-cycle switch shows a one-line advisory for `full\\|lite` transitions | ready_for_review | ok | pending | review |")
	if !ok {
		t.Fatal("parseTasksBoardRow() = false, want true")
	}
	if row.Status != "ready_for_review" {
		t.Fatalf("Status = %q, want %q", row.Status, "ready_for_review")
	}
}

func TestParseTasksBoardRowFiltersSeparatorRow(t *testing.T) {
	if _, ok := parseTasksBoardRow("| --- | --- | --- | --- | --- | --- |"); ok {
		t.Fatal("parseTasksBoardRow() = true, want false")
	}
}

func TestParseTasksBoardRowFiltersHeaderRow(t *testing.T) {
	if _, ok := parseTasksBoardRow("| Task ID | Scope | Status | Acceptance Criteria | Evidence | Next Role |"); ok {
		t.Fatal("parseTasksBoardRow() = true, want false")
	}
}

func TestParseTasksBoardRowHandlesEmptyCells(t *testing.T) {
	row, ok := parseTasksBoardRow("| T-009 | scope | done |  |  | none |")
	if !ok {
		t.Fatal("parseTasksBoardRow() = false, want true")
	}
	if row.AcceptanceCriteria != "" {
		t.Fatalf("AcceptanceCriteria = %q, want empty", row.AcceptanceCriteria)
	}
	if row.Evidence != "" {
		t.Fatalf("Evidence = %q, want empty", row.Evidence)
	}
}
