package list

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cian911/go-merge/pkg/gitclient"
	"github.com/cian911/go-merge/pkg/printer"
	"github.com/cian911/go-merge/pkg/utils"
	"github.com/google/go-github/v45/github"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	org           = ""
	repo          = ""
	approveOnly   = false
	configPresent = false
	mergeMethod   = "merge"
)

const (
	TokenEnvVar = "GITHUB_TOKEN"
)

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
			mergeMethod := viper.GetString("merge-method")
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

			ghClient := gitclient.Client(token, ctx, isEnterprise)
			pullRequestsArray := []*github.PullRequest{}
			table := initTable()

			// If user has passed a config file
			if configPresent {
				org = viper.GetString("organization")

				for _, v := range viper.GetStringSlice("repositories") {
					pullRequests, _, err := ghClient.PullRequests.List(ctx, org, v, nil)
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
				for x, id := range selectedIds {
					p := parsePrId(id)
					prId, _ := strconv.Atoi(p[0])
					if approveOnly {
						gitclient.ApprovePullRequest(ghClient, ctx, org, p[1], prId, skip)
					} else if closePr {
						gitclient.ClosePullRequest(ghClient, ctx, org, p[1], prId, pullRequestsArray[x])
					} else {
						gitclient.MergePullRequest(ghClient, ctx, org, p[1], prId, mergeMethod, skip)

						// delay between merges to allow other active PRs to get synced
						time.Sleep(time.Duration(delay) * time.Second)
					}
				}
			} else {
				org, repo = parseOrgRepo(orgRepo, configPresent)
				// if user has NOT passed a config file
				pullRequests, _, err := ghClient.PullRequests.List(ctx, org, repo, nil)
				if err != nil {
					log.Fatal(err)
				}

				if len(pullRequests) == 0 {
					fmt.Println("No open pull requests found for given repository.")
					os.Exit(0)
				}

				selectedIds := promptAndFormat(pullRequests, table)
				for x, id := range selectedIds {
					p := parsePrId(id)
					prId, _ := strconv.Atoi(p[0])
					if approveOnly {
						gitclient.ApprovePullRequest(ghClient, ctx, org, repo, prId, skip)
					} else if closePr {
						gitclient.ClosePullRequest(ghClient, ctx, org, repo, prId, pullRequests[x])
					} else {
						gitclient.MergePullRequest(ghClient, ctx, org, repo, prId, mergeMethod, skip)

						// delay between merges to allow other active PRs to get synced
						time.Sleep(time.Duration(delay) * time.Second)
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
		prIds = append(prIds, fmt.Sprintf("%d | %s | %s", *pr.Number, repoName, *pr.User.Login))
		data = formatTable(pr, org, repoName)
		if len(data) == 0 {
			// If there is an issue with the pr, skip
			continue
		}
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
			"Author",
			"Created",
		},
	)
	table = printer.HeaderStyle(table)
	return
}

func formatTable(pr *github.PullRequest, org, repo string) (data []string) {
	if (pr.Number == nil) || (pr.State == nil) || (pr.Title == nil) || (pr.CreatedAt == nil) {
		return
	}
	data = []string{
		fmt.Sprintf("#%s", printer.FormatID(pr.Number)),
		printer.FormatString(pr.State),
		printer.FormatString(pr.Title),
		fmt.Sprintf("%s/%s", org, repo),
		*pr.User.Login,
		printer.FormatTime(pr.CreatedAt),
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

func parsePrId(prId string) []string {
	str := strings.Split(strings.ReplaceAll(prId, " ", ""), "|")
	return str
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
