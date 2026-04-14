# raydium-sdk-go

Go SDK for Raydium with parity-oriented package layout against Raydium SDK V2 TypeScript baseline.

## Status
- Canonical TS source baseline: `/home/hamza/Web_Kitchen/TokensHive.com/raydium-sdk-V2` (`main_head`).
- Per-package tests are enforced with a strict `100.0%` statement coverage gate.
- CI, release automation, and dependency update automation are included.

## Install
```bash
go get github.com/TokensHive/raydium-sdk-go@latest
```

## Go Version
- Go `1.25+`

## Package Overview
- `api`: Raydium HTTP client and endpoint models.
- `common`: program IDs, key helpers, logging, and numeric parsing utilities.
- `module`: core math and amount primitives (`Fraction`, `TokenAmount`, `CurrencyAmount`, `Price`, `Percent`).
- `marshmallow`: binary serialization helpers.
- `solana`: cluster abstractions.
- `raydium`: root facade and feature-module wiring.
- `raydium/account`, `raydium/farm`, `raydium/liquidity`, `raydium/clmm`, `raydium/cpmm`, `raydium/tradev2`, `raydium/utils1216`, `raydium/marketv2`, `raydium/ido`, `raydium/launchpad`, `raydium/serum`, `raydium/token`: feature module packages.

## Quickstart
```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/TokensHive/raydium-sdk-go/raydium"
)

func main() {
	sdk, err := raydium.Load(raydium.RaydiumLoadParams{
		APIRequestTimeout: 10 * time.Second,
		DisableFeatureCheck: false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("modules ready:", sdk.Farm.ModuleName, sdk.Clmm.ModuleName)
	_ = context.Background()
}
```

## Development
```bash
go mod tidy
go test ./... -cover
./scripts/coverage_gate.sh
```

## CI/CD
- `ci.yml`: format, vet, tests, coverage gate.
- `live-network.yml`: scheduled/non-blocking network verification.
- `release-please.yml`: automated release PRs and versioning.
- `release.yml`: release publication on `v*` tags.
- `dependabot.yml`: automated `gomod` and GitHub Actions updates.

## Documentation
- Architecture: `docs/architecture.md`
- TS parity mapping: `docs/parity-matrix.md`
- Contribution guide: `CONTRIBUTING.md`
- Security policy: `SECURITY.md`

## Coverage Policy
- Every package containing executable statements must remain at `100.0%` statement coverage.
- Enforced in CI via `scripts/coverage_gate.sh`.

## License
- GNU GPL v3.0

This Go SDK is inspired by the official Raydium TypeScript SDK, and full credit goes to the Raydium team for the original design and implementation.