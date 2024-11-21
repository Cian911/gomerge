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

type commits struct {
	Nodes []struct {
		Commit struct {
			StatusCheckRollup struct {
				State githubv4.String
			}
		}
	}
}

type pullRequests struct {
	Nodes []struct {
		Number    githubv4.Int
		Title     githubv4.String
		State     githubv4.String
		Url       githubv4.URI
		CreatedAt githubv4.DateTime
		ID        githubv4.ID
		Commits   commits `graphql:"commits(last:1)"`
	}
}

// NB: githubv4 uses struct tags to define the GraphQL query. To omit the labels
// argument to pullRequests(), we need to define a compeltely new type. Luckily
// these are convertible to each other as of go 1.8.
type repository struct {
	NameWithOwner githubv4.String
	Owner         struct {
		Login githubv4.String
	}
	Name         githubv4.String
	PullRequests pullRequests `graphql:"pullRequests(states: [OPEN], first: $maxPullRequests, orderBy: {field: CREATED_AT, direction: DESC})"`
}

type repositoryWithPRLabels struct {
	NameWithOwner githubv4.String
	Owner         struct {
		Login githubv4.String
	}
	Name         githubv4.String
	PullRequests pullRequests `graphql:"pullRequests(states: [OPEN], labels: $labels, first: $maxPullRequests, orderBy: {field: CREATED_AT, direction: DESC})"`
}

func GetPullRequests(client *githubv4.Client, ctx context.Context, owner string, repo string, labels *[]githubv4.String) ([]*PullRequest, error) {

	vars := map[string]interface{}{
		"maxPullRequests": githubv4.Int(100),
		"owner":           githubv4.String(owner),
	}
	if len(*labels) > 0 {
		vars["labels"] = labels
	}

	repos := []repository{}

	if repo == "" {
		vars["maxRepositories"] = githubv4.Int(100)
		if len(*labels) > 0 {
			var q struct {
				RepositoryOwner struct {
					Repositories struct {
						Nodes []repositoryWithPRLabels
					} `graphql:"repositories(first: $maxRepositories, isFork: false, isArchived: false, isLocked: false, orderBy: {field: UPDATED_AT, direction: DESC})"`
				} `graphql:"repositoryOwner(login: $owner)"`
			}
			err := client.Query(ctx, &q, vars)
			if err != nil {
				return nil, err
			}
			for _, repo := range q.RepositoryOwner.Repositories.Nodes {
				repos = append(repos, repository(repo))
			}
		} else {
			var q struct {
				RepositoryOwner struct {
					Repositories struct {
						Nodes []repository
					} `graphql:"repositories(first: $maxRepositories, isFork: false, isArchived: false, isLocked: false, orderBy: {field: UPDATED_AT, direction: DESC})"`
				} `graphql:"repositoryOwner(login: $owner)"`
			}
			err := client.Query(ctx, &q, vars)
			if err != nil {
				return nil, err
			}
			for _, repo := range q.RepositoryOwner.Repositories.Nodes {
				repos = append(repos, repository(repo))
			}
		}
	} else {
		vars["name"] = githubv4.String(repo)
		if len(*labels) > 0 {
			var q struct {
				Repository repositoryWithPRLabels `graphql:"repository(owner: $owner, name: $name)"`
			}
			err := client.Query(ctx, &q, vars)
			if err != nil {
				return nil, err
			}
			repos = append(repos, repository(q.Repository))
		} else {
			var q struct {
				Repository repository `graphql:"repository(owner: $owner, name: $name)"`
			}
			err := client.Query(ctx, &q, vars)
			if err != nil {
				return nil, err
			}
			repos = append(repos, q.Repository)
		}
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

	input := githubv4.MergePullRequestInput{
		PullRequestID: prId,
		MergeMethod:   mergeMethod,
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

func DefaultApproveMsg() string {
	return `PR has been approved by [GoMerge](https://github.com/Cian911/gomerge) tool. :rocket:`
}
