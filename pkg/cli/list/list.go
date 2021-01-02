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

var org = ""
var repo = ""

func NewCommand() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "list",
		Short: "List all open pull request for a repository you wish to merge.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			orgRepo := viper.GetString("repo")
			token := viper.GetString("token")
			configFile := viper.GetString("config")

			configPresent := false

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
					mergePullRequest(ghClient, ctx, org, p[1], prId)
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
					mergePullRequest(ghClient, ctx, org, repo, prId)
				}
			}
		},
	}

	return
}

func promptAndFormat(pullRequests []*github.PullRequest, table *tablewriter.Table) []string {
	prIds := []string{}
	data := []string{}

	for _, pr := range pullRequests {
		prIds = append(prIds, fmt.Sprintf("%d | %s", *pr.Number, *pr.Head.Repo.Name))
		data = formatTable(pr, org, *pr.Head.Repo.Name)
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
