package update

import "testing"

func TestExtractSectionsSplitsMarkerFile(t *testing.T) {
	before, managed, after, err := ExtractSections("preamble\n" + markerStart + "\nmanaged\n" + markerEnd + "\nafter\n")
	if err != nil {
		t.Fatalf("ExtractSections() error = %v", err)
	}
	if before != "preamble\n" {
		t.Fatalf("before = %q", before)
	}
	if managed != "managed" {
		t.Fatalf("managed = %q", managed)
	}
	if after != "\nafter\n" {
		t.Fatalf("after = %q", after)
	}
}

func TestReplaceManagedSectionPreservesUserContent(t *testing.T) {
	updated, err := ReplaceManagedSection("before\n"+markerStart+"\nold\n"+markerEnd+"\nafter\n", "new")
	if err != nil {
		t.Fatalf("ReplaceManagedSection() error = %v", err)
	}
	want := "before\n" + markerStart + "\nnew\n" + markerEnd + "\nafter\n"
	if updated != want {
		t.Fatalf("updated = %q, want %q", updated, want)
	}
}

func TestReplaceManagedSectionPrependsWhenMarkersMissing(t *testing.T) {
	updated, err := ReplaceManagedSection("user content\n", "managed")
	if err != nil {
		t.Fatalf("ReplaceManagedSection() error = %v", err)
	}
	want := markerStart + "\nmanaged\n" + markerEnd + "\nuser content\n"
	if updated != want {
		t.Fatalf("updated = %q, want %q", updated, want)
	}
}
