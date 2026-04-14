package tradev2

import "testing"

func TestNew_ModuleName(t *testing.T) {
	const want = "test-tradev2"
	m := New(want)
	if m.ModuleName != want {
		t.Fatalf("ModuleName = %q, want %q", m.ModuleName, want)
	}
}
