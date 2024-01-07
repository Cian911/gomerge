package list

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func List() (c *cobra.Command) {
  c = &cobra.Command{
    Use: "list",
    Short: "List all pull requests for a repo you want to action.",
    Run: func(cmd *cobra.Command, args []string) {
      if viper.ConfigFileUsed() != "" {
        // Use the config
      } else {
        // Use the CLI args
      }
    },
  }
  
  return
}

func initConfig() {}

func validateFlags() {}

func validateConfig() {}
