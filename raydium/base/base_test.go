package base

import "testing"

func TestNewModuleBase(t *testing.T) {
	module := NewModuleBase("test")
	if module.ModuleName != "test" {
		t.Fatalf("module name mismatch")
	}
	if module.Logger == nil {
		t.Fatalf("logger should be set")
	}
}
