package tui

import (
	"path/filepath"
	"sort"
	"strings"

	lipglosstree "github.com/charmbracelet/lipgloss/tree"
)

// navItem represents a navigable item in the tree (file or directory header).
type navItem struct {
	fileIdx int    // >= 0 for files, -1 for directory headers
	dirPath string // non-empty for directory headers (path relative to repo root, no trailing slash)
	lineIdx int    // index into the rendered lines
}

// pathNode is an internal node for building a sorted path hierarchy before
// converting to lipgloss/tree for rendering. Leaf nodes are files (isDir=false);
// internal nodes are directories (isDir=true, fileIdx=-1).
type pathNode struct {
	name     string
	isDir    bool
	children []*pathNode
	fileIdx  int // index into []FileEntry; -1 for directories
}

// buildPathTree constructs a pathNode hierarchy from the given file entries.
func buildPathTree(files []FileEntry) *pathNode {
	root := &pathNode{isDir: true, fileIdx: -1}
	for i, f := range files {
		parts := strings.Split(filepath.ToSlash(f.RelPath), "/")
		insertPathNode(root, parts, i)
	}
	sortPathNode(root)
	return root
}

func insertPathNode(node *pathNode, parts []string, fileIdx int) {
	if len(parts) == 1 {
		node.children = append(node.children, &pathNode{
			name:    parts[0],
			isDir:   false,
			fileIdx: fileIdx,
		})
		return
	}
	dirName := parts[0]
	for _, child := range node.children {
		if child.isDir && child.name == dirName {
			insertPathNode(child, parts[1:], fileIdx)
			return
		}
	}
	dir := &pathNode{name: dirName, isDir: true, fileIdx: -1}
	node.children = append(node.children, dir)
	insertPathNode(dir, parts[1:], fileIdx)
}

// sortPathNode sorts children at every level: directories first, then files, both alphabetically.
func sortPathNode(node *pathNode) {
	sort.Slice(node.children, func(i, j int) bool {
		a, b := node.children[i], node.children[j]
		if a.isDir != b.isDir {
			return a.isDir
		}
		return strings.ToLower(a.name) < strings.ToLower(b.name)
	})
	for _, child := range node.children {
		if child.isDir {
			sortPathNode(child)
		}
	}
}

// tryCollapse checks if a directory node has exactly one file descendant through
// a chain of single-child directories. For example: git/ → git.go collapses to "git/git.go".
// Returns the collapsed display name, the fileIdx, and whether collapse is applicable.
func tryCollapse(node *pathNode) (name string, fileIdx int, ok bool) {
	if len(node.children) != 1 {
		return
	}
	child := node.children[0]
	if !child.isDir {
		return child.name, child.fileIdx, true
	}
	sub, idx, subOk := tryCollapse(child)
	if !subOk {
		return
	}
	return child.name + "/" + sub, idx, true
}

// BuildFileTree builds a lipgloss/tree.Tree from FileEntries using the lipgloss/tree library
// for all box-drawing rendering (├──, └──, │). Returns:
//   - rendered: plain-text rendered tree output (no ANSI codes)
//   - lineToFileIdx: maps each rendered line index → fileIdx (-1 for dir/root rows)
//   - navItems: navigable items (files and directory headers) in top-to-bottom visual order
func BuildFileTree(files []FileEntry, repoName string, collapsedDirs map[string]bool) (rendered string, lineToFileIdx []int, navItems []navItem) {
	root := buildPathTree(files)

	// Line 0 is the root header (repoName/)
	lineToFileIdx = append(lineToFileIdx, -1)

	t := lipglosstree.Root(repoName + "/")
	populateLipglossTree(t, root, &lineToFileIdx, &navItems, collapsedDirs, "")

	rendered = t.String()
	return
}

// populateLipglossTree recursively adds children to the lipgloss/tree node,
// recording the line index → fileIdx mapping in lockstep with the DFS traversal.
func populateLipglossTree(t *lipglosstree.Tree, node *pathNode, lineToFileIdx *[]int, navItems *[]navItem, collapsedDirs map[string]bool, pathPrefix string) {
	for _, child := range node.children {
		if child.isDir {
			fullPath := child.name
			if pathPrefix != "" {
				fullPath = pathPrefix + "/" + child.name
			}
			if collName, fileIdx, ok := tryCollapse(child); ok {
				// Collapsed single-child chain: one leaf line ("dir/file")
				lineIdx := len(*lineToFileIdx)
				t.Child(child.name + "/" + collName)
				*lineToFileIdx = append(*lineToFileIdx, fileIdx)
				*navItems = append(*navItems, navItem{fileIdx: fileIdx, lineIdx: lineIdx})
			} else if collapsedDirs[fullPath] {
				// User-collapsed directory: show header as leaf, no children rendered
				lineIdx := len(*lineToFileIdx)
				t.Child(child.name + "/")
				*lineToFileIdx = append(*lineToFileIdx, -1)
				*navItems = append(*navItems, navItem{fileIdx: -1, dirPath: fullPath, lineIdx: lineIdx})
			} else {
				// Expanded directory: subtree (dir header + children)
				lineIdx := len(*lineToFileIdx)
				*lineToFileIdx = append(*lineToFileIdx, -1) // dir header
				*navItems = append(*navItems, navItem{fileIdx: -1, dirPath: fullPath, lineIdx: lineIdx})
				subT := lipglosstree.Root(child.name + "/")
				populateLipglossTree(subT, child, lineToFileIdx, navItems, collapsedDirs, fullPath)
				t.Child(subT)
			}
		} else {
			lineIdx := len(*lineToFileIdx)
			t.Child(child.name)
			*lineToFileIdx = append(*lineToFileIdx, child.fileIdx)
			*navItems = append(*navItems, navItem{fileIdx: child.fileIdx, lineIdx: lineIdx})
		}
	}
}
