package gomerge

import (
	"github.com/cian911/go-merge/pkg/cli/list"
	"github.com/cian911/go-merge/pkg/cli/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "gomerge",
		Short: "Gomerge makes it simple to merge an open pull request from your terminal.",
	}

	c.PersistentFlags().StringP("repo", "r", "", "Pass name of repository as argument (organization/repo).")
	c.PersistentFlags().StringP("token", "t", "", "Pass your github personal access token (PAT).")
	c.PersistentFlags().StringP("config", "c", "", "Pass an optional config file as an argument with list of repositories.")
	c.PersistentFlags().BoolP("approve", "a", false, "Pass an optional approve flag as an argument which will only approve and not merge selected repos.")
	c.PersistentFlags().StringP("merge-method", "m", "", "Pass an optional merge method for the pull request (merge [default], squash, rebase).")
	c.PersistentFlags().BoolP("skip", "s", false, "Pass an optional flag to skip a pull request and continue if one or more are not mergable.")
	c.PersistentFlags().BoolP("close", "", false, "Pass an optional argument to close a pull request.")

	c.MarkFlagRequired("token")

	viper.BindPFlag("repo", c.PersistentFlags().Lookup("repo"))
	viper.BindPFlag("token", c.PersistentFlags().Lookup("token"))
	viper.BindPFlag("config", c.PersistentFlags().Lookup("config"))
	viper.BindPFlag("approve", c.PersistentFlags().Lookup("approve"))
	viper.BindPFlag("merge-method", c.PersistentFlags().Lookup("merge-method"))
	viper.BindPFlag("skip", c.PersistentFlags().Lookup("skip"))
	viper.BindPFlag("close", c.PersistentFlags().Lookup("close"))

	c.AddCommand(list.NewCommand())
	c.AddCommand(version.NewCommand())

	return
}
