package solana

import "testing"

func TestClusterConstants(t *testing.T) {
	if Mainnet != "mainnet" {
		t.Fatalf("mainnet mismatch")
	}
	if Devnet != "devnet" {
		t.Fatalf("devnet mismatch")
	}
}
