package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
  mainViewStyle = lipgloss.NewStyle().
	  BorderForeground(lipgloss.Color("240")).
    BorderRight(true)

  tableStyle = table.DefaultStyles().
    Header.
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("240")).
    BorderBottom(true).
    Bold(true)

  tableSelectedStyle = table.DefaultStyles().
    Selected.
    Foreground(lipgloss.Color("#fff")).
		Background(lipgloss.Color("57")).
		Bold(true)   

  tableCellStyle = table.DefaultStyles().
    Cell.
    Padding(1)

	detailViewStyle = lipgloss.NewStyle().
			Padding(1, 1).
      BorderLeft(true).
      BorderRight(true).
      BorderStyle(lipgloss.NormalBorder()).
      BorderForeground(lipgloss.Color("63")).
      BorderBottom(true).BorderBottomBackground(lipgloss.Color("#fff"))

  detailViewBranchStyle = lipgloss.NewStyle().
    Italic(true).
    Bold(true).
    Foreground(lipgloss.Color("#CBCBCB")).
    MarginLeft(3)

  detailViewTitleStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#fff")).
    MarginLeft(3).
    Align(lipgloss.Left)

  detailViewStateStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#fff")).
    MarginLeft(3).
    Align(lipgloss.Left)

  detailViewDescriptionStyle = lipgloss.NewStyle().
    MarginLeft(3).
    Foreground(lipgloss.Color("#fff"))

  detailViewDefaultMsgStyle = lipgloss.NewStyle().
    Align(lipgloss.Center).
    Bold(true).
    Foreground(lipgloss.Color("#fff"))

	helpViewStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#24273a", Dark: "#181926"}).
			Align(lipgloss.Center)

  mergeGlyph = ""
  successGlyph = ""
  errorGlyph = "󰅙"
  waitingGlyph = "󰗨"
  prGlyph = ""
  timeGlyph = ""
  changeReqeuestGlyp = "󰌑"
  authorGlyph =""
  repoGlyph = ""
)
