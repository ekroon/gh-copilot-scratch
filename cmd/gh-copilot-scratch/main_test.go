package main

import (
	"strings"
	"testing"
)

func TestParseArgs_Version(t *testing.T) {
	opts := parseArgs([]string{"--version"})
	if !opts.version {
		t.Error("--version not parsed")
	}
}

func TestParseArgs_Help(t *testing.T) {
	opts := parseArgs([]string{"--help"})
	if !opts.help {
		t.Error("--help not parsed")
	}
}

func TestParseArgs_HelpShort(t *testing.T) {
	opts := parseArgs([]string{"-h"})
	if !opts.help {
		t.Error("-h not parsed")
	}
}

func TestParseArgs_ExtraArgs(t *testing.T) {
	opts := parseArgs([]string{"--model", "claude-sonnet-4.5", "-v"})
	if len(opts.copilotArgs) != 3 {
		t.Fatalf("copilotArgs len = %d, want 3", len(opts.copilotArgs))
	}
	if opts.copilotArgs[0] != "--model" {
		t.Errorf("copilotArgs[0] = %q, want --model", opts.copilotArgs[0])
	}
}

func TestParseArgs_MixedFlags(t *testing.T) {
	opts := parseArgs([]string{"--version", "--model", "gpt-4"})
	if !opts.version {
		t.Error("--version not parsed")
	}
	// version flag is consumed, rest forwarded
	if len(opts.copilotArgs) != 2 {
		t.Errorf("copilotArgs = %v, want [--model gpt-4]", opts.copilotArgs)
	}
}

func TestParseArgs_NoArgs(t *testing.T) {
	opts := parseArgs(nil)
	if opts.version || opts.help {
		t.Error("flags should be false for no args")
	}
	if len(opts.copilotArgs) != 0 {
		t.Errorf("copilotArgs = %v, want empty", opts.copilotArgs)
	}
}

func TestUsageText(t *testing.T) {
	text := usageText()
	if !strings.Contains(text, "gh copilot-scratch") {
		t.Error("usage text should mention 'gh copilot-scratch'")
	}
	if !strings.Contains(text, "--version") {
		t.Error("usage text should mention --version")
	}
}
