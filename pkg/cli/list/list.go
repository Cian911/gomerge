package list

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cian911/go-merge/pkg/gitclient"
	"github.com/cian911/go-merge/pkg/printer"
	"github.com/cian911/go-merge/pkg/utils"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	org           = ""
	repo          = ""
	approveOnly   = false
	configPresent = false
)

// TODO: Refactor NewCommnd
func NewCommand() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "list",
		Short: "List all open pull request for a repository you wish to merge.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			orgRepo := viper.GetString("repo")
			token := viper.GetString("token")
			configFile := viper.GetString("config")
			approveOnly = viper.GetBool("approve")

			if len(configFile) > 0 {
				utils.ReadConfigFile(configFile)
				configPresent = true
			}

			if !configPresent && len(orgRepo) <= 0 {
				log.Fatal("You must pass either a config file or repository as argument to continue.")
				os.Exit(1)
			}

			ghClient := gitclient.Client(token, ctx)
			pullRequestsArray := []*github.PullRequest{}
			table := initTable()

			// If user has passed a config file
			if configPresent {
				org = viper.GetString("organization")

				for _, v := range viper.GetStringSlice("repositories") {
					pullRequests, _, err := ghClient.PullRequests.List(ctx, org, v, nil)

					if err != nil {
						log.Fatal(err)
						os.Exit(1)
					}

					// Use variadic notation to append to array here...
					pullRequestsArray = append(pullRequestsArray, pullRequests...)
				}

				if len(pullRequestsArray) == 0 {
					fmt.Println("No open pull requests found for configured repositories.")
					os.Exit(0)
				}

				selectedIds := promptAndFormat(pullRequestsArray, table)
				for _, id := range selectedIds {
					p := parsePrId(id)
					prId, _ := strconv.Atoi(p[0])
					if approveOnly {
						approvePullRequest(ghClient, ctx, org, repo, prId)
					} else {
						mergePullRequest(ghClient, ctx, org, p[1], prId)
					}
				}
			} else {
				org, repo = parseOrgRepo(orgRepo, configPresent)
				// if user has NOT passed a config file
				pullRequests, _, err := ghClient.PullRequests.List(ctx, org, repo, nil)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}

				if len(pullRequests) == 0 {
					fmt.Println("No open pull requests found for given repository.")
					os.Exit(0)
				}

				selectedIds := promptAndFormat(pullRequests, table)
				for _, id := range selectedIds {
					p := parsePrId(id)
					prId, _ := strconv.Atoi(p[0])
					if approveOnly {
						approvePullRequest(ghClient, ctx, org, repo, prId)
					} else {
						mergePullRequest(ghClient, ctx, org, repo, prId)
					}
				}
			}
		},
	}

	return
}

func promptAndFormat(pullRequests []*github.PullRequest, table *tablewriter.Table) []string {
	prIds := []string{}
	data := []string{}
	repoName := ""

	for _, pr := range pullRequests {
		if pr.Head.Repo == nil {
			repoName = "Forked Likely Repository Removed."
		} else {
			repoName = *pr.Head.Repo.Name
		}
		prIds = append(prIds, fmt.Sprintf("%d | %s", *pr.Number, repoName))
		data = formatTable(pr, org, repoName)
		table = printer.SuccessStyle(table, data)
	}
	table.Render()

	prompt, selectedIds := selectPrIds(prIds)
	survey.AskOne(prompt, &selectedIds)
	return selectedIds
}

func initTable() (table *tablewriter.Table) {
	table = printer.NewTable(os.Stdout,
		[]string{
			"PR",
			"State",
			"Title",
			"Repository",
			"Created",
		},
	)
	table = printer.HeaderStyle(table)
	return
}

func formatTable(pr *github.PullRequest, org, repo string) []string {
	data := []string{
		fmt.Sprintf("#%s", printer.FormatID(pr.Number)),
		printer.FormatString(pr.State),
		printer.FormatString(pr.Title),
		fmt.Sprintf("%s/%s", org, repo),
		printer.FormatTime(pr.CreatedAt),
	}

	return data
}

func parseOrgRepo(repo string, configPresent bool) (org, repository string) {
	str := strings.Split(repo, "/")

	if len(str) <= 1 {
		log.Fatal("You must pass your repo name like so: organization/repository to continue.")
		os.Exit(1)
	}

	org = str[0]
	repository = str[1]

	return
}

func parsePrId(prId string) []string {
	str := strings.Split(strings.ReplaceAll(prId, " ", ""), "|")
	return str
}

func selectPrIds(prIds []string) (*survey.MultiSelect, []string) {
	selectedIds := []string{}
	msg := ""
	switch approveOnly {
	case true:
		msg = "Select which Pull Requests you would like to approve."
	default:
		msg = "Select which Pull Requests you would like to merge."
	}
	prompt := &survey.MultiSelect{
		Message: msg,
		Options: prIds,
	}

	return prompt, selectedIds
}

// TODO: Move to gitclient pkg
func mergePullRequest(ghClient *github.Client, ctx context.Context, org, repo string, prId int) {
	result, _, err := ghClient.PullRequests.Merge(ctx, org, repo, prId, gitclient.DefaultCommitMsg(), nil)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("PR #%d: %v.", prId, *result.Message))
}

// TODO: Move to gitclient pkg
func approvePullRequest(ghClient *github.Client, ctx context.Context, org, repo string, prId int) {
	// Create review
	t := fmt.Sprintf(`PR #%d has been approved by [GoMerge](https://github.com/Cian911/gomerge) tool. :rocket:`, prId)
	e := "APPROVE"
	reviewRequest := &github.PullRequestReviewRequest{
		Body:  &t,
		Event: &e,
	}
	review, _, err := ghClient.PullRequests.CreateReview(ctx, org, repo, prId, reviewRequest)
	if err != nil {
		//TODO: Parse error to check if user tried to approve their own PR..
		log.Fatalf("Could not approve pull request, did you try to approve your on pull request? - %v", err)
	}

	fmt.Printf("PR #%d: %v\n", prId, *review.State)
}
