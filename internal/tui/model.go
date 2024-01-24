package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/google/go-github/v57/github"
	"github.com/spf13/viper"

	"github.com/cian911/go-merge/pkg/gclient"
)

type model struct {
	width  int
	height int

	list     list.Model
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

	s := spinner.New()
	s.Spinner = spinner.Dot

	return &model{
		list:    l,
		spinner: s,

		keyMap: defaultKeyMappings(),

		gh:     client,
		loaded: false,
	}, nil
}
