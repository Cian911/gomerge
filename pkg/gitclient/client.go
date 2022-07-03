package gitclient

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func Client(githubToken string, ctx context.Context) (client *github.Client) {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: githubToken,
		},
	)
	tokenContext := oauth2.NewClient(ctx, tokenSource)
	client = github.NewClient(tokenContext)

	return
}

func ApprovePullRequest(ghClient *github.Client, ctx context.Context, org, repo string, prId int, skip bool) {
	// Create review
	t := defaultApproveMsg(prId)
	e := "APPROVE"
	reviewRequest := &github.PullRequestReviewRequest{
		Body:  &t,
		Event: &e,
	}
	review, _, err := ghClient.PullRequests.CreateReview(ctx, org, repo, prId, reviewRequest)
	if err != nil && !skip {
		log.Fatalf("Could not approve pull request, did you try to approve your on pull request? - %v", err)
	} else {
		log.Printf("Could not approve pull request, did you try to approve your on pull request? Skipping: %v \n", err)
	}

	fmt.Printf("PR #%d: %v\n", prId, *review.State)
}

func MergePullRequest(ghClient *github.Client, ctx context.Context, org, repo string, prId int, mergeMethod string, skip bool) {
	result, _, err := ghClient.PullRequests.Merge(ctx, org, repo, prId, defaultCommitMsg(), &github.PullRequestOptions{MergeMethod: mergeMethod})
	if err != nil && !skip {
		log.Fatal(err)
	} else {
		log.Printf("Could not merge PR #%d, skipping: %v\n", prId, err)
	}

	fmt.Println(fmt.Sprintf("PR #%d: %v.", prId, *result.Message))
}

func defaultCommitMsg() string {
	return "Merged by gomerge CLI."
}

func defaultApproveMsg(prId int) string {
	return fmt.Sprintf(`PR #%d has been approved by [GoMerge](https://github.com/Cian911/gomerge) tool. :rocket:`, prId)
}
