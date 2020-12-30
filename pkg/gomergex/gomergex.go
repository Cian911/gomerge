package gomergex

import "context"

type GoMergeContext struct {
	GithubClient *gitclient.GitClient
}

func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, &GoMergeContext{})
}
