package tui

import (
	"context"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type queryMsg struct {
	items []PullRequest
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			switch {
			case key.Matches(msg, m.keyMap.Approve):
				// Approve PR
				rows := m.applySelection(Approve)
				m.table.SetRows(rows)
			case key.Matches(msg, m.keyMap.Merge):
				// Merge PR
				rows := m.applySelection(Merge)
				m.table.SetRows(rows)
			case key.Matches(msg, m.keyMap.Close):
				// Close PR
				rows := m.applySelection(Close)
				m.table.SetRows(rows)
			case key.Matches(msg, m.keyMap.Open):
				// Toggle action view
				m.actionViewSelected = !m.actionViewSelected
			case key.Matches(msg, m.keyMap.Help):
				// Toggle help view
				m.help.ShowAll = !m.help.ShowAll
			case key.Matches(msg, m.keyMap.Quit):
				// Quit
				return m, tea.Quit
			case key.Matches(msg, m.keyMap.Remove):
				// Check if item is actually selected
				if m.isSelected() {
					// Remove PR from selected list
					rows := m.applySelection(None)
					m.table.SetRows(rows)
				}
			}
		case tea.KeyUp, tea.KeyDown:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)

			m.selectedList, cmd = m.selectedList.Update(msg)
			m.viewport.GotoTop()
			// Clear viewport
			m.viewport.SetContent("")
			// Render new content
			m.viewport.SetContent(m.mainViewportContent(m.detailViewWidth))
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		helpBarHeight := lipgloss.Height(m.helpView())

		// Calculate widths for table and detail view
		tableWidth := int(float32(m.width) * 0.7)
		// Subtract a slight offset to the sidebar width
		detailViewWidth := m.width - (tableWidth)

		// Main View Size (Table)
		m.table.SetWidth(tableWidth)
		m.tableWidth = tableWidth
		if m.loaded {
			columns := adaptiveColumnWidths(m.tableWidth)
			m.table.SetColumns(columns)
		}
		// m.table.SetHeight(m.height - helpBarHeight)
		m.tableHeight = m.height - helpBarHeight

		// Set list dimensions
		// m.selectedList.SetWidth(m.width)
		// m.selectedList.SetHeight(m.height)

		// Detail View Size (Viewport for Sidebar)
		m.viewport.Width = detailViewWidth
		m.viewport.Height = m.height - helpBarHeight

		m.detailViewWidth = detailViewWidth
		m.detailViewHeight = m.viewport.Height

		m.actionViewWidth = m.detailViewWidth
		m.actionViewHeight = m.detailViewHeight
		m.viewport.SetContent(m.mainViewportContent(m.detailViewWidth))
	case queryMsg:
		m.prs = msg.items
		columns := adaptiveColumnWidths(m.tableWidth)
		rows := mapToTableRows(msg.items)

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithWidth(m.tableWidth),
			// table.WithHeight(m.tableHeight),
		)
		tStyle := table.DefaultStyles()
		tStyle.Header = tableStyle
		tStyle.Selected = tableSelectedStyle
		// tStyle.Cell = tableCellStyle
		t.SetStyles(tStyle)
		m.table = t
		m.viewport = viewport.New(m.width, m.height)
		m.viewport.SetContent(m.mainViewportContent(m.detailViewWidth))
		m.loaded = true
	default:
		m.selectedList, cmd = m.selectedList.Update(msg)
		cmds = append(cmds, cmd)
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
func (m *model) detailState(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	return cmd
}

// helpState denotes the state of the help view where we display
// helpful information and commands the user can action
func (m *model) helpState(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	return cmd
}

// actionState denotes the state of the action view where we display
// a list of actions taken by the user on a PR
func (m *model) actionState(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	return cmd
}

func (m *model) queryCmd() tea.Cmd {
	return func() tea.Msg {
		var (
			err   error
			items []PullRequest
		)

		prs, _, err := m.gh.PullRequests.List(context.Background(), "Cian911", "gomerge-test", nil)
		if err != nil {
			log.Fatalf("Could not get PRs: %v", err)
		}

		for _, v := range prs {
			item := PullRequest{
				Id:             fmt.Sprintf("%d", *v.Number),
				Title:          string(*v.Title),
				CreatedAt:      v.CreatedAt,
				UpdatedAt:      v.UpdatedAt,
				Repo:           "Cian911/gomerge-test",
				Author:         *v.User.Login,
				Head:           v.Head,
				Base:           v.Base,
				Mergeable:      v.GetMergeable(),
				Body:           v.GetBody(),
				State:          v.GetState(),
				MergeableState: v.GetMergeableState(),
				Additions:      v.GetAdditions(),
				Deletions:      v.GetDeletions(),
				Assignee:       v.GetAssignee().GetName(),
				Draft:          v.GetDraft(),
				selected:       false,
				choice:         None,
			}

			items = append(items, item)
		}

		return queryMsg{items: items}
	}
}

func adaptiveColumnWidths(tableWidth int) []table.Column {
	// Define minimum column width to prevent overly narrow columns
	minColumnWidth := 4

	proportions := []float32{0.01, 0.01, 0.4, 0.07, 0.20, 0.09}
	titles := []string{"", mergeGlyph, "PR Title", timeGlyph, repoGlyph, authorGlyph}

	columns := make([]table.Column, len(proportions))
	for i, proportion := range proportions {
		width := int(float32(tableWidth) * proportion)
		if width < minColumnWidth {
			width = minColumnWidth
		}
		columns[i] = table.Column{Title: titles[i], Width: width}
	}

	return columns
}

func (m *model) applySelection(choice Choice) []table.Row {
	itemId := m.table.SelectedRow()[1]
	m.prs[m.table.Cursor()].choice = choice
	// item := m.prs[m.table.Cursor()]
	i := selectedItem{
		title: "Tile",
		desc:  "This is a test description",
	}
	items := m.selectedList.Items()
	items = append(items, i)

	if choice == None {
		m.selectedList.RemoveItem(m.selectedList.Index())
	} else {
		// m.selectedList.InsertItem(0, i)
		m.selectedList.SetItems(items)
	}
	m.prs = updatedSelectedItem(m.prs, itemId)
	rows := mapToTableRows(m.prs)
	return rows
}

func (m *model) isSelected() bool {
	if m.prs[m.table.Cursor()].selected {
		return true
	}

	return false
}

func stringPtr(str string) *string {
	return &str
}
