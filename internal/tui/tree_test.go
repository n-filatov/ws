package tui

import (
	"strings"
	"testing"
)

// helpers

func makeFiles(relPaths ...string) []FileEntry {
	files := make([]FileEntry, len(relPaths))
	for i, p := range relPaths {
		files[i] = FileEntry{AbsPath: "/" + p, RelPath: p, Exists: true}
	}
	return files
}

// ── buildPathTree ─────────────────────────────────────────────────────────────

func TestBuildPathTree_Empty(t *testing.T) {
	root := buildPathTree(nil)
	if len(root.children) != 0 {
		t.Errorf("expected no children for empty input, got %d", len(root.children))
	}
}

func TestBuildPathTree_FlatFiles(t *testing.T) {
	files := makeFiles("a.go", "b.go", "c.go")
	root := buildPathTree(files)
	if len(root.children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(root.children))
	}
	for _, child := range root.children {
		if child.isDir {
			t.Errorf("expected file node for %q, got dir", child.name)
		}
	}
}

func TestBuildPathTree_NestedFiles(t *testing.T) {
	files := makeFiles("cmd/main.go", "internal/config/config.go", "internal/git/git.go")
	root := buildPathTree(files)

	// Should have two dir children: cmd/ and internal/
	if len(root.children) != 2 {
		t.Fatalf("expected 2 top-level children (cmd, internal), got %d", len(root.children))
	}

	// internal/ should have 2 child dirs
	var internal *pathNode
	for _, c := range root.children {
		if c.name == "internal" {
			internal = c
		}
	}
	if internal == nil {
		t.Fatal("expected 'internal' child")
	}
	if len(internal.children) != 2 {
		t.Errorf("expected 2 children under internal/, got %d", len(internal.children))
	}
}

func TestBuildPathTree_SortedDirsFirst(t *testing.T) {
	files := makeFiles("z.go", "a/b.go")
	root := buildPathTree(files)
	if len(root.children) != 2 {
		t.Fatalf("expected 2 children")
	}
	if !root.children[0].isDir {
		t.Errorf("expected directory first, got file %q", root.children[0].name)
	}
}

// ── tryCollapse ───────────────────────────────────────────────────────────────

func TestTryCollapse_SingleFile(t *testing.T) {
	// dir/ → file.go  →  collapses to "file.go"
	dir := &pathNode{name: "dir", isDir: true, fileIdx: -1, children: []*pathNode{
		{name: "file.go", isDir: false, fileIdx: 2},
	}}
	name, idx, ok := tryCollapse(dir)
	if !ok {
		t.Fatal("expected collapse to succeed")
	}
	if name != "file.go" {
		t.Errorf("want %q, got %q", "file.go", name)
	}
	if idx != 2 {
		t.Errorf("want fileIdx 2, got %d", idx)
	}
}

func TestTryCollapse_ChainedDirs(t *testing.T) {
	// a/ → b/ → file.go  →  collapses to "b/file.go"
	inner := &pathNode{name: "b", isDir: true, fileIdx: -1, children: []*pathNode{
		{name: "file.go", isDir: false, fileIdx: 1},
	}}
	outer := &pathNode{name: "a", isDir: true, fileIdx: -1, children: []*pathNode{inner}}

	name, _, ok := tryCollapse(outer)
	if !ok {
		t.Fatal("expected collapse to succeed")
	}
	if name != "b/file.go" {
		t.Errorf("want %q, got %q", "b/file.go", name)
	}
}

func TestTryCollapse_MultipleChildren_NoCollapse(t *testing.T) {
	dir := &pathNode{name: "dir", isDir: true, fileIdx: -1, children: []*pathNode{
		{name: "a.go", isDir: false, fileIdx: 0},
		{name: "b.go", isDir: false, fileIdx: 1},
	}}
	_, _, ok := tryCollapse(dir)
	if ok {
		t.Error("expected no collapse for node with multiple children")
	}
}

// ── BuildFileTree ─────────────────────────────────────────────────────────────

