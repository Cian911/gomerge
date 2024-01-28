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
		lipgloss.JoinHorizontal(lipgloss.Top, m.mainView(), m.detailView()),
		m.helpView(),
	)
}

func (m model) mainViewportContent(width int) string {
	var builder strings.Builder
	divider := dividerStyle.Render(strings.Repeat("-", width)) + "\n"

	if it := m.list.SelectedItem(); it != nil {
		keyType := fmt.Sprintf("KeyType: %s\n", it.(item).Title())
		key := fmt.Sprintf("Key: \n%s\n", it.(item).State())
		value := fmt.Sprintf("Value: \n%s\n", it.(item).Id())
		builder.WriteString(keyType)
		builder.WriteString(divider)
		builder.WriteString(key)
		builder.WriteString(divider)
		builder.WriteString(value)
	} else {
		builder.WriteString("No item selected")
	}

	return wordwrap.String(builder.String(), width)
}

func (m model) mainView() string {
  // var builder strings.Builder
  //
  // for _, listItem := range m.list.Items() {
  //   it := listItem.(item)
  //
  //   checkbox := "[ ]"
  //   if it.checked {
  //     checkbox = "[x]"
  //   }
  //
  //   checkbox = checkboxStyle.Render(checkbox)
  //   title := titleStyle.Render(fmt.Sprintf("%s", it.Title()))
  //   state := stateStyle.Render(fmt.Sprintf("%s", it.State()))
  //
  //   line := lipgloss.JoinHorizontal(lipgloss.Top, checkbox, title, state)
  //   builder.WriteString(lineStyle.Render(line))
  //   builder.WriteString("\n")
  // }
  //
  // return mainViewStyle.Render(builder.String())
  return mainViewStyle.Render(m.table.View())
}

func (m model) detailView() string {
	return m.viewport.View()
}

func (m model) helpView() string {
	help := "ctrl-m - merge, ctrl-a - approve, ctrl-c close"
	helpValue := helpViewStyle.Copy().Width(m.width).Render(help)

	helpViewBar := lipgloss.JoinHorizontal(lipgloss.Top, helpValue)
	return helpViewStyle.Width(m.width).Render(helpViewBar)
}

func (m model) actionView() string { return "" }
