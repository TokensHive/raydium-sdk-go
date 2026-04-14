# Contributing

## Development Setup
- Install Go `1.25+`.
- Clone the repository and run:
  - `go mod tidy`
  - `go test ./... -cover`
  - `./scripts/coverage_gate.sh`

## Quality Bar
- Every package with executable statements must keep `100.0%` statement coverage.
- Keep behavior parity with the TypeScript SDK baseline located at:
  - `/home/hamza/Web_Kitchen/TokensHive.com/raydium-sdk-V2`
- Add/update tests for each behavior change.

## Pull Requests
- Keep commits focused and small.
- Include test evidence (`go test ./... -cover` output).
- Update `README.md` and docs for user-facing changes.
- Include changelog context in PR body (release automation consumes commit history).
