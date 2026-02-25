package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds user configuration from ~/.wsconfig.
type Config struct {
	Editor      string
	CleanupDays int // days of inactivity before a working set is offered for cleanup; 0 = disabled
}

// defaultContent is written to ~/.wsconfig when the file is created for the
// first time, giving users a self-documenting starting point.
const defaultContent = `# ws configuration

# Editor to open files with.
editor=vim

# Number of days of inactivity before a working set is offered for cleanup.
# Set to 0 to disable.
cleanup_days=7
`

// Path returns the path to the user's config file (~/.wsconfig).
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".wsconfig"), nil
}

// EnsureExists creates ~/.wsconfig with default values if it does not already
// exist. It is a no-op when the file is already present.
func EnsureExists() error {
	path, err := Path()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return nil // already exists
	}
	return os.WriteFile(path, []byte(defaultContent), 0o644)
}

// Load reads ~/.wsconfig and returns a Config.
// Missing file is not an error — defaults are used.
func Load() (*Config, error) {
	cfg := &Config{
		Editor:      "vim",
		CleanupDays: 7,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil
	}

	path := filepath.Join(home, ".wsconfig")
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "editor":
			cfg.Editor = val
		case "cleanup_days":
			if n, err := strconv.Atoi(val); err != nil {
				fmt.Fprintf(os.Stderr, "warning: invalid cleanup_days value %q, using default (%d)\n", val, cfg.CleanupDays)
			} else if n < 0 {
				fmt.Fprintf(os.Stderr, "warning: cleanup_days must be >= 0, ignoring value %d\n", n)
			} else {
				cfg.CleanupDays = n
			}
		}
	}

	return cfg, scanner.Err()
}
