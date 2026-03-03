# gh-copilot-scratch

A [gh CLI](https://cli.github.com/) extension that launches [GitHub Copilot CLI](https://docs.github.com/copilot/how-tos/copilot-cli) sessions in auto-created scratch directories.

## What it does

When you run `gh copilot-scratch`, it:

1. Creates a fresh, timestamped scratch directory under `~/.local/share/copilot-scratch/`
2. Initializes a git repository in it
3. Ensures the directory is trusted by Copilot CLI
4. Launches Copilot CLI in that directory
5. Prints the scratch directory path when the session ends

This is useful for quick experiments, throwaway coding sessions, or when you want a clean workspace without setting up a project first.

## Installation

```sh
gh extension install ekroon/gh-copilot-scratch
```

## Usage

```sh
# Start a scratch session
gh copilot-scratch

# Pass extra arguments to copilot
gh copilot-scratch --model claude-sonnet-4.5

# Show version
gh copilot-scratch --version
```

## Prerequisites

- [gh CLI](https://cli.github.com/) installed and authenticated
- [Copilot CLI](https://docs.github.com/copilot/how-tos/copilot-cli) installed
- `git` in PATH

## Scratch directory location

Directories are created under `$XDG_DATA_HOME/copilot-scratch/` (defaults to `~/.local/share/copilot-scratch/`).

Each session gets a directory named `YYYY-MM-DD_HHMMSS-<random>/` for easy sorting and uniqueness.

Scratch directories are **not** automatically cleaned up — the path is printed after each session so you can find your work.

## Development

```sh
# Build
go build -o gh-copilot-scratch ./cmd/gh-copilot-scratch

# Test
go test -race ./...

# Lint
go vet ./...

# Integration tests (requires copilot CLI)
go test -race -tags=integration ./...
```

## License

MIT
