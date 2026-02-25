package tui

import "github.com/charmbracelet/lipgloss"

var (
	styleSelected = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("237")).
			Foreground(lipgloss.Color("255"))

	styleNormal = lipgloss.NewStyle()

	styleDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Strikethrough(true)

	styleDir = lipgloss.NewStyle().Foreground(lipgloss.Color("111")) // light blue for directories

	styleStatusM = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // yellow
	styleStatusA = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))  // green
	styleStatusQ = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))  // cyan

	styleSearchMatch = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	styleStatusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)

	styleFilePath = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	stylePrompt = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	styleBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("237"))
)
