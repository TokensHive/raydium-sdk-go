package account

import "testing"

func TestNew_ModuleName(t *testing.T) {
	const want = "test-account"
	m := New(want, nil)
	if m.ModuleName != want {
		t.Fatalf("ModuleName = %q, want %q", m.ModuleName, want)
	}
}

func TestResetTokenAccounts(t *testing.T) {
	initial := []TokenAccountData{
		{Mint: "So11111111111111111111111111111111111111112", Amount: "1000"},
		{Mint: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", Amount: "1"},
	}
	m := New("test-account", initial)
	if len(m.TokenAccounts) != len(initial) {
		t.Fatalf("before reset: len(TokenAccounts) = %d, want %d", len(m.TokenAccounts), len(initial))
	}
	m.ResetTokenAccounts()
	if len(m.TokenAccounts) != 0 {
		t.Fatalf("after reset: len(TokenAccounts) = %d, want 0", len(m.TokenAccounts))
	}
}
