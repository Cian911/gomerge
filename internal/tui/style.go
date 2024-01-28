package tui

import "github.com/charmbracelet/lipgloss"

var (
	// mainViewStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, true).Padding(0, 2)
  mainViewStyle = lipgloss.NewStyle().
	  // BorderStyle(lipgloss.NormalBorder()).
    PaddingTop(5).
	  BorderForeground(lipgloss.Color("240"))

	detailViewStyle = lipgloss.NewStyle().
			Padding(1, 1).
			Border(lipgloss.RoundedBorder())

	helpViewStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#24273a", Dark: "#181926"}).
			Align(lipgloss.Center)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"})

  titleStyle = lipgloss.NewStyle().Width(70)
  stateStyle = lipgloss.NewStyle().Width(20)
  checkboxStyle = lipgloss.NewStyle().Width(4)
  lineStyle = lipgloss.NewStyle().
    PaddingTop(1).
    PaddingBottom(1).
    Border(lipgloss.ThickBorder(), true, false, false, false)

  selectedItemStyle = lipgloss.NewStyle().
    Background(lipgloss.Color("#FF00FF")).
    Foreground(lipgloss.Color("#FFF")).
    Bold(true)
)
