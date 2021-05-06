package version

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Hold current application version
	Version string
	// Holds current application build number
	Build string
	// Holds current application build date
	BuildDate string
)

func NewCommand() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "version",
		Short: "Prints the current version and build information.",
		Run: func(cmd *cobra.Command, args []string) {
			printVersionInformation(os.Stdout)
		},
	}

	return c
}

func printVersionInformation(w io.Writer) {
	fmt.Fprintf(w, "\nGomerge: \nversion: %s\nbuild: %s\nbuild date: %s\n", Version, Build, BuildDate)
}
