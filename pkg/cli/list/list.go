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
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var org = ""
var repo = ""

func NewCommand() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "list",
		Short: "List all open pull request for a repo.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			orgRepo := viper.GetString("repo")
			token := viper.GetString("token")
			org, repo = parseOrgRepo(orgRepo)

			ghClient := gitclient.Client(token, ctx)
			pullRequests, _, err := ghClient.PullRequests.List(ctx, org, repo, nil)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			if len(pullRequests) == 0 {
				fmt.Println("No open pull requests found for given repository.")
				os.Exit(0)
			}

			table := printer.NewTable(os.Stdout,
				[]string{
					"PR",
					"State",
					"Title",
					"Repository",
					"Created",
				},
			)
			table = printer.HeaderStyle(table)

			prIds := []string{}

			for _, pr := range pullRequests {
				prIds = append(prIds, fmt.Sprintf("%d", *pr.Number))
				data := formatTable(pr)
				table = printer.SuccessStyle(table, data)
				// table.Append(data)
			}
			table.Render()

			prompt, selectedIds := selectPrIds(prIds)
			survey.AskOne(prompt, &selectedIds)

			for _, id := range selectedIds {
				prId, _ := strconv.Atoi(id)
				mergePullRequest(ghClient, ctx, org, repo, prId)
			}
		},
	}

	return
}

func formatTable(pr *github.PullRequest) []string {
	data := []string{
		fmt.Sprintf("#%s", printer.FormatID(pr.Number)),
		printer.FormatString(pr.State),
		printer.FormatString(pr.Title),
		fmt.Sprintf("%s/%s", org, repo),
		printer.FormatTime(pr.CreatedAt),
	}

	return data
}

func parseOrgRepo(repo string) (org, repository string) {
	str := strings.Split(repo, "/")

	if len(str) <= 1 {
		log.Fatal("You must pass your repo name like so: organization/repository to continue.")
		os.Exit(1)
	}

	org = str[0]
	repository = str[1]

	return
}

func selectPrIds(prIds []string) (*survey.MultiSelect, []string) {
	selectedIds := []string{}
	prompt := &survey.MultiSelect{
		Message: "Select which Pull Requests you would like to merge.",
		Options: prIds,
	}

	return prompt, selectedIds
}

func mergePullRequest(ghClient *github.Client, ctx context.Context, org, repo string, prId int) {
	result, _, err := ghClient.PullRequests.Merge(ctx, org, repo, prId, gitclient.DefaultCommitMsg(), nil)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("PR #%d: %v.", prId, *result.Message))
}
