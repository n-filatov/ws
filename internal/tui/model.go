package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nfilatov/ws/internal/config"
	"github.com/nfilatov/ws/internal/git"
	"github.com/nfilatov/ws/internal/store"
	"github.com/sahilm/fuzzy"
)

type mode int

const (
	modeNormal mode = iota
	modeAddInput
	modeDeleteConfirm
	modeSearch
)

// FileEntry represents a file in the working set.
type FileEntry struct {
	AbsPath string
	RelPath string
	Status  string // "M", "A", "?", ""
	Exists  bool
}

// Model is the Bubbletea model for the TUI.
type Model struct {
	allFiles       []FileEntry // all files from last refresh (source of truth)
	files          []FileEntry // currently displayed files (filtered or all)
	cursor         int         // index into navItems
	scrollOffset   int         // first visible rendered line index
	treeRendered   string      // lipgloss/tree rendered output (plain text)
	lineToFileIdx  []int       // lineToFileIdx[i] = fileIdx for line i, or -1
	navItems       []navItem
	collapsedDirs  map[string]bool // set of dir paths (no trailing slash) that are collapsed
	searchQuery    string          // active filter; "" means no filter
	matchHighlights map[string][]int // RelPath → matched char positions in basename
	wsPath         string
	gitRoot        string
	cfg            *config.Config
	mode           mode
	input          textinput.Model
	confirmMsg     string
	loading        bool
	width          int
	height         int
	err            string
}

// New creates a new Model. Call Init() to start.
func New(wsPath, gitRoot string, cfg *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "path/to/file"
	ti.CharLimit = 512

	return Model{
		wsPath:        wsPath,
		gitRoot:       gitRoot,
		cfg:           cfg,
		input:         ti,
		collapsedDirs: make(map[string]bool),
	}
}

// selectedNavItem returns the navItem currently under the cursor, or nil if none.
func (m Model) selectedNavItem() *navItem {
	if len(m.navItems) == 0 || m.cursor >= len(m.navItems) {
		return nil
	}
	item := m.navItems[m.cursor]
	return &item
}

// selectedEntry returns the FileEntry currently under the cursor, or nil if none/on a dir.
func (m Model) selectedEntry() *FileEntry {
	item := m.selectedNavItem()
	if item == nil || item.fileIdx < 0 {
		return nil
	}
	return &m.files[item.fileIdx]
}

// listHeight returns how many tree lines fit in the current terminal.
func (m Model) listHeight() int {
	h := m.height - 3
	if m.mode == modeAddInput || m.mode == modeDeleteConfirm || m.mode == modeSearch {
		h -= 2
	}
	if h < 1 {
		h = 1
	}
	return h
}

// lineForCursor returns the rendered line index for the currently selected item.
func (m Model) lineForCursor() int {
	if len(m.navItems) == 0 || m.cursor >= len(m.navItems) {
		return 0
	}
	return m.navItems[m.cursor].lineIdx
}

