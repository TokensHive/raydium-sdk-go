package cpmm

import "testing"

func TestNew_ModuleName(t *testing.T) {
	const want = "test-cpmm"
	m := New(want)
	if m.ModuleName != want {
		t.Fatalf("ModuleName = %q, want %q", m.ModuleName, want)
	}
}
