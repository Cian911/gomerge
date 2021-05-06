package main

import (
	"context"

	"github.com/cian911/go-merge/pkg/cli/gomerge"
	"github.com/cian911/go-merge/pkg/cli/version"
)

var (
	Version   string
	Build     string
	BuildDate string
)

func main() {
	version.Version = Version
	version.Build = Build
	version.BuildDate = BuildDate

	ctx := context.Background()
	gomerge.New().ExecuteContext(ctx)
}
