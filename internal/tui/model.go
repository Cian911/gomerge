package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/google/go-github/v57/github"
	"github.com/spf13/viper"

	"github.com/cian911/go-merge/pkg/gclient"
)

type model struct {
	width  int
	height int
  tableWidth int
  tableHeight int
  detailViewWidth int
  detailViewHeight int
  actionViewWidth int
  actionViewHeight int

  prs      []PullRequest
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

	s := spinner.New()
	s.Spinner = spinner.Dot

	return &model{
		spinner: s,

		keyMap: defaultKeyMappings(),

		gh:     client,
		loaded: false,
	}, nil
}
