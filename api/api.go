package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"

	"github.com/TokensHive/raydium-sdk-go/common"
	rsolana "github.com/TokensHive/raydium-sdk-go/solana"
)

var (
	poolKeysCache     = map[string]PoolKeys{}
	cacheLaunchConfig = map[rsolana.Cluster][]ApiLaunchConfig{}
	cacheMu           sync.RWMutex
)

func EndlessRetry[T any](name string, interval time.Duration, call func() (T, error)) (T, error) {
	logger := common.CreateLogger("Raydium_Api")
	for {
		value, err := call()
		if err == nil {
			return value, nil
		}
		logger.Error(fmt.Sprintf("Request %s failed, retry after %s: %v", name, interval, err))
		time.Sleep(interval)
	}
}

type API struct {
	Cluster    rsolana.Cluster
	HTTPClient *http.Client
	BaseURL    string
	URLConfigs map[string]string
}

type APIProps struct {
	Cluster    rsolana.Cluster
	Timeout    time.Duration
	URLConfigs map[string]string
}

func NewAPI(props APIProps) *API {
	baseURL := APIURLs["BASE_HOST"]
	if props.Cluster == rsolana.Devnet {
		baseURL = DevAPIURLs["BASE_HOST"]
	}
	if customBase, ok := props.URLConfigs["BASE_HOST"]; ok && customBase != "" {
		baseURL = customBase
	}
	return &API{
		Cluster: props.Cluster,
		HTTPClient: &http.Client{
			Timeout: props.Timeout,
		},
		BaseURL:    baseURL,
		URLConfigs: props.URLConfigs,
	}
}

func (a *API) endpoint(key string) string {
	if v, ok := a.URLConfigs[key]; ok && v != "" {
		return v
	}
	if a.Cluster == rsolana.Devnet {
		if v, ok := DevAPIURLs[key]; ok {
			return v
		}
	}
	return APIURLs[key]
}

func (a *API) get(ctx context.Context, endpoint string, output any) error {
	url := endpoint
	if !strings.HasPrefix(endpoint, "http") {
		url = a.BaseURL + endpoint
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed (%d): %s", resp.StatusCode, string(body))
	}
	if len(body) == 0 {
		return nil
	}
	return json.Unmarshal(body, output)
}

func (a *API) GetClmmConfigs(ctx context.Context) ([]ApiClmmConfigInfo, error) {
	out := []ApiClmmConfigInfo{}
	err := a.get(ctx, a.endpoint("CLMM_CONFIG"), &out)
	return out, err
}

func (a *API) GetCpmmConfigs(ctx context.Context) ([]ApiCpmmConfigInfo, error) {
	out := []ApiCpmmConfigInfo{}
	err := a.get(ctx, a.endpoint("CPMM_CONFIG"), &out)
	return out, err
}

func (a *API) GetChainTimeOffset(ctx context.Context) (map[string]int64, error) {
	out := map[string]int64{}
	err := a.get(ctx, a.endpoint("CHAIN_TIME"), &out)
	return out, err
}

func (a *API) GetTokenList(ctx context.Context) (ApiV3TokenRes, error) {
	out := ApiV3TokenRes{}
	err := a.get(ctx, a.endpoint("TOKEN_LIST"), &out)
	return out, err
}

func (a *API) GetTokenInfo(ctx context.Context, mint []solana.PublicKey) ([]ApiV3Token, error) {
	mints := make([]string, 0, len(mint))
	for _, m := range mint {
		mints = append(mints, m.String())
	}
	path := fmt.Sprintf("%s?mints=%s", a.endpoint("MINT_INFO_ID"), strings.Join(mints, ","))
	out := []ApiV3Token{}
	err := a.get(ctx, path, &out)
	return out, err
}

func (a *API) GetPoolList(ctx context.Context, params FetchPoolParams) (PoolsAPIResult, error) {
	if params.Type == "" {
		params.Type = "all"
	}
	if params.Sort == "" {
		params.Sort = "liquidity"
	}
	if params.Order == "" {
		params.Order = "desc"
	}
	if params.PageSize == 0 {
		params.PageSize = 100
	}
	query := fmt.Sprintf("%s?poolType=%s&poolSortField=%s&sortType=%s&page=%d&pageSize=%d",
		a.endpoint("POOL_LIST"),
		url.QueryEscape(params.Type),
		url.QueryEscape(params.Sort),
		url.QueryEscape(params.Order),
		params.Page,
		params.PageSize,
	)
	out := PoolsAPIResult{}
	err := a.get(ctx, query, &out)
	return out, err
}

func (a *API) FetchPoolByID(ctx context.Context, ids string) ([]map[string]any, error) {
	out := []map[string]any{}
	err := a.get(ctx, a.endpoint("POOL_SEARCH_BY_ID")+"?ids="+ids, &out)
	return out, err
}

func (a *API) FetchPoolKeysByID(ctx context.Context, ids []string) ([]PoolKeys, error) {
	cacheMu.RLock()
	cached := make([]PoolKeys, 0, len(ids))
	ready := make([]string, 0, len(ids))
	for _, id := range ids {
		if v, ok := poolKeysCache[id]; ok {
			cached = append(cached, v)
		} else {
			ready = append(ready, id)
		}
	}
	cacheMu.RUnlock()

	fetched := []PoolKeys{}
	if len(ready) > 0 {
		err := a.get(ctx, a.endpoint("POOL_KEY_BY_ID")+"?ids="+strings.Join(ready, ","), &fetched)
		if err != nil {
			return nil, err
		}
		cacheMu.Lock()
		for _, item := range fetched {
			poolKeysCache[item.ID] = item
		}
		cacheMu.Unlock()
	}
	return append(cached, fetched...), nil
}

func (a *API) FetchAvailabilityStatus(ctx context.Context) (AvailabilityCheckAPI3, error) {
	out := AvailabilityCheckAPI3{}
	err := a.get(ctx, a.endpoint("CHECK_AVAILABILITY"), &out)
	return out, err
}

func (a *API) FetchLaunchConfigs(ctx context.Context) ([]ApiLaunchConfig, error) {
	cacheMu.RLock()
	if configs, ok := cacheLaunchConfig[a.Cluster]; ok && len(configs) > 0 {
		cacheMu.RUnlock()
		return configs, nil
	}
	cacheMu.RUnlock()

	launchURL := "https://launch-mint-v1.raydium.io/main/configs"
	if a.Cluster == rsolana.Devnet {
		launchURL = "https://launch-mint-v1-devnet.raydium.io/main/configs"
	}
	response := struct {
		Data []ApiLaunchConfig `json:"data"`
	}{}
	if err := a.get(ctx, launchURL, &response); err != nil {
		return nil, err
	}
	cacheMu.Lock()
	cacheLaunchConfig[a.Cluster] = response.Data
	cacheMu.Unlock()
	return response.Data, nil
}
