//go:build live

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

type epochInfoRPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  any             `json:"error"`
}

func TestLiveRPCGetEpochInfo(t *testing.T) {
	endpoint := os.Getenv("SOLANA_RPC_HTTP")
	if endpoint == "" {
		endpoint = "https://api.mainnet-beta.solana.com"
	}

	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "getEpochInfo",
		"params":  []any{},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal rpc payload: %v", err)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create rpc request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("rpc request failed against %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("rpc request failed with status %d", resp.StatusCode)
	}

	var decoded epochInfoRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode rpc response: %v", err)
	}
	if decoded.Error != nil {
		t.Fatalf("rpc returned error: %v", decoded.Error)
	}
	if len(decoded.Result) == 0 || string(decoded.Result) == "null" {
		t.Fatalf("rpc returned empty epoch info result")
	}
}
