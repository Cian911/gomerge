package list

import (
	"os"
	"testing"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/cian911/go-merge/pkg/gitclient"
	"github.com/cian911/go-merge/pkg/printer"
)

func TestInitTable(t *testing.T) {
	t.Run("It returns a tablewriter pointer", func(t *testing.T) {
		got := initTable()
		want := &tablewriter.Table{}

		assert.IsType(t, got, want)
	})
}

func TestFormatTable(t *testing.T) {
	t.Run("It returns a string array", func(t *testing.T) {
		number := 1
		state := "#open"
		title := "My Pr"
		createdAt := time.Now()

		pr := &gitclient.PullRequest{
			RepositoryOwner: "Cian911",
			RepositoryName:  "syncwave",
			Number:          number,
			State:           state,
			Title:           title,
			StatusRollup:    "SUCCESS",
			CreatedAt:       createdAt,
		}

		got, _ := formatTable(pr)
		want := []string{
			"#1",
			"#open ",
			"My Pr",
			"Cian911/syncwave",
			printer.FormatTime(&pr.CreatedAt),
		}

		assert.Equal(t, got, want)
	})
}

func TestListGetToken(t *testing.T) {
	var (
		flag   = "flag@token"
		config = "config@token"
		envVar = "env@token"
	)

	t.Run(
		"When a given token is set by flag, it should return token as the flag value",
		func(t *testing.T) {
			got, err := getToken(flag, "")
			want := flag
			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run(
		"When a given token is set by config, it should return token as defined on the configuration file",
		func(t *testing.T) {
			got, err := getToken("", config)
			want := config
			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run(
		"When a given token is set by environment variable, it should return token as defined on the environment",
		func(t *testing.T) {
			os.Setenv(TokenEnvVar, envVar)
			got, err := getToken("", "")
			want := envVar
			assert.Nil(t, err)
			assert.Equal(t, want, got)
			os.Unsetenv(TokenEnvVar)
		},
	)

	t.Run(
		"When a given token is set on flag and config file, it should return the value set on flag",
		func(t *testing.T) {
			got, err := getToken(flag, config)
			want := flag
			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run(
		"When a given token is set on flag and environment, it should return the value set on the flag",
		func(t *testing.T) {
			os.Setenv(TokenEnvVar, envVar)
			got, err := getToken(flag, "")
			want := flag
			assert.Nil(t, err)
			assert.Equal(t, want, got)
			os.Unsetenv(TokenEnvVar)
		},
	)

	t.Run(
		"When a given token is set on config file and environment, it should return the value set on the config file",
		func(t *testing.T) {
			os.Setenv(TokenEnvVar, envVar)
			got, err := getToken("", config)
			want := config
			assert.Nil(t, err)
			assert.Equal(t, want, got)
			os.Unsetenv(TokenEnvVar)
		},
	)

	t.Run(
		"When a given token is set on flag, config file, and environment, it should return the value set on flag",
		func(t *testing.T) {
			os.Setenv(TokenEnvVar, envVar)
			got, err := getToken(flag, config)
			want := flag
			assert.Nil(t, err)
			assert.Equal(t, want, got)
			os.Unsetenv(TokenEnvVar)
		},
	)

	t.Run("When no token is passed should return error", func(t *testing.T) {
		got, err := getToken("", "")
		assert.Equal(t, "", got)
		assert.Error(t, err)
	})
}

func TestMergeQueueAndMergeMethodMutuallyExclusive(t *testing.T) {
	t.Run("Should error when both merge-queue and merge-method are set", func(t *testing.T) {
		viper.Reset()
		viper.Set("merge-queue", true)
		viper.Set("merge-method", "squash")

		err := validateMergeFlags()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Cannot use --merge-queue with --merge-method")
	})

	t.Run("Should not error when only merge-queue is set", func(t *testing.T) {
		viper.Reset()
		viper.Set("merge-queue", true)
		viper.Set("merge-method", "")

		err := validateMergeFlags()
		assert.NoError(t, err)
	})

	t.Run("Should not error when only merge-method is set", func(t *testing.T) {
		viper.Reset()
		viper.Set("merge-queue", false)
		viper.Set("merge-method", "squash")

		err := validateMergeFlags()
		assert.NoError(t, err)
	})
}