// withScrollClamped adjusts scrollOffset so the selected line stays visible.
func (m Model) withScrollClamped() Model {
	lh := m.listHeight()
	line := m.lineForCursor()

	if line < m.scrollOffset {
		m.scrollOffset = line
	} else if line >= m.scrollOffset+lh {
		m.scrollOffset = line - lh + 1
	}

	maxScroll := len(m.lineToFileIdx) - lh
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scrollOffset > maxScroll {
		m.scrollOffset = maxScroll
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
	return m
}

// applySearchFilter filters m.allFiles by the current searchQuery, rebuilds the tree,
// and recomputes matchHighlights. Call this whenever allFiles or searchQuery changes.
func (m Model) applySearchFilter() Model {
	if m.searchQuery == "" {
		m.files = m.allFiles
		m.matchHighlights = nil
	} else {
		relPaths := make([]string, len(m.allFiles))
		for i, f := range m.allFiles {
			relPaths[i] = f.RelPath
		}
		matches := fuzzy.Find(m.searchQuery, relPaths)

		m.files = make([]FileEntry, 0, len(matches))
		m.matchHighlights = make(map[string][]int, len(matches))
		for _, match := range matches {
			f := m.allFiles[match.Index]
			m.files = append(m.files, f)
			m.matchHighlights[f.RelPath] = basenameMatchPositions(match.MatchedIndexes, f.RelPath)
		}
	}

	repoName := filepath.Base(m.gitRoot)
	m.treeRendered, m.lineToFileIdx, m.navItems = BuildFileTree(m.files, repoName, m.collapsedDirs)
	if len(m.navItems) > 0 && m.cursor >= len(m.navItems) {
		m.cursor = len(m.navItems) - 1
	}
	m = m.withScrollClamped()
	return m
}

// basenameMatchPositions filters fuzzy match positions to only those within the
// basename of relPath and returns them relative to the start of the basename.
func basenameMatchPositions(matchedIndexes []int, relPath string) []int {
	basename := filepath.Base(relPath)
	basenameStart := len(relPath) - len(basename)
	var positions []int
	for _, idx := range matchedIndexes {
		if idx >= basenameStart {
			positions = append(positions, idx-basenameStart)
		}
	}
	return positions
}

// highlightBasename returns the basename string with matched character positions
// styled with styleSearchMatch.
func highlightBasename(basename string, positions []int) string {
	if len(positions) == 0 {
		return basename
	}
	posSet := make(map[int]bool, len(positions))
	for _, p := range positions {
		posSet[p] = true
	}
	var sb strings.Builder
	for i, ch := range basename {
		if posSet[i] {
			sb.WriteString(styleSearchMatch.Render(string(ch)))
		} else {
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}

// rebuildTree reconstructs the rendered tree from the current files and collapsed state.
// Used after fold/unfold so we don't need an async refresh.
func (m Model) rebuildTree() (tea.Model, tea.Cmd) {
	repoName := filepath.Base(m.gitRoot)
	m.treeRendered, m.lineToFileIdx, m.navItems = BuildFileTree(m.files, repoName, m.collapsedDirs)
	if len(m.navItems) > 0 && m.cursor >= len(m.navItems) {
		m.cursor = len(m.navItems) - 1
	}
	m = m.withScrollClamped()
	return m, nil
}

// ── Messages ─────────────────────────────────────────────────────────────────

type refreshMsg struct{}

type refreshResultMsg struct {
	files         []FileEntry
	treeRendered  string
	lineToFileIdx []int
	navItems      []navItem
	err           string
}

type errMsg struct{ err error }

// ── Async I/O ─────────────────────────────────────────────────────────────────

func doRefresh(wsPath, gitRoot string, collapsedDirs map[string]bool) tea.Cmd {
	return func() tea.Msg {
		result := refreshResultMsg{}

		status, err := git.ModifiedFiles(gitRoot)
		if err != nil {
			result.err = err.Error()
			status = map[string]string{}
		}

		for absPath := range status {
			_ = store.Add(wsPath, absPath)
		}

		paths, err := store.Load(wsPath)
		if err != nil {
			if result.err == "" {
				result.err = err.Error()
			}
			paths = []string{}
		}

		result.files = make([]FileEntry, 0, len(paths))
		for _, abs := range paths {
			rel, relErr := filepath.Rel(gitRoot, abs)
			if relErr != nil {
				rel = abs
			}
			_, statErr := os.Stat(abs)
			result.files = append(result.files, FileEntry{
				AbsPath: abs,
				RelPath: rel,
				Status:  status[abs],
				Exists:  statErr == nil,
			})
		}

		// Tree is built after filter is applied in refreshResultMsg handler.
		return result
	}
}

// ── Bubbletea lifecycle ───────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	m.loading = true
	return doRefresh(m.wsPath, m.gitRoot, m.collapsedDirs)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.withScrollClamped()
		return m, nil

	case refreshMsg:
		m.loading = true
		m.err = ""
		return m, doRefresh(m.wsPath, m.gitRoot, m.collapsedDirs)

	case refreshResultMsg:
		m.loading = false
		m.err = msg.err
		m.allFiles = msg.files
		m = m.applySearchFilter() // applies active query (or copies allFiles → files)
		return m, nil

	case errMsg:
		m.err = msg.err.Error()
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		switch m.mode {
		case modeNormal:
			return m.updateNormal(msg)
		case modeAddInput:
			return m.updateAddInput(msg)
		case modeDeleteConfirm:
			return m.updateDeleteConfirm(msg)
		case modeSearch:
			return m.updateSearch(msg)
		}
	}
	return m, nil
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		// Esc clears active search filter before quitting
		if msg.String() == "esc" && m.searchQuery != "" {
			m.searchQuery = ""
			m.input.SetValue("")
			return m.applySearchFilter(), nil
		}
		return m, tea.Quit

	case key.Matches(msg, keys.Up):
		if len(m.navItems) > 0 {
			m.cursor = (m.cursor - 1 + len(m.navItems)) % len(m.navItems)
			m = m.withScrollClamped()
		}

	case key.Matches(msg, keys.Down):
		if len(m.navItems) > 0 {
			m.cursor = (m.cursor + 1) % len(m.navItems)
			m = m.withScrollClamped()
		}

	case key.Matches(msg, keys.Right):
		item := m.selectedNavItem()
		if item != nil && item.fileIdx < 0 && item.dirPath != "" && m.collapsedDirs[item.dirPath] {
			delete(m.collapsedDirs, item.dirPath)
			return m.rebuildTree()
		}

	case key.Matches(msg, keys.Left):
		item := m.selectedNavItem()
		if item != nil && item.fileIdx < 0 && item.dirPath != "" && !m.collapsedDirs[item.dirPath] {
			m.collapsedDirs[item.dirPath] = true
			return m.rebuildTree()
		}

	case key.Matches(msg, keys.Search):
		m.mode = modeSearch
		m.input.Placeholder = "fuzzy search…"
		m.input.SetValue(m.searchQuery) // preserve any previous query
		m.input.CursorEnd()
		m.input.Focus()
		return m, textinput.Blink

	case key.Matches(msg, keys.Refresh):
		return m, func() tea.Msg { return refreshMsg{} }

	case key.Matches(msg, keys.Add):
		m.mode = modeAddInput
		m.input.Placeholder = "path/to/file"
		m.input.SetValue("")
		m.input.Focus()
		return m, textinput.Blink

	case key.Matches(msg, keys.Delete):
		f := m.selectedEntry()
		if f == nil {
			break
		}
		if f.Status != "" {
			m.confirmMsg = fmt.Sprintf("%s has uncommitted changes. Revert with git? [y/N/cancel]", f.RelPath)
			m.mode = modeDeleteConfirm
		} else {
			if err := store.Remove(m.wsPath, f.AbsPath); err != nil {
				m.err = err.Error()
			} else {
				return m, func() tea.Msg { return refreshMsg{} }
			}
		}

	case key.Matches(msg, keys.Edit):
		f := m.selectedEntry()
		if f == nil {
			break
		}
		if !f.Exists {
			m.err = "file does not exist on disk"
			break
		}
		cmd := exec.Command(m.cfg.Editor, f.AbsPath)
		return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return errMsg{err}
			}
			return refreshMsg{}
		})
	}

	return m, nil
}

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.searchQuery = ""
		m.input.SetValue("")
		m.input.Blur()
		m.mode = modeNormal
		return m.applySearchFilter(), nil

	case tea.KeyEnter:
		m.searchQuery = m.input.Value()
		m.input.Blur()
		m.mode = modeNormal
		return m.applySearchFilter(), nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.searchQuery = m.input.Value()
	return m.applySearchFilter(), cmd
}

