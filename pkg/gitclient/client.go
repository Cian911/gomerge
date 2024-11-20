package gitclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/shurcooL/githubv4"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type PullRequest struct {
	RepositoryOwner string
	RepositoryName  string
	Number          int
	Title           string
	State           string
	CreatedAt       time.Time
	ID              githubv4.ID
	StatusRollup    string
}

type Commits struct {
	Nodes []struct {
		Commit struct {
			StatusCheckRollup struct {
				State githubv4.String
			}
		}
	}
}

type Repository struct {
	NameWithOwner githubv4.String
	Owner         struct {
		Login githubv4.String
	}
	Name         githubv4.String
	PullRequests struct {
		Nodes []struct {
			Number    githubv4.Int
			Title     githubv4.String
			State     githubv4.String
			Url       githubv4.URI
			CreatedAt githubv4.DateTime
			ID        githubv4.ID
			Commits   Commits `graphql:"commits(last:1)"`
		}
	} `graphql:"pullRequests(states: [OPEN], first: $maxPullRequests, orderBy: {field: CREATED_AT, direction: DESC})"`
}

func Client(githubToken string, ctx context.Context, isEnterprise bool) (client *github.Client) {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: githubToken,
		},
	)

	tokenContext := oauth2.NewClient(ctx, tokenSource)

	if isEnterprise {
		baseUrl := viper.GetString("enterprise-base-url")
		c, err := github.NewEnterpriseClient(baseUrl, baseUrl, tokenContext)

		if err != nil {
			log.Fatalf("Could not auth enterprise client: %v", err)
		}

		client = c
	} else {
		client = github.NewClient(tokenContext)
	}

	return
}

func ClientV4(githubToken string, ctx context.Context, isEnterprise bool) (client *githubv4.Client) {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: githubToken,
		},
	)

	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	if isEnterprise {
		baseUrl := viper.GetString("enterprise-base-url")
		client = githubv4.NewEnterpriseClient(baseUrl, httpClient)
	} else {
		client = githubv4.NewClient(httpClient)

	}

	return
}

func GetPullRequests(client *githubv4.Client, ctx context.Context, owner string, repo string) ([]*PullRequest, error) {
	vars := map[string]interface{}{
		"maxPullRequests": githubv4.Int(100),
		"owner":           githubv4.String(owner),
	}

	repos := []Repository{}

	if repo == "" {
		var q struct {
			RepositoryOwner struct {
				Repositories struct {
					Nodes []Repository
				} `graphql:"repositories(first: $maxRepositories, isFork: false, isArchived: false, isLocked: false, orderBy: {field: UPDATED_AT, direction: DESC})"`
			} `graphql:"repositoryOwner(login: $owner)"`
		}
		vars["maxRepositories"] = githubv4.Int(100)
		err := client.Query(ctx, &q, vars)
		if err != nil {
			return nil, err
		}
		repos = append(repos, q.RepositoryOwner.Repositories.Nodes...)
		for _, repo := range repos {
			fmt.Printf("Found repo: %v\n", repo.NameWithOwner)
		}
	} else {
		var q struct {
			Repository Repository `graphql:"repository(owner: $owner, name: $name)"`
		}
		vars["name"] = githubv4.String(repo)
		err := client.Query(ctx, &q, vars)
		if err != nil {
			return nil, err
		}
		repos = append(repos, q.Repository)
	}

	pullRequests := []*PullRequest{}
	for _, repo := range repos {
		for _, pr := range repo.PullRequests.Nodes {
			pullRequest := &PullRequest{
				RepositoryOwner: string(repo.Owner.Login),
				RepositoryName:  string(repo.Name),
				Number:          int(pr.Number),
				Title:           string(pr.Title),
				State:           string(pr.State),
				CreatedAt:       pr.CreatedAt.Time,
				ID:              pr.ID,
				StatusRollup:    string(pr.Commits.Nodes[0].Commit.StatusCheckRollup.State),
			}

			pullRequests = append(pullRequests, pullRequest)
		}
	}

	return pullRequests, nil
}

func ApprovePullRequest(ghClient *githubv4.Client, ctx context.Context, prId githubv4.ID, skip bool) {

	commitMessage := ctx.Value("message").(githubv4.String)
	event := githubv4.PullRequestReviewEventApprove

	input := githubv4.AddPullRequestReviewInput{
		PullRequestID: prId,
		Body:          &commitMessage,
		Event:         &event,
	}

	var m struct {
		AddPullRequestReview struct {
			PullRequest struct {
				Number     githubv4.Int
				Repository struct {
					NameWithOwner githubv4.String
				}
			}
			PullRequestReview struct {
				State githubv4.String
			}
		} `graphql:"addPullRequestReview(input: $input)"`
	}

	err := ghClient.Mutate(ctx, &m, input, nil)
	if err != nil && !skip {
		log.Printf("Could not approve pull request %v, did you try to approve your on pull request? - %v\n", prId, err)
	}

	review := m.AddPullRequestReview.PullRequestReview

	if err != nil && skip {
		fmt.Printf("Could not approve pull request, skipping.")
	} else {
		fmt.Printf("PR %v: %v\n", prId, review.State)
	}
}
func MergePullRequest(ghClient *githubv4.Client, ctx context.Context, prId githubv4.ID, mergeMethod *githubv4.PullRequestMergeMethod, skip bool) {

	commitMessage := githubv4.String(defaultCommitMsg())

	input := githubv4.MergePullRequestInput{
		PullRequestID:  prId,
		CommitHeadline: &commitMessage,
		MergeMethod:    mergeMethod,
	}

	var m struct {
		MergePullRequest struct {
			PullRequest struct {
				Merged     bool
				Number     githubv4.Int
				Repository struct {
					NameWithOwner githubv4.String
				}
			}
		} `graphql:"mergePullRequest(input: $input)"`
	}

	err := ghClient.Mutate(ctx, &m, input, nil)
	if err != nil {
		log.Printf("Could not merge PR %s, skipping: %v\n", prId, err)
	}

	pr := m.MergePullRequest.PullRequest

	fmt.Printf("PR %s#%d merged: %v\n", pr.Repository.NameWithOwner, pr.Number, pr.Merged)
}

func ClosePullRequest(ghClient *githubv4.Client, ctx context.Context, prId githubv4.ID, skip bool) {

	input := githubv4.ClosePullRequestInput{
		PullRequestID: prId,
	}

	var m struct {
		ClosePullRequest struct {
			PullRequest struct {
				Closed     bool
				Number     githubv4.Int
				Repository struct {
					NameWithOwner githubv4.String
				}
			}
		} `graphql:"closePullRequest(input: $input)"`
	}

	err := ghClient.Mutate(ctx, &m, input, nil)
	if err != nil {
		log.Printf("Could not close PR %s, skipping: %v\n", prId, err)
	}

	pr := m.ClosePullRequest.PullRequest

	fmt.Printf("PR %s#%d closed: %v\n", pr.Repository.NameWithOwner, pr.Number, pr.Closed)

}

func defaultCommitMsg() string {
	return "Merged by gomerge CLI."
}

func DefaultApproveMsg() string {
	return `PR has been approved by [GoMerge](https://github.com/Cian911/gomerge) tool. :rocket:`
}
