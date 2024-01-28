package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v57/github"
	"github.com/spf13/viper"

	"github.com/cian911/go-merge/pkg/gclient"
)

type model struct {
	width  int
	height int

	list     list.Model
  table    table.Model
	viewport viewport.Model
	spinner  spinner.Model

	keyMap
	loaded bool

	// Github Client
	gh *github.Client
}

func New() (*model, error) {
	client := gclient.Client(viper.GetString("token"), context.Background(), false)
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Pull Requests"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

  columns := []table.Column{
    {Title: "Title", Width: 12},
    {Title: "Status", Width: 4},
    {Title: "Age", Width: 4},
    {Title: "Author", Width: 8},
    {Title: "Checks", Width: 4},
  }

  t := table.New(
    table.WithColumns(columns),
    table.WithFocused(true),
    table.WithHeight(1),
    table.WithWidth(32),
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

	s := spinner.New()
	s.Spinner = spinner.Dot

	return &model{
		list:    l,
		spinner: s,
    table: t,

		keyMap: defaultKeyMappings(),

		gh:     client,
		loaded: false,
	}, nil
}
