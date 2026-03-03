package scratch

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// DataDir returns the base directory for scratch sessions.
// Uses $XDG_DATA_HOME/copilot-scratch/ or ~/.local/share/copilot-scratch/.
func DataDir() string {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "copilot-scratch")
}

// NewDir creates a new timestamped scratch directory with git init.
func NewDir() (string, error) {
	base := DataDir()

	name := generateName()
	dir := filepath.Join(base, name)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("creating scratch dir: %w", err)
	}

	if err := gitInit(dir); err != nil {
		return "", fmt.Errorf("git init: %w", err)
	}

	return dir, nil
}

func generateName() string {
	ts := time.Now().Format("2006-01-02_150405")
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("%s-%s", ts, hex.EncodeToString(b))
}

func gitInit(dir string) error {
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = dir
	return cmd.Run()
}
