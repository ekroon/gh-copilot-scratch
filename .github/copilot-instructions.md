# Copilot Instructions

## Build, test, and lint

```bash
go build -o gh-copilot-scratch ./cmd/gh-copilot-scratch   # build
go vet ./...                                                # lint
go test -race ./...                                         # unit tests
go test -race -tags=integration ./...                       # all tests (including integration)
```

## Architecture

Single Go binary (`gh-copilot-scratch`) that serves as a `gh` CLI extension.

### Packages

- `cmd/gh-copilot-scratch/` — Entry point. Parses flags, orchestrates scratch dir creation + copilot exec.
- `internal/scratch/` — XDG-compliant scratch directory creation with timestamped naming and git init.
- `internal/copilot/` — Copilot CLI trust management, binary lookup, and argument building.

### Flow

1. Parse args (consume `--version`/`--help`, forward rest to copilot)
2. Create timestamped scratch dir under `$XDG_DATA_HOME/copilot-scratch/`
3. `git init` in the scratch dir
4. Add scratch dir to copilot's trusted folders config
5. `chdir` to scratch dir
6. `syscall.Exec` copilot with forwarded args (replaces process)

## Conventions

- **TDD**: All code is written test-first. Tests live alongside implementation (`_test.go`).
- **Integration tests**: Use `//go:build integration` tag. Not run in CI by default.
- **XDG compliance**: Respect `$XDG_DATA_HOME` and `$XDG_CONFIG_HOME`.
- **No dependencies**: Pure stdlib. No third-party Go modules.
- **Process replacement**: Use `syscall.Exec` so the extension process doesn't stay resident.

## Release flow

Tag-triggered releases using `cli/gh-extension-precompile`:
1. Push changes to `main`
2. Create a semver tag (`v0.x.y`)
3. Push tag → triggers Release workflow → cross-compiles + creates GitHub Release
4. Users install with `gh extension install ekroon/gh-copilot-scratch`

Use the release skill (`.github/skills/release/`) for automated release orchestration.
