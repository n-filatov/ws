package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// RootDir returns the root directory of the current git repository.
func RootDir() (string, error) {
	out, err := run("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not inside a git repository")
	}
	return strings.TrimSpace(out), nil
}

// CurrentBranch returns the name of the current git branch.
// Works on repos with no commits yet by falling back to symbolic-ref.
func CurrentBranch() (string, error) {
	out, err := run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err == nil {
		if b := strings.TrimSpace(out); b != "" && b != "HEAD" {
			return b, nil
		}
	}
	// Fallback for empty repos (no commits yet)
	out, err = run("git", "symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(out), nil
}

// ModifiedFiles returns a map of absolute file path → git status code.
// Status codes: "M" (modified), "A" (added/staged), "?" (untracked).
func ModifiedFiles(root string) (map[string]string, error) {
	out, err := run("git", "-C", root, "status", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("failed to run git status: %w", err)
	}

	result := make(map[string]string)
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 4 {
			continue
		}
		xy := line[:2]
		path := strings.TrimSpace(line[3:])

		// Handle renamed files (old -> new)
		if strings.Contains(path, " -> ") {
			parts := strings.SplitN(path, " -> ", 2)
			path = parts[1]
		}

		absPath := filepath.Join(root, path)

		var status string
		switch {
		case strings.Contains(xy, "?"):
			status = "?"
		case xy[0] == 'A' || xy[1] == 'A':
			status = "A"
		default:
			status = "M"
		}

		result[absPath] = status
	}
	return result, nil
}

// Checkout reverts a file to its last committed state via git checkout.
func Checkout(file string) error {
	_, err := run("git", "checkout", "--", file)
	if err != nil {
		return fmt.Errorf("git checkout failed: %w", err)
	}
	return nil
}

func run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	return string(out), err
}
