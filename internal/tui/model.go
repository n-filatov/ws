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
)

type mode int

const (
	modeNormal mode = iota
	modeAddInput
	modeDeleteConfirm
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
	files         []FileEntry
	cursor        int    // index into treeOrder
	scrollOffset  int    // first visible rendered line index
	treeRendered  string // lipgloss/tree rendered output (plain text)
	lineToFileIdx []int  // lineToFileIdx[i] = fileIdx for line i, or -1
	treeOrder     []int  // fileIdx values in top-to-bottom visual order
	wsPath        string
	gitRoot       string
	cfg           *config.Config
	mode          mode
	input         textinput.Model
	confirmMsg    string
	loading       bool
	width         int
	height        int
	err           string
}

// New creates a new Model. Call Init() to start.
func New(wsPath, gitRoot string, cfg *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "path/to/file"
	ti.CharLimit = 512

	return Model{
		wsPath:  wsPath,
		gitRoot: gitRoot,
		cfg:     cfg,
		input:   ti,
	}
}

// selectedEntry returns the FileEntry currently under the cursor, or nil if none.
func (m Model) selectedEntry() *FileEntry {
	if len(m.treeOrder) == 0 || m.cursor >= len(m.treeOrder) {
		return nil
	}
	return &m.files[m.treeOrder[m.cursor]]
}

// listHeight returns how many tree lines fit in the current terminal.
func (m Model) listHeight() int {
	h := m.height - 3
	if m.mode == modeAddInput || m.mode == modeDeleteConfirm {
		h -= 2
	}
	if h < 1 {
		h = 1
	}
	return h
}

// lineForCursor returns the rendered line index for the currently selected file.
func (m Model) lineForCursor() int {
	if len(m.treeOrder) == 0 || m.cursor >= len(m.treeOrder) {
		return 0
	}
	target := m.treeOrder[m.cursor]
	for i, idx := range m.lineToFileIdx {
		if idx == target {
			return i
		}
	}
	return 0
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

// ── Messages ─────────────────────────────────────────────────────────────────

// refreshMsg signals that an async refresh should start.
type refreshMsg struct{}

// refreshResultMsg carries the completed result of an async refresh.
type refreshResultMsg struct {
	files         []FileEntry
	treeRendered  string
	lineToFileIdx []int
	treeOrder     []int
	err           string
}

// errMsg carries an error from a background operation.
type errMsg struct{ err error }

// ── Async I/O ─────────────────────────────────────────────────────────────────

// doRefresh runs all git/store I/O in a goroutine and returns the result as a Cmd.
// Nothing in this function touches the model — it is pure data fetching.
func doRefresh(wsPath, gitRoot string) tea.Cmd {
	return func() tea.Msg {
		result := refreshResultMsg{}

		status, err := git.ModifiedFiles(gitRoot)
		if err != nil {
			result.err = err.Error()
			status = map[string]string{}
		}

		// Auto-add all git-modified files (silent, best-effort)
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

		repoName := filepath.Base(gitRoot)
		result.treeRendered, result.lineToFileIdx, result.treeOrder =
			BuildFileTree(result.files, repoName)

		return result
	}
}

// ── Bubbletea lifecycle ───────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	m.loading = true
	return doRefresh(m.wsPath, m.gitRoot)
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
		return m, doRefresh(m.wsPath, m.gitRoot)

	case refreshResultMsg:
		m.loading = false
		m.err = msg.err
		m.files = msg.files
		m.treeRendered = msg.treeRendered
		m.lineToFileIdx = msg.lineToFileIdx
		m.treeOrder = msg.treeOrder
		if m.cursor >= len(m.treeOrder) && m.cursor > 0 {
			m.cursor = len(m.treeOrder) - 1
		}
		m = m.withScrollClamped()
		return m, nil

	case errMsg:
		m.err = msg.err.Error()
		return m, nil

	case tea.KeyMsg:
		// Ctrl+C always quits regardless of mode
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
		}
	}
	return m, nil
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.Up):
		if len(m.treeOrder) > 0 {
			m.cursor = (m.cursor - 1 + len(m.treeOrder)) % len(m.treeOrder)
			m = m.withScrollClamped()
		}

	case key.Matches(msg, keys.Down):
		if len(m.treeOrder) > 0 {
			m.cursor = (m.cursor + 1) % len(m.treeOrder)
			m = m.withScrollClamped()
		}

	case key.Matches(msg, keys.Refresh):
		return m, func() tea.Msg { return refreshMsg{} }

	case key.Matches(msg, keys.Add):
		m.mode = modeAddInput
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

	// Footer occupies 3 lines: prompt/error area (2) + status bar (1).
	listHeight := m.height - 3
	if m.mode == modeAddInput || m.mode == modeDeleteConfirm {
		listHeight -= 2
	}
	if listHeight < 1 {
		listHeight = 1
	}

	if m.loading && len(m.files) == 0 {
		// Show a simple loading indicator on first load
		sb.WriteString(styleStatusBar.Render("  Loading…") + "\n")
		for i := 1; i < listHeight; i++ {
			sb.WriteString("\n")
		}
	} else {
		// Determine which file is selected
		selectedFileIdx := -1
		if len(m.treeOrder) > 0 && m.cursor < len(m.treeOrder) {
			selectedFileIdx = m.treeOrder[m.cursor]
		}

		// Render tree lines starting from scrollOffset
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

			displayLine := "  " + line
			isSelected := fileIdx >= 0 && fileIdx == selectedFileIdx

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

		// Fill remaining list space
		for rendered < listHeight {
			sb.WriteString("\n")
			rendered++
		}
	}

	// Prompt / error area
	switch m.mode {
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
	hint := "  j/k move  e edit  a add  d delete  r refresh  q quit"
	if m.loading {
		hint = "  Loading…"
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
