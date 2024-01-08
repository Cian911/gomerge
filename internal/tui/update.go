package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

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
		}
	default:
		// Do something as default
	}

	return m, cmd
}

// mainState denotes the state of the main view where we display
// the PR list to the user.
func (m model) mainState(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	return cmd
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
