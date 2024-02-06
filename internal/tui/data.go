package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/cian911/go-merge/internal/utils"
	"github.com/google/go-github/v57/github"
)

type PullRequest struct {
  Id string
  Title string
  Repo string
  Author string
  State string
  Body string
  CreatedAt *github.Timestamp
  UpdatedAt *github.Timestamp

  // Branch details
  Head *github.PullRequestBranch
  Base *github.PullRequestBranch

  Mergeable bool
  selected bool
}

func (pr *PullRequest) Timestamp() string {
  return string (utils.TimeElapsed(pr.CreatedAt.Time))
}

func (pr *PullRequest) PrAuthor() string {
  return fmt.Sprintf("%s", pr.Author)
}

func (pr *PullRequest) HeadBranch() string {
  return fmt.Sprintf("%s", *pr.Head.Ref)
}

func (pr *PullRequest) BaseBranch() string {
  return fmt.Sprintf("%s", *pr.Base.Ref)
}

func (pr *PullRequest) Branch() string {
  return fmt.Sprintf("%s -> %s", pr.HeadBranch(), pr.BaseBranch())
}

func (pr *PullRequest) IsMergable() string {
  if pr.Mergeable {
    return lipgloss.NewStyle().
      Background(lipgloss.Color("#18CF15")).
      Align(lipgloss.Center).
      Bold(true).
      Width(12).
      Render("MERGEABLE")
  }

  return lipgloss.NewStyle().
    Background(lipgloss.Color("#D80E07")).
    Align(lipgloss.Center).
    Bold(true).
    Width(12).
    Render("MERGEABLE")
}

func (pr *PullRequest) Description() string {
  return pr.Body
}

func updatedSelectedItem(prs []PullRequest, itemId string) []PullRequest {
  for i := range prs {
    if prs[i].Id == itemId {
      prs[i].selected = !prs[i].selected
    }
  }

  return prs
}

func mapToTableRows(prs []PullRequest) []table.Row {
  rows := []table.Row{}

  for _, v := range prs {
    selected := "[]"
    if v.selected {
      selected = "[x]"
    }

    row := table.Row{
      selected,
      v.Id,
      v.Title,
      v.Timestamp(),
      v.Repo,
      v.PrAuthor(),
    }

    rows = append(rows, row)
  }

  return rows
}
