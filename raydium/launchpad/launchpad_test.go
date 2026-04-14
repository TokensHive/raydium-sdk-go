package launchpad

import "testing"

func TestNew_ModuleName(t *testing.T) {
	const want = "test-launchpad"
	m := New(want)
	if m.ModuleName != want {
		t.Fatalf("ModuleName = %q, want %q", m.ModuleName, want)
	}
}