func (m Model) updateAddInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		val := strings.TrimSpace(m.input.Value())
		if val != "" {
			abs, err := filepath.Abs(val)
			if err == nil {
				err = store.Add(m.wsPath, abs)
			}
			if err != nil {
				m.err = err.Error()
			}
		}
		m.input.Blur()
		m.mode = modeNormal
		return m, func() tea.Msg { return refreshMsg{} }

	case tea.KeyEsc:
		m.input.Blur()
		m.mode = modeNormal
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) updateDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f := m.selectedEntry()
	if f == nil {
		m.mode = modeNormal
		return m, nil
	}

	switch msg.String() {
	case "y", "Y":
		if err := git.Checkout(f.AbsPath); err != nil {
			m.err = err.Error()
		}
		if err := store.Remove(m.wsPath, f.AbsPath); err != nil {
			m.err = err.Error()
		}
		m.mode = modeNormal
		m.confirmMsg = ""
		return m, func() tea.Msg { return refreshMsg{} }

	case "n", "N":
		if err := store.Remove(m.wsPath, f.AbsPath); err != nil {
			m.err = err.Error()
		}
		m.mode = modeNormal
		m.confirmMsg = ""
		return m, func() tea.Msg { return refreshMsg{} }

	case "esc", "c", "q":
		m.mode = modeNormal
		m.confirmMsg = ""
	}

	return m, nil
}

