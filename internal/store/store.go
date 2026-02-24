package store

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WorkingSetPath returns the path to the .workingset-<branch> file.
func WorkingSetPath(root, branch string) string {
	// Sanitize branch name for use as filename (replace / with -)
	safe := strings.ReplaceAll(branch, "/", "-")
	return filepath.Join(root, ".workingset-"+safe)
}

// Load reads the working set file and returns a slice of absolute paths.
// Returns an empty slice if the file does not exist.
func Load(path string) ([]string, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return []string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open working set: %w", err)
	}
	defer f.Close()

	var files []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			files = append(files, line)
		}
	}
	return files, scanner.Err()
}

// Save writes the list of absolute paths to the working set file, one per line.
func Save(path string, files []string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to write working set: %w", err)
	}
	defer f.Close()

	for _, file := range files {
		if _, err := fmt.Fprintln(f, file); err != nil {
			return err
		}
	}
	return nil
}

// Add resolves the given file paths to absolute paths, deduplicates, and saves.
func Add(wsPath string, files ...string) error {
	existing, err := Load(wsPath)
	if err != nil {
		return err
	}

	seen := make(map[string]struct{}, len(existing))
	for _, f := range existing {
		seen[f] = struct{}{}
	}

	for _, f := range files {
		abs, err := filepath.Abs(f)
		if err != nil {
			return fmt.Errorf("failed to resolve path %q: %w", f, err)
		}
		if _, ok := seen[abs]; !ok {
			existing = append(existing, abs)
			seen[abs] = struct{}{}
		}
	}

	return Save(wsPath, existing)
}

// Remove removes the given absolute path from the working set file.
func Remove(wsPath string, file string) error {
	abs, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("failed to resolve path %q: %w", file, err)
	}

	existing, err := Load(wsPath)
	if err != nil {
		return err
	}

	filtered := existing[:0]
	for _, f := range existing {
		if f != abs {
			filtered = append(filtered, f)
		}
	}

	return Save(wsPath, filtered)
}
