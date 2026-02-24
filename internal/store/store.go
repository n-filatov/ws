package store

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// dataHome returns the XDG data home directory.
func dataHome() string {
	if x := os.Getenv("XDG_DATA_HOME"); x != "" {
		return x
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share")
}

// repoSlug returns a stable, human-readable directory name for a repo root.
func repoSlug(root string) string {
	sum := sha1.Sum([]byte(root))
	return filepath.Base(root) + "-" + fmt.Sprintf("%x", sum)[:8]
}

// WorkingSetPath returns the path to the .workingset-<branch> file
// stored under the user's data directory (never inside the project repo).
func WorkingSetPath(root, branch string) string {
	safe := strings.ReplaceAll(branch, "/", "-")
	return filepath.Join(dataHome(), "ws", repoSlug(root), ".workingset-"+safe)
}

// WriteRepoPath records the repo root in the ws data dir so future gc can
// detect orphaned entries. Safe to call on every startup.
func WriteRepoPath(root string) {
	dir := filepath.Join(dataHome(), "ws", repoSlug(root))
	_ = os.MkdirAll(dir, 0o755)
	_ = os.WriteFile(filepath.Join(dir, ".repo-path"), []byte(root+"\n"), 0o644)
}

// MigrateIfNeeded moves a legacy .workingset-<branch> file from the repo root
// to the new data directory location. Silent no-op if nothing to migrate.
func MigrateIfNeeded(root, branch string) {
	newPath := WorkingSetPath(root, branch)
	if _, err := os.Stat(newPath); err == nil {
		return // already at new location
	}
	safe := strings.ReplaceAll(branch, "/", "-")
	oldPath := filepath.Join(root, ".workingset-"+safe)
	if _, err := os.Stat(oldPath); err != nil {
		return // nothing to migrate
	}
	_ = os.MkdirAll(filepath.Dir(newPath), 0o755)
	_ = os.Rename(oldPath, newPath)
}

// StaleCandidate describes a working set file that hasn't been used recently.
type StaleCandidate struct {
	Branch   string    // sanitized branch name (slashes replaced with dashes)
	WsPath   string    // full path to the .workingset-* file
	LastUsed time.Time // file modification time
}

// StaleCandidates returns working set files for the given repo that haven't
// been modified in more than maxAgeDays, excluding the current branch.
func StaleCandidates(root, currentBranch string, maxAgeDays int) []StaleCandidate {
	dir := filepath.Join(dataHome(), "ws", repoSlug(root))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	currentFile := ".workingset-" + strings.ReplaceAll(currentBranch, "/", "-")
	cutoff := time.Now().AddDate(0, 0, -maxAgeDays)

	var result []StaleCandidate
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, ".workingset-") || name == currentFile {
			continue
		}
		info, err := e.Info()
		if err != nil || info.ModTime().After(cutoff) {
			continue
		}
		result = append(result, StaleCandidate{
			Branch:   strings.TrimPrefix(name, ".workingset-"),
			WsPath:   filepath.Join(dir, name),
			LastUsed: info.ModTime(),
		})
	}
	return result
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
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create working set directory: %w", err)
	}
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
