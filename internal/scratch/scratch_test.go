package scratch

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestNewDir_CreatesUnderXDGDataHome(t *testing.T) {
	// Set up a temp dir as XDG_DATA_HOME
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	dir, err := NewDir()
	if err != nil {
		t.Fatalf("NewDir() error: %v", err)
	}

	// Should be under $XDG_DATA_HOME/copilot-scratch/
	expectedPrefix := filepath.Join(tmpDir, "copilot-scratch")
	if !strings.HasPrefix(dir, expectedPrefix) {
		t.Errorf("dir = %q, want prefix %q", dir, expectedPrefix)
	}
}

func TestNewDir_DefaultsToHomeDotLocalShare(t *testing.T) {
	// Unset XDG_DATA_HOME to use default
	t.Setenv("XDG_DATA_HOME", "")

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() error: %v", err)
	}

	dir, err := NewDir()
	if err != nil {
		t.Fatalf("NewDir() error: %v", err)
	}
	defer os.RemoveAll(dir)

	expectedPrefix := filepath.Join(home, ".local", "share", "copilot-scratch")
	if !strings.HasPrefix(dir, expectedPrefix) {
		t.Errorf("dir = %q, want prefix %q", dir, expectedPrefix)
	}
}

func TestNewDir_TimestampedNaming(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	dir, err := NewDir()
	if err != nil {
		t.Fatalf("NewDir() error: %v", err)
	}

	// Dir name should match YYYY-MM-DD_HHMMSS-<random>
	name := filepath.Base(dir)
	pattern := `^\d{4}-\d{2}-\d{2}_\d{6}-[a-z0-9]{6}$`
	matched, err := regexp.MatchString(pattern, name)
	if err != nil {
		t.Fatalf("regexp error: %v", err)
	}
	if !matched {
		t.Errorf("dir name = %q, want match for %s", name, pattern)
	}
}

func TestNewDir_DirectoryExists(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	dir, err := NewDir()
	if err != nil {
		t.Fatalf("NewDir() error: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Stat(%q) error: %v", dir, err)
	}
	if !info.IsDir() {
		t.Errorf("%q is not a directory", dir)
	}
}

func TestNewDir_GitInitialized(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	dir, err := NewDir()
	if err != nil {
		t.Fatalf("NewDir() error: %v", err)
	}

	gitDir := filepath.Join(dir, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		t.Fatalf("Stat(%q) error: %v (git init not done?)", gitDir, err)
	}
	if !info.IsDir() {
		t.Errorf("%q is not a directory", gitDir)
	}
}

func TestNewDir_UniqueDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	dir1, err := NewDir()
	if err != nil {
		t.Fatalf("NewDir() #1 error: %v", err)
	}
	dir2, err := NewDir()
	if err != nil {
		t.Fatalf("NewDir() #2 error: %v", err)
	}

	if dir1 == dir2 {
		t.Errorf("two calls returned same dir: %q", dir1)
	}
}

func TestDataDir_RespectsXDG(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	got := DataDir()
	want := filepath.Join(tmpDir, "copilot-scratch")
	if got != want {
		t.Errorf("DataDir() = %q, want %q", got, want)
	}
}

func TestDataDir_FallsBackToHome(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "")

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() error: %v", err)
	}

	got := DataDir()
	want := filepath.Join(home, ".local", "share", "copilot-scratch")
	if got != want {
		t.Errorf("DataDir() = %q, want %q", got, want)
	}
}
