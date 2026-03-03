package copilot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildArgs_NoExtraArgs(t *testing.T) {
	args := BuildArgs(nil)
	if len(args) != 0 {
		t.Errorf("BuildArgs(nil) = %v, want empty", args)
	}
}

func TestBuildArgs_ForwardsArgs(t *testing.T) {
	input := []string{"--model", "claude-sonnet-4.5", "-v"}
	args := BuildArgs(input)

	if len(args) != len(input) {
		t.Fatalf("BuildArgs() len = %d, want %d", len(args), len(input))
	}
	for i, a := range args {
		if a != input[i] {
			t.Errorf("args[%d] = %q, want %q", i, a, input[i])
		}
	}
}

func TestEnsureTrust_CreatesConfigIfNotExists(t *testing.T) {
	tmpHome := t.TempDir()
	configDir := filepath.Join(tmpHome, ".config", "github-copilot")

	scratchDir := filepath.Join(t.TempDir(), "scratch-session")
	os.MkdirAll(scratchDir, 0o755)

	err := EnsureTrust(scratchDir, configDir)
	if err != nil {
		t.Fatalf("EnsureTrust() error: %v", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error: %v", configPath, err)
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	trusted, ok := config["trusted_folders"]
	if !ok {
		t.Fatal("config missing trusted_folders key")
	}

	folders, ok := trusted.([]any)
	if !ok {
		t.Fatalf("trusted_folders type = %T, want []any", trusted)
	}

	found := false
	for _, f := range folders {
		if f == scratchDir {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("trusted_folders = %v, want to contain %q", folders, scratchDir)
	}
}

func TestEnsureTrust_AppendsToExistingConfig(t *testing.T) {
	configDir := t.TempDir()
	existingDir := "/some/existing/dir"

	// Write existing config with one trusted folder
	configPath := filepath.Join(configDir, "config.json")
	existing := map[string]any{
		"trusted_folders": []string{existingDir},
		"otherKey":        "preserved",
	}
	data, _ := json.Marshal(existing)
	os.WriteFile(configPath, data, 0o644)

	scratchDir := filepath.Join(t.TempDir(), "new-scratch")
	os.MkdirAll(scratchDir, 0o755)

	err := EnsureTrust(scratchDir, configDir)
	if err != nil {
		t.Fatalf("EnsureTrust() error: %v", err)
	}

	data, _ = os.ReadFile(configPath)
	var config map[string]any
	json.Unmarshal(data, &config)

	// Check both folders are present
	folders := config["trusted_folders"].([]any)
	if len(folders) != 2 {
		t.Fatalf("trustedFolders len = %d, want 2", len(folders))
	}

	// Check otherKey is preserved
	if config["otherKey"] != "preserved" {
		t.Errorf("otherKey = %v, want 'preserved'", config["otherKey"])
	}
}

func TestEnsureTrust_NoDuplicates(t *testing.T) {
	configDir := t.TempDir()
	scratchDir := filepath.Join(t.TempDir(), "scratch")
	os.MkdirAll(scratchDir, 0o755)

	// Trust twice
	EnsureTrust(scratchDir, configDir)
	EnsureTrust(scratchDir, configDir)

	configPath := filepath.Join(configDir, "config.json")
	data, _ := os.ReadFile(configPath)
	var config map[string]any
	json.Unmarshal(data, &config)

	folders := config["trusted_folders"].([]any)
	if len(folders) != 1 {
		t.Errorf("trusted_folders len = %d after double trust, want 1", len(folders))
	}
}

func TestFindCopilot_ReturnsPath(t *testing.T) {
	path, err := FindCopilot()
	if err != nil {
		t.Skipf("copilot not in PATH, skipping: %v", err)
	}
	if path == "" {
		t.Error("FindCopilot() returned empty path")
	}
}
