package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/ekroon/gh-copilot-scratch/internal/copilot"
	"github.com/ekroon/gh-copilot-scratch/internal/scratch"
)

var version = "dev"

type options struct {
	version    bool
	help       bool
	copilotArgs []string
}

func parseArgs(args []string) options {
	var opts options
	for _, a := range args {
		switch a {
		case "--version":
			opts.version = true
		case "--help", "-h":
			opts.help = true
		default:
			opts.copilotArgs = append(opts.copilotArgs, a)
		}
	}
	return opts
}

func usageText() string {
	return `gh copilot-scratch — launch Copilot CLI in a scratch directory

Usage:
  gh copilot-scratch [flags] [-- copilot-args...]

Flags:
  --help, -h      Show this help
  --version       Show version

All other arguments are forwarded to the copilot CLI.

Examples:
  gh copilot-scratch
  gh copilot-scratch --model claude-sonnet-4.5
`
}

func run(args []string) error {
	opts := parseArgs(args)

	if opts.help {
		fmt.Print(usageText())
		return nil
	}

	if opts.version {
		fmt.Printf("gh-copilot-scratch %s\n", version)
		return nil
	}

	// Find copilot binary
	copilotPath, err := copilot.FindCopilot()
	if err != nil {
		return fmt.Errorf("copilot CLI not found in PATH: %w", err)
	}

	// Create scratch directory
	dir, err := scratch.NewDir()
	if err != nil {
		return fmt.Errorf("creating scratch directory: %w", err)
	}

	fmt.Printf("Scratch directory: %s\n", dir)

	// Trust the directory in copilot config
	configDir := copilotConfigDir()
	if err := copilot.EnsureTrust(dir, configDir); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not auto-trust directory: %v\n", err)
	}

	// Build copilot command args
	execArgs := []string{"copilot"}
	execArgs = append(execArgs, copilot.BuildArgs(opts.copilotArgs)...)

	// Set working directory and exec copilot
	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("changing to scratch dir: %w", err)
	}

	return syscall.Exec(copilotPath, execArgs, os.Environ())
}

func copilotConfigDir() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "github-copilot")
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