func TestBuildFileTree_Empty(t *testing.T) {
	rendered, lineToFileIdx, navItems := BuildFileTree(nil, "myrepo", nil)
	if rendered == "" {
		t.Error("expected non-empty rendered string (at least root header)")
	}
	// Line 0 is the root header
	if len(lineToFileIdx) == 0 || lineToFileIdx[0] != -1 {
		t.Errorf("line 0 should map to -1 (root header), got %v", lineToFileIdx)
	}
	if len(navItems) != 0 {
		t.Errorf("expected no nav items for empty tree, got %d", len(navItems))
	}
}

func TestBuildFileTree_ContainsFilenames(t *testing.T) {
	files := makeFiles("main.go", "config.go")
	rendered, _, _ := BuildFileTree(files, "myrepo", nil)
	if !strings.Contains(rendered, "main.go") {
		t.Errorf("rendered tree missing 'main.go':\n%s", rendered)
	}
	if !strings.Contains(rendered, "config.go") {
		t.Errorf("rendered tree missing 'config.go':\n%s", rendered)
	}
}

func TestBuildFileTree_NavItemsMatchFiles(t *testing.T) {
	files := makeFiles("a.go", "b.go", "c.go")
	_, lineToFileIdx, navItems := BuildFileTree(files, "myrepo", nil)

	if len(navItems) != 3 {
		t.Fatalf("expected 3 nav items for 3 flat files, got %d", len(navItems))
	}
	for i, item := range navItems {
		if item.fileIdx < 0 {
			t.Errorf("navItems[%d]: expected file idx >= 0, got %d", i, item.fileIdx)
		}
		if item.lineIdx >= len(lineToFileIdx) {
			t.Errorf("navItems[%d].lineIdx %d out of range (len=%d)", i, item.lineIdx, len(lineToFileIdx))
		}
		if lineToFileIdx[item.lineIdx] != item.fileIdx {
			t.Errorf("lineToFileIdx mismatch at line %d: want %d, got %d",
				item.lineIdx, item.fileIdx, lineToFileIdx[item.lineIdx])
		}
	}
}

func TestBuildFileTree_CollapsedDirShowsAsLeaf(t *testing.T) {
	files := makeFiles("internal/a.go", "internal/b.go")
	collapsed := map[string]bool{"internal": true}
	_, _, navItems := BuildFileTree(files, "myrepo", collapsed)

	// With internal/ collapsed we expect exactly 1 nav item (the dir header).
	if len(navItems) != 1 {
		t.Fatalf("expected 1 nav item for collapsed dir, got %d", len(navItems))
	}
	if navItems[0].fileIdx != -1 {
		t.Errorf("collapsed dir nav item should have fileIdx=-1, got %d", navItems[0].fileIdx)
	}
	if navItems[0].dirPath != "internal" {
		t.Errorf("expected dirPath %q, got %q", "internal", navItems[0].dirPath)
	}
}

func TestBuildFileTree_RepoNameInOutput(t *testing.T) {
	rendered, _, _ := BuildFileTree(nil, "myawesomerepo", nil)
	if !strings.Contains(rendered, "myawesomerepo") {
		t.Errorf("expected repo name in rendered output:\n%s", rendered)
	}
}

// ── basenameMatchPositions ────────────────────────────────────────────────────

func TestBasenameMatchPositions_ReturnsRelativePositions(t *testing.T) {
	// relPath = "internal/config.go", basename = "config.go" (starts at idx 9)
	// matched indexes [9, 10] → relative positions [0, 1]
	positions := basenameMatchPositions([]int{9, 10}, "internal/config.go")
	if len(positions) != 2 || positions[0] != 0 || positions[1] != 1 {
		t.Errorf("expected [0 1], got %v", positions)
	}
}

func TestBasenameMatchPositions_FiltersOutDirPositions(t *testing.T) {
	// Matches inside the directory portion should be excluded.
	positions := basenameMatchPositions([]int{0, 1, 9}, "internal/config.go")
	// Only idx 9 (start of "config.go") survives → relative position 0
	if len(positions) != 1 || positions[0] != 0 {
		t.Errorf("expected [0], got %v", positions)
	}
}

func TestBasenameMatchPositions_FlatFile(t *testing.T) {
	// No directory separator — all positions pass through unchanged.
	positions := basenameMatchPositions([]int{0, 2, 4}, "main.go")
	if len(positions) != 3 {
		t.Errorf("expected 3 positions, got %v", positions)
	}
}
