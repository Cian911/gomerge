package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

func (m model) View() string {
	if !m.loaded {
		return m.spinner.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, m.mainView(), lipgloss.JoinVertical(lipgloss.Top, m.detailView(), m.actionView())),
		m.helpView(),
	)
}

func (m model) mainViewportContent(width int) string {
	var builder strings.Builder

	if m.loaded {
    title := detailViewTitleStyle.Render(m.table.SelectedRow()[2])
    mergable := detailViewStateStyle.Render(m.prs[m.table.Cursor()].IsMergable())
    state := lipgloss.NewStyle().
      Background(lipgloss.Color("#1CA4D6")).
      Width(lipgloss.Width(mergable)-5).
      Bold(true).
      Align(lipgloss.Center).
      Render(strings.ToUpper(m.prs[m.table.Cursor()].State))
    state = detailViewStateStyle.Render(state)
    stateRow := lipgloss.JoinHorizontal(lipgloss.Top, mergable, state)

    branch := detailViewBranchStyle.Render(m.prs[m.table.Cursor()].Branch())
    // description := detailViewDescriptionStyle.Render(m.prs[m.table.Cursor()].Description())
		
    builder.WriteString("\n\n\n")
    builder.WriteString(title)
    builder.WriteString("\n\n")
    builder.WriteString(branch)
    builder.WriteString("\n\n")
    builder.WriteString(stateRow)
    builder.WriteString("\n\n")
    // builder.WriteString(description)
	} else {
    defaultMsg := detailViewDefaultMsgStyle.Render("Content not loaded.")
		builder.WriteString(defaultMsg)
	}

	return wordwrap.String(builder.String(), width)
}

func (m model) mainView() string {
	return mainViewStyle.
    // Width(m.tableWidth).
    // Height(m.tableHeight).
    MaxWidth(m.tableWidth).
    MaxHeight(m.tableHeight).
    Render(m.table.View())
}

func (m model) detailView() string {
  styledDetail := lipgloss.NewStyle().
    // Width(m.detailViewWidth).
    MaxHeight(m.detailViewHeight).
    MaxWidth(m.detailViewWidth).
    BorderLeft(true).
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("63")).
    Render(m.viewport.View())

	return styledDetail
}

func (m model) helpView() string {
	help := "ctrl-m - merge, ctrl-a - approve, ctrl-c close"
	helpValue := helpViewStyle.Copy().Width(m.width).Render(help)

	helpViewBar := lipgloss.JoinHorizontal(lipgloss.Top, helpValue)
	return helpViewStyle.Width(m.width).Render(helpViewBar)
}

func (m model) actionView() string {
  items := "[x] #234 Issue 69 revert"
  actionView := lipgloss.NewStyle().
    // Width(m.actionViewWidth).
    // Height(10).
    MaxHeight(m.actionViewHeight).
    MaxWidth(m.actionViewWidth).
    Align(lipgloss.Left).
    BorderLeft(true).
    BorderTop(true).
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("63")).
    Render(items)
  
  return actionView
}
