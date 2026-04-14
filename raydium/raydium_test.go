package raydium

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"

	"github.com/TokensHive/raydium-sdk-go/raydium/account"
	rsolana "github.com/TokensHive/raydium-sdk-go/solana"
)

type mockConnection struct{}

func (m mockConnection) GetEpochInfo(ctx context.Context) (any, error) {
	return map[string]string{"epoch": "1"}, nil
}

func TestNewRaydiumDefaults(t *testing.T) {
	instance := New(RaydiumLoadParams{
		Connection: mockConnection{},
	})
	if instance.Cluster != rsolana.Mainnet {
		t.Fatalf("default cluster mismatch")
	}
	if instance.API == nil || instance.Farm == nil || instance.Account == nil {
		t.Fatalf("raydium modules should be initialized")
	}
	if instance.APICacheTime != 5*time.Minute {
		t.Fatalf("default cache mismatch")
	}
}

func TestLoadAndOwnerHelpers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/main/AvailabilityCheckAPI" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]bool{"all": true})
	}))
	defer server.Close()

	loaded, err := Load(RaydiumLoadParams{
		Connection: mockConnection{},
		URLConfigs: map[string]string{
			"BASE_HOST":          server.URL,
			"CHECK_AVAILABILITY": "/v3/main/AvailabilityCheckAPI",
		},
	})
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if !loaded.Availability["all"] {
		t.Fatalf("availability should be populated")
	}

	skipped, err := Load(RaydiumLoadParams{
		Connection:          mockConnection{},
		DisableFeatureCheck: true,
	})
	if err != nil {
		t.Fatalf("load with skip check failed: %v", err)
	}
	if len(skipped.Availability) != 0 {
		t.Fatalf("availability should remain empty when skipped")
	}
	if err := skipped.CheckOwner(); err == nil {
		t.Fatalf("check owner should fail without owner")
	}
	privateKey := solana.NewWallet().PrivateKey
	skipped.Account = account.New("acct", []account.TokenAccountData{{Mint: "m", Amount: "1"}})
	skipped.SetOwner(&privateKey)
	if skipped.Owner == nil {
		t.Fatalf("owner should be set")
	}
	if len(skipped.Account.TokenAccounts) != 0 {
		t.Fatalf("set owner should reset token accounts")
	}
	if err := skipped.CheckOwner(); err != nil {
		t.Fatalf("check owner should pass after setting owner")
	}

	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	}))
	defer errorServer.Close()
	if _, err := Load(RaydiumLoadParams{
		Connection: mockConnection{},
		URLConfigs: map[string]string{
			"BASE_HOST":          errorServer.URL,
			"CHECK_AVAILABILITY": "/v3/main/AvailabilityCheckAPI",
		},
	}); err == nil {
		t.Fatalf("expected load availability error")
	}
}
