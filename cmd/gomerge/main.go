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
  Kofi string
  BMAC string
  Github string
)

func main() {
	version.Version = Version
	version.Build = Build
	version.BuildDate = BuildDate
  version.Kofi = Kofi
  version.BMAC = BMAC
  version.Github = Github

	ctx := context.Background()
	gomerge.New().ExecuteContext(ctx)
}
