package tui

import (
	"fmt"

	"github.com/google/go-github/v57/github"
)

type item struct {
	id        *int64
	number    *int
	state     *string
	title     *string
	body      *string
	createdAt *github.Timestamp
	updatedAt *github.Timestamp

  checked bool
}

func (i item) Title() string {
	return *i.title
}

func (i item) State() string {
	return *i.state
}

func (i item) Id() string {
	return fmt.Sprintf("%d", *i.id)
}

func (i item) Description() string {
	return ""
}

func (i item) FilterValue() string { return *i.title }
