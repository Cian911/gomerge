package gitclient

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	MergeLabel   = "gomerge-merged"
	ApproveLabel = "gomerge-approved"
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

func ApprovePullRequest(ghClient *github.Client, ctx context.Context, org, repo string, prId int) {
	// Create review
	t := defaultApproveMsg(prId)
	e := "APPROVE"
	reviewRequest := &github.PullRequestReviewRequest{
		Body:  &t,
		Event: &e,
	}
	review, _, err := ghClient.PullRequests.CreateReview(ctx, org, repo, prId, reviewRequest)
	if err != nil {
		// TODO: Parse error to check if user tried to approve their own PR..
		log.Fatalf("Could not approve pull request, did you try to approve your on pull request? - %v", err)
	}

	AddLabel(ghClient, ctx, org, repo, prId, ApproveLabel)

	fmt.Printf("PR #%d: %v\n", prId, *review.State)
}

func MergePullRequest(ghClient *github.Client, ctx context.Context, org, repo string, prId int, mergeMethod string) {
	result, _, err := ghClient.PullRequests.Merge(ctx, org, repo, prId, defaultCommitMsg(), &github.PullRequestOptions{MergeMethod: mergeMethod})
	if err != nil {
		log.Fatal(err)
	}

	AddLabel(ghClient, ctx, org, repo, prId, MergeLabel)

	fmt.Println(fmt.Sprintf("PR #%d: %v.", prId, *result.Message))
}

func AddLabel(ghClient *github.Client, ctx context.Context, org, repo string, prId int, label string) {
	// Get label
	_, _, err := ghClient.Issues.GetLabel(ctx, org, repo, label)

	// If label does not exist, create it
	if err != nil {
		labelDesc := "Merged/Approved by Gomerge."
		_, _, err := ghClient.Issues.CreateLabel(ctx, org, repo, &github.Label{
			Name:        &label,
			Description: &labelDesc,
		})

		if err != nil {
			log.Fatal(err)
		}
	}

	// Add label
	_, _, err = ghClient.Issues.AddLabelsToIssue(ctx, org, repo, prId, []string{label})

	if err != nil {
		log.Fatal(err)
	}
}

func defaultCommitMsg() string {
	return "Merged by gomerge CLI."
}

func defaultApproveMsg(prId int) string {
	return fmt.Sprintf(`PR #%d has been approved by [GoMerge](https://github.com/Cian911/gomerge) tool. :rocket:`, prId)
}
