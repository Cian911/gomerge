package tui

import (
	"fmt"
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
    title := detailViewTitleStyle.Width(m.detailViewWidth).Render(m.table.SelectedRow()[2])
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
    description := detailViewDescriptionStyle.Render(m.prs[m.table.Cursor()].Description(m.detailViewWidth))
		
    builder.WriteString("\n")
    builder.WriteString(title)
    builder.WriteString("\n\n")
    builder.WriteString(branch)
    builder.WriteString("\n\n")
    builder.WriteString(stateRow)
    builder.WriteString("\n\n")
    builder.WriteString(description)
	} else {
    defaultMsg := detailViewDefaultMsgStyle.Width(m.detailViewWidth).Render("Content not loaded.")
		builder.WriteString(defaultMsg)
	}

	return wordwrap.String(builder.String(), m.detailViewWidth)
}

func (m model) actionViewportContent(width int) string {
	var builder strings.Builder
  selectedPrStyle := lipgloss.NewStyle().
    Align(lipgloss.Left).
    Foreground(lipgloss.Color("#fff"))  

	if m.loaded {
    for _, pr := range m.prs {
      if pr.selected {
        selected := fmt.Sprintf("#%s %s\n", pr.Id, pr.Title)
        s := actionViewBackgroundStyles(selected, pr.choice, m.actionViewWidth)
        builder.WriteString(selectedPrStyle.Render(s))
      }
    }
	} else {
    defaultMsg := detailViewDefaultMsgStyle.Render("Content not loaded.")
		builder.WriteString(defaultMsg)
	}

	return wordwrap.String(builder.String(), width)
}

func (m model) mainView() string {
	return mainViewStyle.
    MaxWidth(m.tableWidth).
    MaxHeight(m.tableHeight).
    Render(m.table.View())
}

func (m model) detailView() string {
  styledDetail := lipgloss.NewStyle().
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
  actionView := lipgloss.NewStyle().
    MaxHeight(m.actionViewHeight).
    Height(m.tableHeight/2).
    MaxWidth(m.actionViewWidth).
    Align(lipgloss.Left).
    BorderLeft(true).
    BorderTop(true).
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("63")).
    Render(m.actionViewportContent(m.actionViewWidth))
  
  return actionView
}

func actionViewBackgroundStyles(str string, choice Choice, width int) string {
  listStyle := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#fff")).
    Width(width).
    Padding(0).
    Align(lipgloss.Left).
    BorderTop(true).
    // BorderForeground(lipgloss.Color("#fff")).
    BorderForeground(lipgloss.Color("240")).
    BorderStyle(lipgloss.NormalBorder())

  approvedStyle := listStyle.Copy().Background(lipgloss.Color("#12B910"))
  mergeStyle := listStyle.Copy().Background(lipgloss.Color("#8C2CB6"))
  closeStyle := listStyle.Copy().Background(lipgloss.Color("#C1122A"))

  switch choice {
    case Merge:
      return mergeStyle.Render(fmt.Sprintf("%s - %s", mergeGlyph, str))
    case Approve:
      return approvedStyle.Render(fmt.Sprintf("%s - %s", successGlyph, str))
    case Close:
    return closeStyle.Render(fmt.Sprintf("%s - %s", errorGlyph, str))
    default:
      return str
  }
}
