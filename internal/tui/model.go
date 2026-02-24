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
	files      []FileEntry
	cursor     int
	gitStatus  map[string]string
	wsPath     string
	gitRoot    string
	cfg        *config.Config
	mode       mode
	input      textinput.Model
	confirmMsg string
	width      int
	height     int
	err        string
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

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return refreshMsg{}
	}
}

// refreshMsg triggers a reload of files + git status.
type refreshMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case refreshMsg:
		m = m.refresh()
		return m, nil

	case tea.KeyMsg:
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
		if len(m.files) > 0 {
			m.cursor = (m.cursor - 1 + len(m.files)) % len(m.files)
		}

	case key.Matches(msg, keys.Down):
		if len(m.files) > 0 {
			m.cursor = (m.cursor + 1) % len(m.files)
		}

	case key.Matches(msg, keys.Refresh):
		return m, func() tea.Msg { return refreshMsg{} }

	case key.Matches(msg, keys.Add):
		m.mode = modeAddInput
		m.input.SetValue("")
		m.input.Focus()
		return m, textinput.Blink

	case key.Matches(msg, keys.Delete):
		if len(m.files) == 0 {
			break
		}
		f := m.files[m.cursor]
		if f.Status != "" {
			m.confirmMsg = fmt.Sprintf("%s has uncommitted changes. Revert with git? [y/N/cancel]", f.RelPath)
			m.mode = modeDeleteConfirm
		} else {
			if err := store.Remove(m.wsPath, f.AbsPath); err != nil {
				m.err = err.Error()
			} else {
				m = m.refresh()
				if m.cursor >= len(m.files) && m.cursor > 0 {
					m.cursor = len(m.files) - 1
				}
			}
		}

	case key.Matches(msg, keys.Edit):
		if len(m.files) == 0 {
			break
		}
		f := m.files[m.cursor]
		if !f.Exists {
			m.err = "file does not exist on disk"
			break
		}
		editor := m.cfg.Editor
		cmd := exec.Command(editor, f.AbsPath)
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
			// Resolve relative to cwd, not git root
			abs, err := filepath.Abs(val)
			if err == nil {
				err = store.Add(m.wsPath, abs)
			}
			if err != nil {
				m.err = err.Error()
			} else {
				m.err = ""
			}
		}
		m.input.Blur()
		m.mode = modeNormal
		m = m.refresh()
		return m, nil

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
	if len(m.files) == 0 {
		m.mode = modeNormal
		return m, nil
	}
	f := m.files[m.cursor]

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
		m = m.refresh()
		if m.cursor >= len(m.files) && m.cursor > 0 {
			m.cursor = len(m.files) - 1
		}

	case "n", "N":
		if err := store.Remove(m.wsPath, f.AbsPath); err != nil {
			m.err = err.Error()
		}
		m.mode = modeNormal
		m.confirmMsg = ""
		m = m.refresh()
		if m.cursor >= len(m.files) && m.cursor > 0 {
			m.cursor = len(m.files) - 1
		}

	case "esc", "c", "q":
		m.mode = modeNormal
		m.confirmMsg = ""
	}

	return m, nil
}

// refresh syncs git status, auto-adds modified files, reloads the list.
func (m Model) refresh() Model {
	m.err = ""

	status, err := git.ModifiedFiles(m.gitRoot)
	if err != nil {
		m.err = err.Error()
		status = map[string]string{}
	}
	m.gitStatus = status

	// Auto-add all git-modified files (silent)
	for absPath := range status {
		_ = store.Add(m.wsPath, absPath)
	}

	// Reload working set
	paths, err := store.Load(m.wsPath)
	if err != nil {
		m.err = err.Error()
		paths = []string{}
	}

	m.files = make([]FileEntry, 0, len(paths))
	for _, abs := range paths {
		rel, err := filepath.Rel(m.gitRoot, abs)
		if err != nil {
			rel = abs
		}
		_, statErr := os.Stat(abs)
		m.files = append(m.files, FileEntry{
			AbsPath: abs,
			RelPath: rel,
			Status:  status[abs],
			Exists:  statErr == nil,
		})
	}

	return m
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	var sb strings.Builder

	// File list
	listHeight := m.height - 3 // reserve space for bottom bar + padding
	if m.mode == modeAddInput || m.mode == modeDeleteConfirm {
		listHeight -= 2
	}

	for i, f := range m.files {
		if i >= listHeight {
			break
		}

		// Build status badge
		var statusStr string
		switch f.Status {
		case "M":
			statusStr = styleStatusM.Render("M")
		case "A":
			statusStr = styleStatusA.Render("A")
		case "?":
			statusStr = styleStatusQ.Render("?")
		}

		// Build path display
		pathStr := f.RelPath
		var row string
		if !f.Exists {
			row = styleDim.Render(fmt.Sprintf("  %-*s  [gone]", m.width-12, pathStr))
		} else {
			padWidth := m.width - 8
			if padWidth < 1 {
				padWidth = 1
			}
			row = fmt.Sprintf("  %-*s  %s", padWidth, pathStr, statusStr)
		}

		if i == m.cursor {
			// Pad the selected row to full width for highlight
			raw := lipgloss.NewStyle().Render(row)
			padded := raw + strings.Repeat(" ", max(0, m.width-lipgloss.Width(raw)))
			sb.WriteString(styleSelected.Render(padded))
		} else {
			sb.WriteString(row)
		}
		sb.WriteString("\n")
	}

	// Fill remaining list space
	rendered := len(m.files)
	if rendered > listHeight {
		rendered = listHeight
	}
	for i := rendered; i < listHeight; i++ {
		sb.WriteString("\n")
	}

	// Inline prompt area
	switch m.mode {
	case modeAddInput:
		sb.WriteString("\n")
		sb.WriteString(stylePrompt.Render("  Add file: ") + m.input.View())
		sb.WriteString("\n")
	case modeDeleteConfirm:
		sb.WriteString("\n")
		sb.WriteString(stylePrompt.Render("  " + m.confirmMsg))
		sb.WriteString("\n")
	default:
		if m.err != "" {
			sb.WriteString("\n")
			sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("  Error: "+m.err))
			sb.WriteString("\n")
		} else {
			sb.WriteString("\n\n")
		}
	}

	// Status bar
	hint := "  j/k move  e edit  a add  d delete  r refresh  q quit"
	sb.WriteString(styleStatusBar.Render(hint))

	return sb.String()
}

// errMsg carries an error back into the update loop.
type errMsg struct{ err error }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
