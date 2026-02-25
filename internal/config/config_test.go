package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── EnsureExists ─────────────────────────────────────────────────────────────

func TestEnsureExists_CreatesFileWithDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	if err := EnsureExists(); err != nil {
		t.Fatalf("EnsureExists: %v", err)
	}

	path := filepath.Join(home, ".wsconfig")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("config file not created: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "editor=vim") {
		t.Errorf("expected default editor in file:\n%s", content)
	}
	if !strings.Contains(content, "cleanup_days=7") {
		t.Errorf("expected default cleanup_days in file:\n%s", content)
	}
}

func TestEnsureExists_NoopWhenFileExists(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	path := filepath.Join(home, ".wsconfig")
	custom := "editor=nano\n"
	if err := os.WriteFile(path, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := EnsureExists(); err != nil {
		t.Fatalf("EnsureExists: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != custom {
		t.Errorf("EnsureExists overwrote existing file: got %q", string(data))
	}
}

func TestEnsureExists_DefaultsAreLoadable(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	if err := EnsureExists(); err != nil {
		t.Fatalf("EnsureExists: %v", err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load after EnsureExists: %v", err)
	}
	if cfg.Editor != "vim" {
		t.Errorf("editor: want vim, got %q", cfg.Editor)
	}
	if cfg.CleanupDays != 7 {
		t.Errorf("cleanup_days: want 7, got %d", cfg.CleanupDays)
	}
}

// writeConfig writes content to ~/.wsconfig inside a temp home directory and
// sets HOME so config.Load() picks it up. Returns a cleanup function.
func writeConfig(t *testing.T, content string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	// Also clear USERPROFILE on platforms that use it.
	t.Setenv("USERPROFILE", home)

	path := filepath.Join(home, ".wsconfig")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ── Defaults ─────────────────────────────────────────────────────────────────

func TestLoad_MissingFile_UsesDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Editor != "vim" {
		t.Errorf("default editor: want %q, got %q", "vim", cfg.Editor)
	}
	if cfg.CleanupDays != 7 {
		t.Errorf("default cleanup_days: want 7, got %d", cfg.CleanupDays)
	}
}

// ── Valid values ──────────────────────────────────────────────────────────────

func TestLoad_ValidEditor(t *testing.T) {
	writeConfig(t, "editor=nano\n")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Editor != "nano" {
		t.Errorf("want %q, got %q", "nano", cfg.Editor)
	}
}

func TestLoad_ValidCleanupDays(t *testing.T) {
	writeConfig(t, "cleanup_days=14\n")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CleanupDays != 14 {
		t.Errorf("want 14, got %d", cfg.CleanupDays)
	}
}

func TestLoad_ZeroCleanupDays_Disabled(t *testing.T) {
	writeConfig(t, "cleanup_days=0\n")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CleanupDays != 0 {
		t.Errorf("want 0 (disabled), got %d", cfg.CleanupDays)
	}
}

func TestLoad_Comments_Ignored(t *testing.T) {
	writeConfig(t, "# this is a comment\neditor=emacs\n# another comment\n")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Editor != "emacs" {
		t.Errorf("want %q, got %q", "emacs", cfg.Editor)
	}
}

func TestLoad_UnknownKeys_Ignored(t *testing.T) {
	writeConfig(t, "unknown_key=value\neditor=nano\n")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Editor != "nano" {
		t.Errorf("want %q, got %q", "nano", cfg.Editor)
	}
}

func TestLoad_MalformedLine_Ignored(t *testing.T) {
	writeConfig(t, "this-has-no-equals-sign\neditor=nano\n")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Editor != "nano" {
		t.Errorf("want %q, got %q", "nano", cfg.Editor)
	}
}

// ── Invalid / warning-worthy values ──────────────────────────────────────────

func TestLoad_InvalidCleanupDays_KeepsDefault(t *testing.T) {
	// Capture stderr to verify warning is emitted.
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	writeConfig(t, "cleanup_days=abc\n")
	cfg, err := Load()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	warning := buf.String()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CleanupDays != 7 {
		t.Errorf("expected default 7, got %d", cfg.CleanupDays)
	}
	if warning == "" {
		t.Error("expected a warning on stderr for invalid cleanup_days, got none")
	}
}

func TestLoad_NegativeCleanupDays_KeepsDefault(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	writeConfig(t, "cleanup_days=-3\n")
	cfg, err := Load()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	warning := buf.String()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CleanupDays != 7 {
		t.Errorf("expected default 7, got %d", cfg.CleanupDays)
	}
	if warning == "" {
		t.Error("expected a warning on stderr for negative cleanup_days, got none")
	}
}

func TestLoad_MultipleKeys(t *testing.T) {
	writeConfig(t, fmt.Sprintf("editor=code\ncleanup_days=30\n"))
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Editor != "code" {
		t.Errorf("editor: want %q, got %q", "code", cfg.Editor)
	}
	if cfg.CleanupDays != 30 {
		t.Errorf("cleanup_days: want 30, got %d", cfg.CleanupDays)
	}
}
