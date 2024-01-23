package tui

import (
	"context"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type queryMsg struct {
	items []list.Item
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
			m.list, cmd = m.list.Update(msg)
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
		m.list.SetSize(m.width, m.height-helpBarHeight)

		// Detail View Size
		m.viewport = viewport.New(m.width, m.height-helpBarHeight)
		m.viewport.SetContent(m.mainViewportContent(m.viewport.Width))
	case queryMsg:
		m.list.SetItems(msg.items)
		m.list.SetWidth(m.width / 2)
		m.loaded = true
	default:
		// Do something as default
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)

		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// mainState denotes the state of the main view where we display
// the PR list to the user.
func (m model) mainState() tea.Cmd {
	// cmd  tea.Cmd
	var cmds []tea.Cmd

	m.list.Title = "Pull Requests"
	prs, _, err := m.gh.PullRequests.List(context.Background(), "Cian911", "gomerge-test", nil)
	if err != nil {
		log.Fatalf("Could not get PRs: %v", err)
	}
	idx := 0
	items := []list.Item{}

	for _, v := range prs {
		item := item{
			id:        v.ID,
			number:    v.Number,
			state:     v.State,
			title:     v.Title,
			body:      v.Body,
			createdAt: v.CreatedAt,
			updatedAt: v.UpdatedAt,
		}
		items = append(items, item)

		idx += 1
	}
	m.list.SetItems(items)

	return tea.Batch(cmds...)
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
			items []list.Item
		)

		m.list.Title = "Pull Requests"
		prs, _, err := m.gh.PullRequests.List(context.Background(), "Cian911", "gomerge-test", nil)
		if err != nil {
			log.Fatalf("Could not get PRs: %v", err)
		}
		idx := 0

		for _, v := range prs {
			item := item{
				id:        v.ID,
				number:    v.Number,
				state:     v.State,
				title:     v.Title,
				body:      v.Body,
				createdAt: v.CreatedAt,
				updatedAt: v.UpdatedAt,
			}
			items = append(items, item)

			idx += 1
		}

		return queryMsg{items: items}
	}
}

func stringPtr(str string) *string {
	return &str
}