// ── View ─────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	var sb strings.Builder

	listHeight := m.height - 3
	if m.mode == modeAddInput || m.mode == modeDeleteConfirm || m.mode == modeSearch {
		listHeight -= 2
	}
	if listHeight < 1 {
		listHeight = 1
	}

	if m.loading && len(m.allFiles) == 0 {
		sb.WriteString(styleStatusBar.Render("  Loading…") + "\n")
		for i := 1; i < listHeight; i++ {
			sb.WriteString("\n")
		}
	} else if len(m.files) == 0 && m.searchQuery != "" {
		// No matches
		sb.WriteString(styleDim.Render("  no matches") + "\n")
		for i := 1; i < listHeight; i++ {
			sb.WriteString("\n")
		}
	} else {
		selectedLineIdx := -1
		if len(m.navItems) > 0 && m.cursor < len(m.navItems) {
			selectedLineIdx = m.navItems[m.cursor].lineIdx
		}

		lineToDirPath := make(map[int]string)
		for _, item := range m.navItems {
			if item.fileIdx < 0 && item.dirPath != "" {
				lineToDirPath[item.lineIdx] = item.dirPath
			}
		}

		lines := strings.Split(m.treeRendered, "\n")
		rendered := 0
		start := m.scrollOffset
		if start < 0 {
			start = 0
		}
		for i := start; i < len(lines) && rendered < listHeight; i++ {
			line := lines[i]

			fileIdx := -1
			if i < len(m.lineToFileIdx) {
				fileIdx = m.lineToFileIdx[i]
			}

			dirPath := lineToDirPath[i]
			isCollapsed := dirPath != "" && m.collapsedDirs[dirPath]

			displayLine := "  " + line
			if isCollapsed {
				displayLine += " ▶"
			}

			// Apply fuzzy highlights to the filename portion
			if fileIdx >= 0 && len(m.matchHighlights) > 0 {
				relPath := m.files[fileIdx].RelPath
				if positions, ok := m.matchHighlights[relPath]; ok && len(positions) > 0 {
					basename := filepath.Base(relPath)
					highlighted := highlightBasename(basename, positions)
					if idx := strings.LastIndex(displayLine, basename); idx >= 0 {
						displayLine = displayLine[:idx] + highlighted + displayLine[idx+len(basename):]
					}
				}
			}

			isSelected := i == selectedLineIdx

			var entry *FileEntry
			if fileIdx >= 0 {
				entry = &m.files[fileIdx]
			}

			if isSelected {
				suffix := ""
				if entry != nil {
					if !entry.Exists {
						suffix = "  [gone]"
					} else if entry.Status != "" {
						suffix = "  " + entry.Status
					}
				}
				plain := displayLine + suffix
				padded := plain + strings.Repeat(" ", max(0, m.width-lipgloss.Width(plain)))
				sb.WriteString(styleSelected.Render(padded))
			} else if entry != nil && !entry.Exists {
				sb.WriteString(styleDim.Render(displayLine + "  [gone]"))
			} else if entry != nil && entry.Status != "" {
				var badge string
				switch entry.Status {
				case "M":
					badge = styleStatusM.Render("M")
				case "A":
					badge = styleStatusA.Render("A")
				case "?":
					badge = styleStatusQ.Render("?")
				}
				sb.WriteString(displayLine + "  " + badge)
			} else if fileIdx == -1 && i > 0 {
				sb.WriteString(styleDir.Render(displayLine))
			} else {
				sb.WriteString(displayLine)
			}
			sb.WriteString("\n")
			rendered++
		}

		for rendered < listHeight {
			sb.WriteString("\n")
			rendered++
		}
	}

	// Prompt / error area
	switch m.mode {
	case modeSearch:
		sb.WriteString(stylePrompt.Render("  Search: ") + m.input.View() + "\n")
		sb.WriteString("\n")
	case modeAddInput:
		sb.WriteString(stylePrompt.Render("  Add file: ") + m.input.View() + "\n")
		sb.WriteString("\n")
	case modeDeleteConfirm:
		sb.WriteString(stylePrompt.Render("  "+m.confirmMsg) + "\n")
		sb.WriteString("\n")
	default:
		if m.err != "" {
			sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("  Error: "+m.err) + "\n")
		} else {
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Status bar
	var hint string
	switch {
	case m.loading:
		hint = "  Loading…"
	case m.mode == modeSearch:
		hint = "  type to filter  enter confirm  esc clear"
	case m.searchQuery != "":
		hint = fmt.Sprintf("  %d match(es)  esc to clear  j/k move  e edit  a add  d delete  q quit", len(m.files))
	default:
		hint = "  j/k move  ←/→ fold  / search  e edit  a add  d delete  r refresh  q quit"
	}
	sb.WriteString(styleStatusBar.Render(hint))

	return sb.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
