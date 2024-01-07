package gitclient

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v56/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

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

func ApprovePullRequest(ghClient *github.Client, ctx context.Context, org, repo string, prId int, skip bool) {
	// Create review
	commitMsg := ctx.Value("message").(string)
	e := "APPROVE"
	reviewRequest := &github.PullRequestReviewRequest{
		Body:  &commitMsg,
		Event: &e,
	}
	review, _, err := ghClient.PullRequests.CreateReview(ctx, org, repo, prId, reviewRequest)
	if err != nil && !skip {
		log.Fatalf("Could not approve pull request, did you try to approve your on pull request? - %v", err)
	}

	if err != nil && skip {
		fmt.Printf("Could not approve pull request, skipping.")
	} else {
		fmt.Printf("PR #%d: %v\n", prId, *review.State)
	}
}

func MergePullRequest(ghClient *github.Client, ctx context.Context, org, repo string, prId int, mergeMethod string, skip bool) {
	shouldSkip, err := shouldSkip(ghClient, ctx, org, repo, prId)
	if err != nil {
		log.Printf("Could not get PR #%d, skipping: %v\n", prId, err)
		return
	}
	if shouldSkip && skip {
		fmt.Printf("PR #%d: %v\n", prId, "PR is not mergable, skipping.")
		return
	}
	result, _, err := ghClient.PullRequests.Merge(ctx, org, repo, prId, defaultCommitMsg(), &github.PullRequestOptions{MergeMethod: mergeMethod})
	if err != nil {
		log.Printf("Could not merge PR #%d, skipping: %v\n", prId, err)

		return
	}

	fmt.Sprintf("PR #%d: %v.\n", prId, *result.Message)
}

func shouldSkip(ghClient *github.Client, ctx context.Context, org string, repo string, prId int) (bool, error) {
	resultz, _, err := ghClient.PullRequests.Get(ctx, org, repo, prId)
	if err != nil {
		return false, fmt.Errorf("failed to check if pr should be skipped %w", err)
	}

	unstable := *resultz.MergeableState == "unstable" || !*resultz.Mergeable
	return unstable, err
}

func ClosePullRequest(ghClient *github.Client, ctx context.Context, org, repo string, prId int, prRef *github.PullRequest) {
	// Set Closed state for PR
	*prRef.State = "closed"
	result, _, err := ghClient.PullRequests.Edit(ctx, org, repo, prId, prRef)
	if err != nil {
		log.Printf("Could not close PR #%d - %v", prId, err)
	} else {
		fmt.Sprintf("PR #%d: %v.\n", prId, *result.State)
	}
}

func defaultCommitMsg() string {
	return "Merged by gomerge CLI."
}

func DefaultApproveMsg() string {
	return `PR has been approved by [GoMerge](https://github.com/Cian911/gomerge) tool. :rocket:`
}
