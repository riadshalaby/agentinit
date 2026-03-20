package overlay

import (
	"testing"
)

func TestGetBase(t *testing.T) {
	o, err := Get("")
	if err != nil {
		t.Fatalf("Get(\"\") error: %v", err)
	}
	if len(o.ValidationCommands) != 0 {
		t.Errorf("base overlay should have no validation commands, got %d", len(o.ValidationCommands))
	}
	if len(o.PRTestPlanItems) != 1 || o.PRTestPlanItems[0] != "All validations pass" {
		t.Errorf("base overlay PRTestPlanItems = %v, want [All validations pass]", o.PRTestPlanItems)
	}
}

func TestGetGo(t *testing.T) {
	o, err := Get("go")
	if err != nil {
		t.Fatalf("Get(\"go\") error: %v", err)
	}
	if len(o.ValidationCommands) != 3 {
		t.Errorf("go overlay should have 3 validation commands, got %d", len(o.ValidationCommands))
	}
	if o.ValidationCommands[0].Command != "go fmt ./..." {
		t.Errorf("first command = %q, want \"go fmt ./...\"", o.ValidationCommands[0].Command)
	}
}

func TestGetJava(t *testing.T) {
	o, err := Get("java")
	if err != nil {
		t.Fatalf("Get(\"java\") error: %v", err)
	}
	if len(o.ValidationCommands) != 3 {
		t.Errorf("java overlay should have 3 validation commands, got %d", len(o.ValidationCommands))
	}
}

func TestGetNode(t *testing.T) {
	o, err := Get("node")
	if err != nil {
		t.Fatalf("Get(\"node\") error: %v", err)
	}
	if len(o.ValidationCommands) != 3 {
		t.Errorf("node overlay should have 3 validation commands, got %d", len(o.ValidationCommands))
	}
}

func TestGetUnknown(t *testing.T) {
	_, err := Get("python")
	if err == nil {
		t.Error("Get(\"python\") should return error")
	}
}
