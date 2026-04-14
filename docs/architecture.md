# Architecture

## Package Topology
- `api`: Raydium HTTP API client and endpoint models.
- `common`: shared constants, numeric parsing, logging, and Solana key/program helpers.
- `module`: amount/fraction/price/currency/token primitives.
- `marshmallow`: binary layout encode/decode helpers.
- `solana`: SDK cluster definitions.
- `raydium`: façade and feature modules (`account`, `farm`, `liquidity`, `clmm`, `cpmm`, `tradev2`, `utils1216`, `marketv2`, `ido`, `launchpad`, `serum`, `token`).

## Testing Strategy
- Unit tests for all code paths (including error branches).
- Deterministic API behavior tests using `httptest`.
- Package-by-package coverage gate via `scripts/coverage_gate.sh`.

## Release Strategy
- CI gate on every PR/push.
- Automated release PRs with Release Please.
- Tagged releases publish GitHub release notes.
