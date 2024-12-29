package list

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cian911/go-merge/pkg/gitclient"
	"github.com/cian911/go-merge/pkg/printer"
	"github.com/cian911/go-merge/pkg/utils"
)

var (
	org           = ""
	repo          = ""
	approveOnly   = false
	configPresent = false
)

const (
	TokenEnvVar    = "GITHUB_TOKEN"
	STATUS_SUCCESS = 0
	STATUS_WAITING = 1
	STATUS_FAILED  = 2
)

func NewCommand() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "list",
		Short: "List all open pull request for a repository you wish to merge.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			orgRepo := viper.GetString("repo")
			labels := getLabels()
			configFile := viper.GetString("config")
			approveOnly = viper.GetBool("approve")
			mergeMethod := getMergeMethod()
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
				log.Fatal(
					"You must pass either a config file or repository as argument to continue.",
				)
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
				pullRequests, err := gitclient.GetPullRequests(ghClient, ctx, org, v, &labels)
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
			for i, pr := range selectedIds {
				if approveOnly {
					gitclient.ApprovePullRequest(ghClient, ctx, pr, skip)
				} else if closePr {
					gitclient.ClosePullRequest(ghClient, ctx, pr, skip)
				} else {
					// delay between merges to allow other active PRs to get synced
					if i > 0 {
						time.Sleep(time.Duration(delay) * time.Second)
					}
					if pr.NeedsReview {
						gitclient.ApprovePullRequest(ghClient, ctx, pr, skip)
					}
					gitclient.MergePullRequest(ghClient, ctx, pr, &mergeMethod, skip)

				}
			}
		},
	}

	return
}

func promptAndFormat(
	pullRequests []*gitclient.PullRequest,
	table *tablewriter.Table,
) (selectedPullRequests []*gitclient.PullRequest) {
	prIds := []string{}

	for _, pr := range pullRequests {
		prIds = append(
			prIds,
			fmt.Sprintf("%d | %s/%s", pr.Number, pr.RepositoryOwner, pr.RepositoryName),
		)

		data, status := formatTable(pr)
		if len(data) == 0 {
			// If there is an issue with the pr, skip
			continue
		}
		switch status {
		case STATUS_SUCCESS:
			table = printer.SuccessStyle(table, data)
		case STATUS_WAITING:
			table = printer.WaitingStyle(table, data)
		case STATUS_FAILED:
			table = printer.ErrorStyle(table, data)
		}
	}
	table.Render()

	prompt, selectedIds := selectPrIds(prIds)
	survey.AskOne(prompt, &selectedIds)
	selectedPullRequests = make([]*gitclient.PullRequest, len(selectedIds))
	for idIndex, id := range selectedIds {
		for prIndex, prId := range prIds {
			if id == prId {
				selectedPullRequests[idIndex] = pullRequests[prIndex]
				break
			}
		}
	}
	return
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

func statusIcon(state string) (icon string, status int) {
	switch state {
	case "SUCCESS":
		icon = ""
		status = STATUS_SUCCESS
	case "IN_PROGRESS":
		icon = ""
		status = STATUS_WAITING
	case "FAILURE":
		icon = "󰅙"
		status = STATUS_FAILED
	default:
		icon = ""
	}

	return
}

func formatTable(pr *gitclient.PullRequest) (data []string, status int) {
	icon, status := statusIcon(pr.StatusRollup)
	data = []string{
		fmt.Sprintf("#%d", pr.Number),
		fmt.Sprintf("%s %s", pr.State, icon),
		pr.Title,
		fmt.Sprintf("%s/%s", pr.RepositoryOwner, pr.RepositoryName),
		printer.FormatTime(&pr.CreatedAt),
	}

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

func getMergeMethod() githubv4.PullRequestMergeMethod {
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
		log.Fatalf(
			"Unknown merge method %s. Please use one of the following: merge, rebase, squash",
			method,
		)
	}
	return githubv4.PullRequestMergeMethodMerge
}

func getLabels() (labels []githubv4.String) {
	raw_labels := viper.GetStringSlice("label")
	labels = make([]githubv4.String, len(raw_labels))
	for i, label := range raw_labels {
		labels[i] = githubv4.String(label)
	}
	return
}
