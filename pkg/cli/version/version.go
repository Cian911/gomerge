package version

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/savioxavier/termlink"
	"github.com/spf13/cobra"
)

var (
	// Hold current application version
	Version string
	// Holds current application build number
	Build string
	// Holds current application build date
	BuildDate string
  // Sponsors
  Kofi string
  BMAC string
  Github string
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
  file, err := ioutil.ReadFile("./gomerge_ascii")
  if err != nil {
    log.Println("Could not read ascii art. Skipping.\n")
  }

  fmt.Fprintln(w, string(file))
	fmt.Fprintf(w, "Version: %s\nBuild: %s\nBuild Date: %s\n\n", Version, Build, BuildDate)
	fmt.Fprintln(w, "---\n")
  fmt.Fprintln(w, "Help support me! Any donations no matter how much is greatly appreciated.\n")
  fmt.Printf(
		`
	✔︎ %s

	✔︎ %s

	✔︎ %s
		`,
    termlink.ColorLink(Kofi, Kofi, "italic magenta"),
    termlink.ColorLink(BMAC, BMAC, "italic yellow"),
		termlink.ColorLink(Github, Github, "italic red"),
	)
  fmt.Println()
}
