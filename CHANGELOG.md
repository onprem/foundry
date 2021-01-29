# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

NOTE: As semantic versioning states all 0.y.z releases can contain breaking changes in API (flags, gRPC API, any backward compatibility)

We use _breaking :warning:_ to mark changes that are not backward compatible (relates only to v0.y.z releases.)

## Unreleased

### Added
- [#7](https://github.com/prmsrswt/foundry/pull/7) Furnace: Add gRPC protobuf. Add `foundry furnace` sub-command to run the Furnace component.
- [#12](https://github.com/prmsrswt/foundry/pull/12) Furnace: Implement a package builder based on `makepkg`.
- [#15](https://github.com/prmsrswt/foundry/pull/15) Add internal endpoints (healthchecks, metrics, pprof) and gracefully shutdown everything.
- [#17](https://github.com/prmsrswt/foundry/pull/17) Furnace: Instrument makepkg builder.
### Fixed

### Changed
