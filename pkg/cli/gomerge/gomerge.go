package gomerge

import (
	"github.com/cian911/go-merge/pkg/cli/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "gomerge",
		Short: "gomerge makes it simple to merge an open pull request from your terminal",
	}

	c.PersistentFlags().StringP("repo", "r", "", "Pass name of repository as argument (organization/repo).")
	c.PersistentFlags().StringP("token", "t", "", "Pass your github personal access token (PAT).")

	c.MarkFlagRequired("repo")
	c.MarkFlagRequired("token")

	viper.BindPFlag("repo", c.PersistentFlags().Lookup("repo"))
	viper.BindPFlag("token", c.PersistentFlags().Lookup("token"))

	c.AddCommand(list.NewCommand())

	return
}
