package gitclient

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func Client(github_token string, ctx context.Context) (client *github.Client) {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: github_token,
		},
	)
	tokenContext := oauth2.NewClient(ctx, tokenSource)
	client = github.NewClient(tokenContext)

	return
}

func DefaultCommitMsg() string {
	return "Merged by gomerge CLI."
}
