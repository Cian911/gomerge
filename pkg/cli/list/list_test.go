package list

import (
	"testing"
	"time"

	"github.com/cian911/go-merge/pkg/printer"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"github.com/stretchr/testify/assert"
)

func TestParseOrgRepo(t *testing.T) {
	t.Run("It returns a valid tuple when no config is present", func(t *testing.T) {
		repo := "Cian911/syncwave"
		configPresent := false

		want1 := "Cian911"
		want2 := "syncwave"

		got1, got2 := parseOrgRepo(repo, configPresent)

		if got1 != want1 || got2 != want2 {
			t.Errorf("got1: %s, got2: %s, want1: %s, want2: %s", got1, got2, want1, want2)
		}
	})
}

func TestInitTable(t *testing.T) {
	t.Run("It returns a tablewriter pointer", func(t *testing.T) {
		got := initTable()
		want := &tablewriter.Table{}

		assert.IsType(t, got, want)
	})
}

func TestFormatTable(t *testing.T) {
	var (
		org  = "Cian911"
		repo = "syncwave"
	)

	t.Run("It returns a string array", func(t *testing.T) {
		number := 1
		state := "#open"
		title := "My Pr"
		createdAt := time.Now()

		pr := &github.PullRequest{
			Number:    &number,
			State:     &state,
			Title:     &title,
			CreatedAt: &createdAt,
		}

		got := formatTable(pr, org, repo)
		want := []string{
			"#1",
			"#open",
			"My Pr",
			"Cian911/syncwave",
			printer.FormatTime(pr.CreatedAt),
		}

		assert.Equal(t, got, want)
	})

	t.Run("It logs failure when attrs are not present in pr struct", func(t *testing.T) {
		state := "#open"
		title := "My Pr"

		pr := &github.PullRequest{
			State: &state,
			Title: &title,
		}

		got := formatTable(pr, org, repo)
		want := []string(nil)

		assert.Equal(t, got, want)
	})
}
