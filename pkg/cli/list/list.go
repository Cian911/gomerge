package list

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cian911/go-merge/internal/tui"
)

func List() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "list",
		Short: "List all pull requests for a repo you want to action.",
		Run: func(cmd *cobra.Command, args []string) {
			if viper.ConfigFileUsed() != "" {
				// Use the config
			} else {
				f, err := tea.LogToFile("debug.log", "debug")
				if err != nil {
					log.Fatalf("%v", err)
				}
				defer f.Close()

				// Use the CLI args
				model, err := tui.New()
				model.Init()

				if err != nil {
					log.Fatal(err)
				}

				p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
				if err := p.Start(); err != nil {
					log.Fatalf("Could not start tui: %v", err)
				}
			}
		},
	}

	return
}

func initConfig() {}

func validateFlags() {}

func validateConfig() {}
