package serum

import "testing"

func TestNew_ModuleName(t *testing.T) {
	const want = "test-serum"
	m := New(want)
	if m.ModuleName != want {
		t.Fatalf("ModuleName = %q, want %q", m.ModuleName, want)
	}
}
