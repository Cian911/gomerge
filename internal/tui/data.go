package tui

import (
	"fmt"
	"regexp"
	"strings"

  "github.com/charmbracelet/glamour"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/cian911/go-merge/internal/utils"
	"github.com/google/go-github/v57/github"
)

const (
  None Choice = iota
  Merge
  Approve
  Close
)

type Choice int

type PullRequest struct {
  Id string
  NodeId string
  Title string
  Repo string
  Author string
  State string
  Body string
  CreatedAt *github.Timestamp
  UpdatedAt *github.Timestamp
  Draft bool
  Additions int
  Deletions int
  Assignee string

  // Branch details
  Head *github.PullRequestBranch
  Base *github.PullRequestBranch

  Mergeable bool
  MergeableState string

  choice Choice
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
      Render(fmt.Sprintf("%s MERGEABLE", mergeGlyph))
  }

  return lipgloss.NewStyle().
    Background(lipgloss.Color("#D80E07")).
    Align(lipgloss.Center).
    Bold(true).
    Width(12).
    Render(fmt.Sprintf("%s MERGEABLE", mergeGlyph))
}

func (pr *PullRequest) Description(width int) string {
  g,_ := glamour.NewTermRenderer(
		glamour.WithWordWrap(width),
	)
  regex := regexp.MustCompile("(?U)<!--(.|[[:space:]])*-->")
	body := regex.ReplaceAllString(pr.Body, "")

	regex = regexp.MustCompile(`((\n)+|^)([^\r\n]*\|[^\r\n]*(\n)?)+`)
	body = regex.ReplaceAllString(body, "")

	body = strings.TrimSpace(body)
  rendered, _ := g.Render(body)
  return rendered
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
