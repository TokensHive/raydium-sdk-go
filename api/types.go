package api

type ApiV3Token struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	ChainID  int    `json:"chainId"`
	ProgramID string `json:"programId"`
}

type ApiV3TokenRes struct {
	MintList  []ApiV3Token `json:"mintList"`
	Blacklist []string     `json:"blacklist"`
	WhiteList []string     `json:"whiteList"`
}

type FetchPoolParams struct {
	Type     string `json:"type"`
	Sort     string `json:"sort"`
	Order    string `json:"order"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type PoolsAPIResult struct {
	Data any `json:"data"`
}

type PoolKeys struct {
	ID string `json:"id"`
}

type ApiClmmConfigInfo map[string]any
type ApiCpmmConfigInfo map[string]any
type AvailabilityCheckAPI3 map[string]bool
type ApiLaunchConfig map[string]any
