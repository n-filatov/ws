package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Config holds user configuration from ~/.wsconfig.
type Config struct {
	Editor string
}

// Load reads ~/.wsconfig and returns a Config.
// Missing file is not an error — defaults are used.
func Load() (*Config, error) {
	cfg := &Config{
		Editor: "vim",
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
		}
	}

	return cfg, scanner.Err()
}
