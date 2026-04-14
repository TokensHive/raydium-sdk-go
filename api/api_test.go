package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"

	rsolana "github.com/TokensHive/raydium-sdk-go/solana"
)

type rewriteTransport struct {
	target *url.URL
	base   http.RoundTripper
}

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = t.target.Scheme
	req.URL.Host = t.target.Host
	req.Host = t.target.Host
	return t.base.RoundTrip(req)
}

type errorBody struct{}

func (errorBody) Read(_ []byte) (int, error) { return 0, errors.New("read error") }
func (errorBody) Close() error               { return nil }

type readErrorTransport struct{}

func (readErrorTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(errorBody{}),
		Header:     http.Header{},
	}, nil
}

type requestErrorTransport struct{}

func (requestErrorTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, errors.New("request failed")
}

func TestEndlessRetry(t *testing.T) {
	attempt := 0
	result, err := EndlessRetry("test", time.Millisecond, func() (string, error) {
		attempt++
		if attempt == 1 {
			return "", context.DeadlineExceeded
		}
		return "ok", nil
	})
	if err != nil || result != "ok" || attempt != 2 {
		t.Fatalf("unexpected endless retry result: %v %v %d", result, err, attempt)
	}
}

func TestAPIRequestsAndCache(t *testing.T) {
	cacheMu.Lock()
	poolKeysCache = map[string]PoolKeys{}
	cacheLaunchConfig = map[rsolana.Cluster][]ApiLaunchConfig{}
	cacheMu.Unlock()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/main/clmm-config":
			_ = json.NewEncoder(w).Encode([]ApiClmmConfigInfo{{"a": "b"}})
		case r.URL.Path == "/main/cpmm-config":
			_ = json.NewEncoder(w).Encode([]ApiCpmmConfigInfo{{"c": "d"}})
		case r.URL.Path == "/main/chain-time":
			_ = json.NewEncoder(w).Encode(map[string]int64{"offset": 1})
		case r.URL.Path == "/mint/list":
			_ = json.NewEncoder(w).Encode(ApiV3TokenRes{MintList: []ApiV3Token{{Address: "a"}}})
		case r.URL.Path == "/mint/ids":
			_ = json.NewEncoder(w).Encode([]ApiV3Token{{Address: "x"}})
		case r.URL.Path == "/pools/info/list":
			_ = json.NewEncoder(w).Encode(PoolsAPIResult{Data: "ok"})
		case r.URL.Path == "/pools/info/ids":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "p1"}})
		case r.URL.Path == "/pools/key/ids":
			_ = json.NewEncoder(w).Encode([]PoolKeys{{ID: "p1"}})
		case r.URL.Path == "/v3/main/AvailabilityCheckAPI":
			_ = json.NewEncoder(w).Encode(AvailabilityCheckAPI3{"all": true})
		case r.URL.Path == "/main/configs":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []ApiLaunchConfig{{"id": "cfg"}}})
		case r.URL.Path == "/error":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"bad"}`))
		case r.URL.Path == "/empty":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewAPI(APIProps{
		Cluster: rsolana.Mainnet,
		Timeout: time.Second,
		URLConfigs: map[string]string{
			"BASE_HOST":          server.URL,
			"CLMM_CONFIG":        "/main/clmm-config",
			"CPMM_CONFIG":        "/main/cpmm-config",
			"CHAIN_TIME":         "/main/chain-time",
			"TOKEN_LIST":         "/mint/list",
			"MINT_INFO_ID":       "/mint/ids",
			"POOL_LIST":          "/pools/info/list",
			"POOL_SEARCH_BY_ID":  "/pools/info/ids",
			"POOL_KEY_BY_ID":     "/pools/key/ids",
			"CHECK_AVAILABILITY": "/v3/main/AvailabilityCheckAPI",
		},
	})
	targetURL, _ := url.Parse(server.URL)
	client.HTTPClient.Transport = rewriteTransport{target: targetURL, base: http.DefaultTransport}

	ctx := context.Background()
	if _, err := client.GetClmmConfigs(ctx); err != nil {
		t.Fatalf("clmm config error: %v", err)
	}
	if _, err := client.GetCpmmConfigs(ctx); err != nil {
		t.Fatalf("cpmm config error: %v", err)
	}
	if _, err := client.GetChainTimeOffset(ctx); err != nil {
		t.Fatalf("chain time error: %v", err)
	}
	if _, err := client.GetTokenList(ctx); err != nil {
		t.Fatalf("token list error: %v", err)
	}
	if _, err := client.GetTokenInfo(ctx, []solana.PublicKey{solana.PublicKey{}}); err != nil {
		t.Fatalf("token info error: %v", err)
	}
	if _, err := client.GetPoolList(ctx, FetchPoolParams{}); err != nil {
		t.Fatalf("pool list error: %v", err)
	}
	if _, err := client.FetchPoolByID(ctx, "p1"); err != nil {
		t.Fatalf("pool by id error: %v", err)
	}
	keys, err := client.FetchPoolKeysByID(ctx, []string{"p1"})
	if err != nil || len(keys) != 1 {
		t.Fatalf("pool keys fetch error: %v", err)
	}
	// cache branch
	if _, err := client.FetchPoolKeysByID(ctx, []string{"p1"}); err != nil {
		t.Fatalf("pool key cache fetch error: %v", err)
	}
	if _, err := client.FetchAvailabilityStatus(ctx); err != nil {
		t.Fatalf("availability fetch error: %v", err)
	}
	launch, err := client.FetchLaunchConfigs(ctx)
	if err != nil || len(launch) != 1 {
		t.Fatalf("launch config error: %v", err)
	}
	// cache branch
	if _, err := client.FetchLaunchConfigs(ctx); err != nil {
		t.Fatalf("launch config cache error: %v", err)
	}
	if _, err := client.FetchPoolByID(ctx, ""); err != nil {
		t.Fatalf("empty id should still call endpoint: %v", err)
	}

	if err := client.get(ctx, "/error", &map[string]any{}); err == nil {
		t.Fatalf("expected get status error")
	}
	if err := client.get(ctx, "/empty", &map[string]any{}); err != nil {
		t.Fatalf("empty body path should succeed: %v", err)
	}
}

func TestEndpointResolution(t *testing.T) {
	mainnet := NewAPI(APIProps{
		Cluster: rsolana.Mainnet,
		Timeout: time.Second,
	})
	if mainnet.endpoint("BASE_HOST") != APIURLs["BASE_HOST"] {
		t.Fatalf("mainnet endpoint mismatch")
	}
	devnet := NewAPI(APIProps{
		Cluster: rsolana.Devnet,
		Timeout: time.Second,
	})
	if devnet.endpoint("SWAP_HOST") != DevAPIURLs["SWAP_HOST"] {
		t.Fatalf("devnet endpoint mismatch")
	}
	custom := NewAPI(APIProps{
		Cluster: rsolana.Devnet,
		Timeout: time.Second,
		URLConfigs: map[string]string{
			"SWAP_HOST": "https://custom",
		},
	})
	if custom.endpoint("SWAP_HOST") != "https://custom" {
		t.Fatalf("custom endpoint mismatch")
	}
}

func TestAPIErrors(t *testing.T) {
	client := NewAPI(APIProps{
		Cluster: rsolana.Mainnet,
		Timeout: time.Second,
	})
	if err := client.get(context.Background(), "http://%", &map[string]any{}); err == nil {
		t.Fatalf("expected URL parse error")
	}

	client.HTTPClient.Transport = readErrorTransport{}
	if err := client.get(context.Background(), "/main/clmm-config", &map[string]any{}); err == nil {
		t.Fatalf("expected read body error")
	}
	client.HTTPClient.Transport = requestErrorTransport{}
	if err := client.get(context.Background(), "/main/clmm-config", &map[string]any{}); err == nil {
		t.Fatalf("expected request transport error")
	}

	cacheMu.Lock()
	poolKeysCache = map[string]PoolKeys{}
	cacheLaunchConfig = map[rsolana.Cluster][]ApiLaunchConfig{}
	cacheMu.Unlock()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	}))
	defer server.Close()

	errorClient := NewAPI(APIProps{
		Cluster: rsolana.Devnet,
		Timeout: time.Second,
		URLConfigs: map[string]string{
			"BASE_HOST":      server.URL,
			"POOL_KEY_BY_ID": "/pools/key/ids",
		},
	})
	targetURL, _ := url.Parse(server.URL)
	errorClient.HTTPClient.Transport = rewriteTransport{target: targetURL, base: http.DefaultTransport}
	if _, err := errorClient.FetchPoolKeysByID(context.Background(), []string{"x"}); err == nil {
		t.Fatalf("expected pool keys fetch error")
	}
	if _, err := errorClient.FetchLaunchConfigs(context.Background()); err == nil {
		t.Fatalf("expected launch config fetch error")
	}
}
