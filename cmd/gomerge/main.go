package main

import (
	"context"

	"github.com/cian911/go-merge/pkg/cli/gomerge"
)

func main() {
	ctx := context.Background()
	gomerge.New().ExecuteContext(ctx)
}
