package tui

import (
	"context"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type queryMsg struct {
	items []table.Row
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			switch {
			case key.Matches(msg, m.keyMap.approve):
				// Approve PR
			case key.Matches(msg, m.keyMap.merge):
				// Merge PR
			case key.Matches(msg, m.keyMap.close):
				// Close PR
			}
		case tea.KeyCtrlC:
			cmd = tea.Quit
			cmds = append(cmds, cmd)
		case tea.KeyUp, tea.KeyDown:
			m.table, cmd = m.table.Update(msg)
			m.viewport.GotoTop()
			m.viewport.SetContent(m.mainViewportContent(m.viewport.Width))
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		helpBarHeight := lipgloss.Height(m.helpView())

		// Main View Size
		// mainViewWidth := cast.ToInt(0.2 * float64(m.width))
		// mainViewSize := mainViewWidth - mainViewStyle.GetHorizontalFrameSize()
		// m.table.SetSize(m.width, m.height-helpBarHeight)
		m.table.SetWidth(m.width)
		m.table.SetHeight(m.height - helpBarHeight)

		// Detail View Size
		m.viewport = viewport.New(m.width, m.height-helpBarHeight)
		m.viewport.SetContent(m.mainViewportContent(m.viewport.Width))
	case queryMsg:
		// m.table.SetRows(msg.items)
		// m.table.SetWidth(m.width)
		// helpBarHeight := lipgloss.Height(m.helpView())
		// m.table.SetHeight(m.height - helpBarHeight)
		// m.list.Select(0)
		// m.list.SetWidth(m.width / 2)
		// m.table.SetWidth(m.width/ 2)
		columns := []table.Column{
			{Title: "Title", Width: 50},
			{Title: "Status", Width: 8},
			{Title: "Age", Width: 12},
			{Title: "Author", Width: 12},
			{Title: "Checks", Width: 15},
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(msg.items),
			table.WithFocused(true),
			// table.WithHeight(10),
			table.WithWidth(100),
		)
		tableStyle := table.DefaultStyles()
		tableStyle.Header = tableStyle.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		tableStyle.Selected = tableStyle.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(tableStyle)
		m.table = t
		m.viewport.SetContent(m.mainViewportContent(m.viewport.Width))
		m.loaded = true
	default:
		// Do something as default
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)

		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// detailView denotes the state of the detail view where we display
// further detailed information for the selected PR
func (m model) detailState(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	return cmd
}

// helpState denotes the state of the help view where we display
// helpful information and commands the user can action
func (m model) helpState(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	return cmd
}

// actionState denotes the state of the action view where we display
// a list of actions taken by the user on a PR
func (m model) actionState(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	return cmd
}

func (m model) queryCmd() tea.Cmd {
	return func() tea.Msg {
		var (
			err   error
			items []table.Row
		)

		prs, _, err := m.gh.PullRequests.List(context.Background(), "Cian911", "gomerge-test", nil)
		if err != nil {
			log.Fatalf("Could not get PRs: %v", err)
		}
		idx := 0

		for _, v := range prs {
			// item := table.Row{
			// 	id:        v.ID,
			// 	number:    v.Number,
			// 	state:     v.State,
			// 	title:     v.Title,
			// 	body:      v.Body,
			// 	createdAt: v.CreatedAt,
			// 	updatedAt: v.UpdatedAt,
			// }
			item := table.Row{
				string(*v.Title),
				string(*v.State),
				"2 weeks ago",
				string(v.User.GetID()),
				"passing",
			}
			//    item := table.Row{
			// 	string(*v.Title),
			// 	string(*v.State),
			// 	v.CreatedAt.String(),
			// 	v.User.GetName(),
			// 	fmt.Sprintf("%t", v.Mergeable),
			// }
			items = append(items, item)

			idx += 1
		}

		return queryMsg{items: items}
	}
}

func stringPtr(str string) *string {
	return &str
}
