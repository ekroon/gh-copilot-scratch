package copilot

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// BuildArgs returns the arguments to pass to the copilot binary.
func BuildArgs(extra []string) []string {
	if extra == nil {
		return nil
	}
	args := make([]string, len(extra))
	copy(args, extra)
	return args
}

// FindCopilot locates the copilot binary in PATH.
func FindCopilot() (string, error) {
	return exec.LookPath("copilot")
}

// EnsureTrust adds the given directory to copilot's trusted folders config.
func EnsureTrust(dir string, configDir string) error {
	configPath := filepath.Join(configDir, "config.json")

	config := make(map[string]any)

	data, err := os.ReadFile(configPath)
	if err == nil {
		json.Unmarshal(data, &config)
	}

	var folders []string
	if raw, ok := config["trusted_folders"]; ok {
		if arr, ok := raw.([]any); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok {
					folders = append(folders, s)
				}
			}
		}
	}

	// Check for duplicates
	for _, f := range folders {
		if f == dir {
			return nil
		}
	}

	folders = append(folders, dir)
	config["trusted_folders"] = folders

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	out, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	return os.WriteFile(configPath, out, 0o644)
}
