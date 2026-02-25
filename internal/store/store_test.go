package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// tempWsPath returns a temporary file path inside t.TempDir() for use as a
// working set path. The file does not need to exist beforehand.
func tempWsPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), ".workingset-test")
}

// ── Load ─────────────────────────────────────────────────────────────────────

func TestLoad_MissingFile(t *testing.T) {
	wsPath := tempWsPath(t)
	files, err := Load(wsPath)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected empty slice, got %v", files)
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	wsPath := tempWsPath(t)
	if err := os.WriteFile(wsPath, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	files, err := Load(wsPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected empty slice, got %v", files)
	}
}

func TestLoad_PopulatedFile(t *testing.T) {
	wsPath := tempWsPath(t)
	content := "/repo/a.go\n/repo/b.go\n"
	if err := os.WriteFile(wsPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	files, err := Load(wsPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"/repo/a.go", "/repo/b.go"}
	if len(files) != len(want) {
		t.Fatalf("expected %v, got %v", want, files)
	}
	for i := range want {
		if files[i] != want[i] {
			t.Errorf("files[%d]: want %q, got %q", i, want[i], files[i])
		}
	}
}

func TestLoad_IgnoresBlankLines(t *testing.T) {
	wsPath := tempWsPath(t)
	content := "\n/repo/a.go\n\n/repo/b.go\n\n"
	if err := os.WriteFile(wsPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	files, err := Load(wsPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %v", files)
	}
}

// ── Save / Load roundtrip ─────────────────────────────────────────────────────

func TestSave_CreatesParentDir(t *testing.T) {
	wsPath := filepath.Join(t.TempDir(), "nested", "dir", ".workingset-main")
	if err := Save(wsPath, []string{"/a.go"}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if _, err := os.Stat(wsPath); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestSaveLoad_Roundtrip(t *testing.T) {
	wsPath := tempWsPath(t)
	want := []string{"/repo/main.go", "/repo/internal/config.go"}
	if err := Save(wsPath, want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(wsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("want %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] want %q, got %q", i, want[i], got[i])
		}
	}
}

// ── Add ───────────────────────────────────────────────────────────────────────

func TestAdd_DeduplicatesFiles(t *testing.T) {
	wsPath := tempWsPath(t)
	dir := t.TempDir()
	f := filepath.Join(dir, "a.go")

	if err := Add(wsPath, f, f, f); err != nil {
		t.Fatalf("Add: %v", err)
	}

	files, err := Load(wsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file after dedup, got %d: %v", len(files), files)
	}
}

func TestAdd_DeduplicatesAcrossCalls(t *testing.T) {
	wsPath := tempWsPath(t)
	dir := t.TempDir()
	f := filepath.Join(dir, "a.go")

	if err := Add(wsPath, f); err != nil {
		t.Fatalf("first Add: %v", err)
	}
	if err := Add(wsPath, f); err != nil {
		t.Fatalf("second Add: %v", err)
	}

	files, err := Load(wsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d: %v", len(files), files)
	}
}

func TestAdd_PreservesOrder(t *testing.T) {
	wsPath := tempWsPath(t)
	dir := t.TempDir()
	a := filepath.Join(dir, "a.go")
	b := filepath.Join(dir, "b.go")
	c := filepath.Join(dir, "c.go")

	if err := Add(wsPath, a, b, c); err != nil {
		t.Fatalf("Add: %v", err)
	}
	files, err := Load(wsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := []string{a, b, c}
	for i := range want {
		if files[i] != want[i] {
			t.Errorf("[%d] want %q, got %q", i, want[i], files[i])
		}
	}
}

func TestAdd_ResolvesRelativePaths(t *testing.T) {
	wsPath := tempWsPath(t)

	// Add a relative path that resolves to an absolute one.
	if err := Add(wsPath, "."); err != nil {
		t.Fatalf("Add: %v", err)
	}
	files, err := Load(wsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %v", files)
	}
	if !filepath.IsAbs(files[0]) {
		t.Errorf("expected absolute path, got %q", files[0])
	}
}

// ── Remove ────────────────────────────────────────────────────────────────────

func TestRemove_ExistingFile(t *testing.T) {
	wsPath := tempWsPath(t)
	dir := t.TempDir()
	a := filepath.Join(dir, "a.go")
	b := filepath.Join(dir, "b.go")

	if err := Save(wsPath, []string{a, b}); err != nil {
		t.Fatal(err)
	}
	if err := Remove(wsPath, a); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	files, err := Load(wsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(files) != 1 || files[0] != b {
		t.Fatalf("expected [%s], got %v", b, files)
	}
}

func TestRemove_NonExistentFile_NoError(t *testing.T) {
	wsPath := tempWsPath(t)
	dir := t.TempDir()
	a := filepath.Join(dir, "a.go")
	ghost := filepath.Join(dir, "ghost.go")

	if err := Save(wsPath, []string{a}); err != nil {
		t.Fatal(err)
	}
	if err := Remove(wsPath, ghost); err != nil {
		t.Fatalf("Remove of non-existent path should not error, got: %v", err)
	}
	files, err := Load(wsPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected list unchanged, got %v", files)
	}
}

// ── StaleCandidates ───────────────────────────────────────────────────────────

func TestStaleCandidates_ReturnsOldBranches(t *testing.T) {
	// Build a fake ws data directory matching the layout expected by StaleCandidates.
	dir := t.TempDir()
	root := filepath.Join(dir, "myrepo")

	// Override XDG_DATA_HOME so dataHome() returns our temp dir.
	t.Setenv("XDG_DATA_HOME", dir)

	slug := repoSlug(root)
	wsDir := filepath.Join(dir, "ws", slug)
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	oldFile := filepath.Join(wsDir, ".workingset-old-branch")
	if err := os.WriteFile(oldFile, []byte("/repo/a.go\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Back-date the file to simulate staleness.
	old := time.Now().Add(-30 * 24 * time.Hour)
	if err := os.Chtimes(oldFile, old, old); err != nil {
		t.Fatal(err)
	}

	candidates := StaleCandidates(root, "main", 7)
	if len(candidates) != 1 {
		t.Fatalf("expected 1 stale candidate, got %d", len(candidates))
	}
	if candidates[0].Branch != "old-branch" {
		t.Errorf("expected branch %q, got %q", "old-branch", candidates[0].Branch)
	}
}

func TestStaleCandidates_ExcludesCurrentBranch(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "myrepo")
	t.Setenv("XDG_DATA_HOME", dir)

	slug := repoSlug(root)
	wsDir := filepath.Join(dir, "ws", slug)
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a stale file for the "current" branch — should NOT appear in results.
	currentFile := filepath.Join(wsDir, ".workingset-main")
	if err := os.WriteFile(currentFile, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-30 * 24 * time.Hour)
	if err := os.Chtimes(currentFile, old, old); err != nil {
		t.Fatal(err)
	}

	candidates := StaleCandidates(root, "main", 7)
	if len(candidates) != 0 {
		t.Fatalf("current branch should be excluded, got %v", candidates)
	}
}

func TestStaleCandidates_ExcludesRecentBranches(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "myrepo")
	t.Setenv("XDG_DATA_HOME", dir)

	slug := repoSlug(root)
	wsDir := filepath.Join(dir, "ws", slug)
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a recently-touched branch file — should NOT be stale.
	recentFile := filepath.Join(wsDir, ".workingset-feature-x")
	if err := os.WriteFile(recentFile, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	// Modification time is "now" by default, so within the 7-day window.

	candidates := StaleCandidates(root, "main", 7)
	if len(candidates) != 0 {
		t.Fatalf("recent branch should not be stale, got %v", candidates)
	}
}
