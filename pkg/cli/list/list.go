package list

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cian911/go-merge/pkg/gitclient"
	"github.com/cian911/go-merge/pkg/printer"
	"github.com/cian911/go-merge/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	org           = ""
	repo          = ""
	approveOnly   = false
	configPresent = false
)

const (
	TokenEnvVar = "GITHUB_TOKEN"
)

func GetMergeMethod() githubv4.PullRequestMergeMethod {
	method := viper.GetString("merge-method")
	switch method {
	case "merge":
		return githubv4.PullRequestMergeMethodMerge
	case "rebase":
		return githubv4.PullRequestMergeMethodRebase
	case "squash":
		return githubv4.PullRequestMergeMethodSquash
	}
	if method != "" {
		log.Fatalf("Unknown merge method %s. Please use one of the following: merge, rebase, squash", method)
	}
	return githubv4.PullRequestMergeMethodMerge
}

// TODO: Refactor NewCommnd
func NewCommand() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "list",
		Short: "List all open pull request for a repository you wish to merge.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			orgRepo := viper.GetString("repo")
			configFile := viper.GetString("config")
			approveOnly = viper.GetBool("approve")
			mergeMethod := GetMergeMethod()
			flagToken := viper.GetString("token")
			skip := viper.GetBool("skip")
			closePr := viper.GetBool("close")
			enterpriseUrl := viper.GetString("enterprise-base-url")
			delay := viper.GetInt("delay")

			if len(configFile) > 0 {
				utils.ReadConfigFile(configFile)
				configPresent = true
			}

			if !configPresent && len(orgRepo) <= 0 {
				log.Fatal("You must pass either a config file or repository as argument to continue.")
			}
			configToken := viper.GetString("token")

			token, err := getToken(flagToken, configToken)
			if err != nil {
				log.Fatal(err)
			}

			isEnterprise := false
			if len(enterpriseUrl) > 0 {
				isEnterprise = true
			}

			ghClient := gitclient.ClientV4(token, ctx, isEnterprise)

			pullRequestsArray := []*gitclient.PullRequest{}
			table := initTable()

			var org string
			var repositories []string = nil

			if configPresent {
				org = viper.GetString("organization")
				if len(viper.GetStringSlice("repositories")) > 0 {
					repositories = viper.GetStringSlice("repositories")
				} else {
					repositories = append(repositories, "")
				}
			} else {
				parts := strings.Split(orgRepo, "/")

				if len(parts) == 1 {
					org = parts[0]
					repositories = append(repositories, "")
				} else if len(parts) == 2 {
					org = parts[0]
					repositories = append(repositories, parts[1])
				} else {
					log.Fatal("You must pass your repo name like so: organization/repository to continue.")
				}
			}

			for _, v := range repositories {
				pullRequests, err := gitclient.GetPullRequests(ghClient, ctx, org, v)
				if err != nil {
					log.Fatal(err)
				}

				// Use variadic notation to append to array here...
				pullRequestsArray = append(pullRequestsArray, pullRequests...)
			}

			if len(pullRequestsArray) == 0 {
				fmt.Println("No open pull requests found for configured repositories.")
				os.Exit(0)
			}

			selectedIds := promptAndFormat(pullRequestsArray, table)
			for i, id := range selectedIds {

				if approveOnly {
					gitclient.ApprovePullRequest(ghClient, ctx, id, skip)
				} else if closePr {
					gitclient.ClosePullRequest(ghClient, ctx, id, skip)
				} else {
					// delay between merges to allow other active PRs to get synced
					if i > 0 {
						time.Sleep(time.Duration(delay) * time.Second)
					}
					gitclient.MergePullRequest(ghClient, ctx, id, &mergeMethod, skip)

				}
			}
		},
	}

	return
}

func promptAndFormat(pullRequests []*gitclient.PullRequest, table *tablewriter.Table) []githubv4.ID {
	prIds := []string{}

	for _, pr := range pullRequests {
		prIds = append(prIds, fmt.Sprintf("%d | %s/%s", pr.Number, pr.RepositoryOwner, pr.RepositoryName))

		data := formatTable(pr)
		if len(data) == 0 {
			// If there is an issue with the pr, skip
			continue
		}
		table = printer.SuccessStyle(table, data)
	}
	table.Render()

	prompt, selectedIds := selectPrIds(prIds)
	survey.AskOne(prompt, &selectedIds)
	indices := []githubv4.ID{}
	for _, id := range selectedIds {
		for i, prId := range prIds {
			if id == prId {
				indices = append(indices, pullRequests[i].ID)
				break
			}
		}
	}
	return indices
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

func statusIcon(state string) (icon string) {
	switch state {
	case "SUCCESS":
		icon = "✅"
	case "IN_PROGRESS":
		icon = "🟠"
	case "FAILURE":
		icon = "❌"
	default:
		icon = ""
	}

	return
}

func formatTable(pr *gitclient.PullRequest) (data []string) {
	data = []string{
		fmt.Sprintf("#%d", pr.Number),
		fmt.Sprintf("%s %s", pr.State, statusIcon(pr.StatusRollup)),
		pr.Title,
		fmt.Sprintf("%s/%s", pr.RepositoryOwner, pr.RepositoryName),
		printer.FormatTime(&pr.CreatedAt),
	}

	return
}

func parseOrgRepo(repo string, configPresent bool) (org, repository string) {
	str := strings.Split(repo, "/")

	if len(str) <= 1 {
		log.Fatal("You must pass your repo name like so: organization/repository to continue.")
	}

	org = str[0]
	repository = str[1]

	return
}

func getToken(flag, config string) (str string, err error) {
	if flag != str {
		return flag, nil
	}
	if config != str {
		return config, nil
	}
	if env, ok := os.LookupEnv(TokenEnvVar); ok {
		return env, nil
	}

	err = fmt.Errorf("you must pass a github token to continue")
	return
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

func commitMsg(ctx context.Context, msg string) context.Context {
	if len(msg) != 0 {
		return context.WithValue(ctx, "message", msg)
	}

	return context.WithValue(ctx, "message", gitclient.DefaultApproveMsg())
}
