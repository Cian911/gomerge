package tui

import "github.com/charmbracelet/lipgloss"

var (
	mainViewStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, true).Padding(0, 2)

	detailViewStyle = lipgloss.NewStyle().
			Padding(1, 1).
			Border(lipgloss.RoundedBorder())

	helpViewStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#24273a", Dark: "#181926"}).
			MarginTop(1).
			MarginBottom(1).
			Align(lipgloss.Center)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"})
)
