package gitclient

import (
	"context"
	"fmt"
	"log"
	"time"

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
	NeedsReview     bool
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
		// NB: this only works for classic branch protection rules. Repository
		// rulesets don't appear to be visible in the API
		BaseRef struct {
			BranchProtectionRule struct {
				RequiredApprovingReviewCount githubv4.Int
			}
		}
		Reviews struct {
			Nodes []struct {
				AuthorCanPushToRepository githubv4.Boolean
			}
		} `graphql:"reviews(first: 100, states: [APPROVED])"`
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
			reviews := 0
			for _, review := range pr.Reviews.Nodes {
				if review.AuthorCanPushToRepository {
					reviews++
				}
			}
			pullRequest := &PullRequest{
				RepositoryOwner: string(repo.Owner.Login),
				RepositoryName:  string(repo.Name),
				Number:          int(pr.Number),
				Title:           string(pr.Title),
				State:           string(pr.State),
				CreatedAt:       pr.CreatedAt.Time,
				ID:              pr.ID,
				StatusRollup:    string(pr.Commits.Nodes[0].Commit.StatusCheckRollup.State),
				NeedsReview:     reviews < int(pr.BaseRef.BranchProtectionRule.RequiredApprovingReviewCount),
			}

			pullRequests = append(pullRequests, pullRequest)
		}
	}

	return pullRequests, nil
}

func ApprovePullRequest(ghClient *githubv4.Client, ctx context.Context, pr *PullRequest, skip bool) {

	commitMessage := githubv4.String(DefaultApproveMsg())
	event := githubv4.PullRequestReviewEventApprove

	input := githubv4.AddPullRequestReviewInput{
		PullRequestID: pr.ID,
		Body:          &commitMessage,
		Event:         &event,
	}

	var m struct {
		AddPullRequestReview struct {
			PullRequestReview struct {
				State githubv4.PullRequestReviewState
			}
		} `graphql:"addPullRequestReview(input: $input)"`
	}

	err := ghClient.Mutate(ctx, &m, input, nil)
	if err != nil && !skip {
		log.Printf("Could not approve pull request %s/%s#%d - %v\n", pr.RepositoryOwner, pr.RepositoryName, pr.Number, err)
	}

	review := m.AddPullRequestReview.PullRequestReview

	if err != nil && skip {
		fmt.Printf("Could not approve pull request, skipping.")
	} else {
		fmt.Printf("%s/%s#%d: %v\n", pr.RepositoryOwner, pr.RepositoryName, pr.Number, review.State)
	}
}

func MergePullRequest(ghClient *githubv4.Client, ctx context.Context, pr *PullRequest, mergeMethod *githubv4.PullRequestMergeMethod, skip bool) {

	input := githubv4.MergePullRequestInput{
		PullRequestID: pr.ID,
		MergeMethod:   mergeMethod,
	}

	var m struct {
		MergePullRequest struct {
			PullRequest struct {
				Merged     bool
				State      githubv4.PullRequestState
				Number     githubv4.Int
				Repository struct {
					NameWithOwner githubv4.String
				}
			}
		} `graphql:"mergePullRequest(input: $input)"`
	}

	err := ghClient.Mutate(ctx, &m, input, nil)
	if err != nil {
		log.Printf("Could not merge %s/%s#%d - %v\n", pr.RepositoryOwner, pr.RepositoryName, pr.Number, err)
	}

	prOut := m.MergePullRequest.PullRequest
	fmt.Printf("%s/%s#%d merged: %v\n", pr.RepositoryOwner, pr.RepositoryName, pr.Number, prOut.State)

}

func ClosePullRequest(ghClient *githubv4.Client, ctx context.Context, pr *PullRequest, skip bool) {

	input := githubv4.ClosePullRequestInput{
		PullRequestID: pr.ID,
	}

	var m struct {
		ClosePullRequest struct {
			PullRequest struct {
				Closed     bool
				State      githubv4.PullRequestState
				Number     githubv4.Int
				Repository struct {
					NameWithOwner githubv4.String
				}
			}
		} `graphql:"closePullRequest(input: $input)"`
	}

	err := ghClient.Mutate(ctx, &m, input, nil)
	if err != nil {
		log.Printf("Could not close %s/%s#%d - %v\n", pr.RepositoryOwner, pr.RepositoryName, pr.Number, err)
	}

	prOut := m.ClosePullRequest.PullRequest

	fmt.Printf("%s/%s#%d merged: %v\n", pr.RepositoryOwner, pr.RepositoryName, pr.Number, prOut.State)

}

func DefaultApproveMsg() string {
	return `PR has been approved by [GoMerge](https://github.com/Cian911/gomerge) tool. :rocket:`
}
