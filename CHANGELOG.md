# Changelog

## [1.0.1](https://github.com/TokensHive/raydium-sdk-go/compare/v1.0.0...v1.0.1) (2026-04-14)


### Build System

* **deps:** bump actions/checkout from 4 to 6 ([40bb87c](https://github.com/TokensHive/raydium-sdk-go/commit/40bb87c9a612666f1e851c1998612d0b073034bd))
* **deps:** bump actions/setup-go from 5 to 6 ([cee8654](https://github.com/TokensHive/raydium-sdk-go/commit/cee8654357a5db3ac455dca9dff03d3955e32f3e))
* **deps:** bump softprops/action-gh-release from 2 to 3 ([3308e0d](https://github.com/TokensHive/raydium-sdk-go/commit/3308e0d11041116fa98feb0622ad0163b7e5d4f4))


### CI

* fix release triggering and add live RPC fallback test ([e28e04c](https://github.com/TokensHive/raydium-sdk-go/commit/e28e04c74ef1917b94170fb4545cd31b94f29054))
* force node24 and patch-bump release-please behavior ([eb67063](https://github.com/TokensHive/raydium-sdk-go/commit/eb67063fe2f53197ddceef8c6e2511348393e09e))
* harden release workflows and token handling ([cc139c8](https://github.com/TokensHive/raydium-sdk-go/commit/cc139c80a33625a613b65bf3474013a5369faee0))
* switch release-please workflow to CLI on node24 runtime ([997d047](https://github.com/TokensHive/raydium-sdk-go/commit/997d047cc16d019c0b1da55e3aee72c0b1ea652a))

## 1.0.0 (2026-04-14)


### Features

* add core sdk primitives and utilities ([3bd226d](https://github.com/TokensHive/raydium-sdk-go/commit/3bd226d78cb6ec3397802e645b60000998937e4b))
* add raydium facade, api client, and module scaffolds ([fc8cecf](https://github.com/TokensHive/raydium-sdk-go/commit/fc8cecfcff629c085813f45d44ebbaf7aea98b98))

## Changelog

All notable changes to this project are documented in this file.

## Unreleased
- Initial Go SDK scaffolding and parity-oriented package layout.
- Core `module`, `common`, `api`, `marshmallow`, and `raydium` facade implementation.
- Full unit test suite with per-package `100%` statement coverage enforcement.
- GitHub Actions CI, release workflows, Dependabot, and documentation baseline.
