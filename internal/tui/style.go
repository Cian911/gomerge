package tui

import "github.com/charmbracelet/lipgloss"

var (
  mainViewStyle = lipgloss.NewStyle().
    PaddingRight(2).
    MarginRight(2).
    Border(
      lipgloss.RoundedBorder(),
      false,
      true,
      false,
      false,
    )
  helpViewStyle = lipgloss.NewStyle().
    Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
    Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
  dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"})
)
