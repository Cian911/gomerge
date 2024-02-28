package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

func (m *model) View() string {
	if !m.loaded {
		return m.spinner.View()
	}

	if m.actionViewSelected {
		m.selectedList.SetHeight(20)
		return m.selectedList.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, m.mainView(), m.detailView()),
		m.helpView(),
	)
}

func (m *model) mainViewportContent(width int) string {
	var builder strings.Builder

	if m.loaded {
		title := detailViewTitleStyle.Width(m.detailViewWidth).Render(m.table.SelectedRow()[2])
		mergable := detailViewStateStyle.Render(m.prs[m.table.Cursor()].IsMergable())
		state := lipgloss.NewStyle().
			Background(lipgloss.Color("#1CA4D6")).
			Width(lipgloss.Width(mergable) - 5).
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

func (m *model) actionViewportContent(width int) string {
	if m.loaded {
		return lipgloss.NewStyle().
			Height(m.height).
			Width(m.width).
			Render(m.selectedList.View()) + "\n"
	} else {
		//   defaultMsg := detailViewDefaultMsgStyle.Render("No pull requests selected.")
		// builder.WriteString(defaultMsg)
	}

	return ""
}

func (m *model) mainView() string {
	return mainViewStyle.
		MaxWidth(m.tableWidth).
		MaxHeight(m.tableHeight).
		Render(m.table.View())
}

func (m *model) detailView() string {
	styledDetail := lipgloss.NewStyle().
		MaxHeight(m.detailViewHeight).
		MaxWidth(m.detailViewWidth).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(m.viewport.View())

	return styledDetail
}

func (m *model) helpView() string {
	helpValue := helpViewStyle.Copy().Width(m.width).Render(m.help.View(defaultKeyMappings()))

	helpViewBar := lipgloss.JoinHorizontal(lipgloss.Top, helpValue)
	return helpViewStyle.Width(m.width).Render(helpViewBar)
}

func (m *model) actionView() string {
	m.selectedList.Title = "Pull Requests"
	return m.selectedList.View()
}
