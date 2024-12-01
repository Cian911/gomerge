package tui

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v57/github"

	"github.com/cian911/go-merge/internal/utils"
)

const (
	None Choice = iota
	Merge
	Approve
	Close
)

type Choice int

type PullRequest struct {
	Id        string
	NodeId    string
	Title     string
	Repo      string
	Author    string
	State     string
	Body      string
	CreatedAt *github.Timestamp
	UpdatedAt *github.Timestamp
	Draft     bool
	Additions int
	Deletions int
	Assignee  string

	// Branch details
	Head *github.PullRequestBranch
	Base *github.PullRequestBranch

	Mergeable      bool
	MergeableState string

	choice   Choice
	selected bool
}

func (pr *PullRequest) Timestamp() string {
	return string(utils.TimeElapsed(pr.CreatedAt.Time))
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
	g, _ := glamour.NewTermRenderer(
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

type selectedItem struct {
	title, repo string
	choice      Choice
}

func (i selectedItem) Title() string       { return i.title }
func (i selectedItem) Description() string { return i.repo }
func (i selectedItem) FilterValue() string { return i.title }

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2).Bold(true)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000")).Bold(true)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	approvedStyle     = lipgloss.NewStyle().Background(lipgloss.Color("#12B910"))
	mergeStyle        = lipgloss.NewStyle().Background(lipgloss.Color("#8C2CB6"))
	closeStyle        = lipgloss.NewStyle().Background(lipgloss.Color("#C1122A"))
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(selectedItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s: %s", index+1, i.repo, i.title)

	switch i.choice {
	case Merge:
		str = mergeStyle.Render(str)
	case Approve:
		str = approvedStyle.Render(str)
	case Close:
		str = closeStyle.Render(str)
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
