package raydium

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"

	"github.com/TokensHive/raydium-sdk-go/api"
	"github.com/TokensHive/raydium-sdk-go/common"
	"github.com/TokensHive/raydium-sdk-go/raydium/account"
	"github.com/TokensHive/raydium-sdk-go/raydium/clmm"
	"github.com/TokensHive/raydium-sdk-go/raydium/cpmm"
	"github.com/TokensHive/raydium-sdk-go/raydium/farm"
	"github.com/TokensHive/raydium-sdk-go/raydium/ido"
	"github.com/TokensHive/raydium-sdk-go/raydium/launchpad"
	"github.com/TokensHive/raydium-sdk-go/raydium/liquidity"
	"github.com/TokensHive/raydium-sdk-go/raydium/marketv2"
	"github.com/TokensHive/raydium-sdk-go/raydium/serum"
	tokenmodule "github.com/TokensHive/raydium-sdk-go/raydium/token"
	"github.com/TokensHive/raydium-sdk-go/raydium/tradev2"
	"github.com/TokensHive/raydium-sdk-go/raydium/utils1216"
	rsolana "github.com/TokensHive/raydium-sdk-go/solana"
)

type Connection interface {
	GetEpochInfo(ctx context.Context) (any, error)
}

type RaydiumLoadParams struct {
	Connection          Connection
	Cluster             rsolana.Cluster
	Owner               *solana.PrivateKey
	APIRequestTimeout   time.Duration
	APICacheTime        time.Duration
	URLConfigs          map[string]string
	DisableFeatureCheck bool
	DisableLoadToken    bool
}

type Raydium struct {
	Cluster rsolana.Cluster

	Farm      *farm.Module
	Account   *account.Module
	Liquidity *liquidity.Module
	Clmm      *clmm.Module
	Cpmm      *cpmm.Module
	TradeV2   *tradev2.Module
	Utils1216 *utils1216.Module
	MarketV2  *marketv2.Module
	Ido       *ido.Module
	Token     *tokenmodule.Module
	Launchpad *launchpad.Module
	Serum     *serum.Module

	Connection   Connection
	API          *api.API
	Owner        *solana.PrivateKey
	APICacheTime time.Duration
	Availability api.AvailabilityCheckAPI3
}

func New(config RaydiumLoadParams) *Raydium {
	cluster := config.Cluster
	if cluster == "" {
		cluster = rsolana.Mainnet
	}
	timeout := config.APIRequestTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	cache := config.APICacheTime
	if cache <= 0 {
		cache = 5 * time.Minute
	}
	client := api.NewAPI(api.APIProps{
		Cluster:    cluster,
		Timeout:    timeout,
		URLConfigs: config.URLConfigs,
	})
	return &Raydium{
		Cluster:      cluster,
		Connection:   config.Connection,
		API:          client,
		Owner:        config.Owner,
		APICacheTime: cache,
		Farm:         farm.New("Raydium_Farm"),
		Account:      account.New("Raydium_Account", nil),
		Liquidity:    liquidity.New("Raydium_LiquidityV2"),
		Clmm:         clmm.New("Raydium_clmm"),
		Cpmm:         cpmm.New("Raydium_cpmm"),
		TradeV2:      tradev2.New("Raydium_tradeV2"),
		Utils1216:    utils1216.New("Raydium_utils1216"),
		MarketV2:     marketv2.New("Raydium_marketV2"),
		Ido:          ido.New("Raydium_ido"),
		Token:        tokenmodule.New("Raydium_tokenV2"),
		Launchpad:    launchpad.New("Raydium_launchpad"),
		Serum:        serum.New("Raydium_serum"),
		Availability: api.AvailabilityCheckAPI3{},
	}
}

func Load(config RaydiumLoadParams) (*Raydium, error) {
	instance := New(config)
	if !config.DisableFeatureCheck {
		availability, err := instance.API.FetchAvailabilityStatus(context.Background())
		if err != nil {
			return nil, err
		}
		instance.Availability = availability
	}
	return instance, nil
}

func (r *Raydium) CheckOwner() error {
	if r.Owner == nil {
		return fmt.Errorf(common.EmptyOwner)
	}
	return nil
}

func (r *Raydium) SetOwner(owner *solana.PrivateKey) {
	r.Owner = owner
	r.Account.ResetTokenAccounts()
}
