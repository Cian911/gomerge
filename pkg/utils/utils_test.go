package utils

import (
	"testing"

	"github.com/spf13/viper"
)

func TestReadConfigFile(t *testing.T) {
	t.Run("It successfully reads a config file", func(t *testing.T) {
		path := "./config_test.yaml"

		ReadConfigFile(path)

		got := viper.Get("organization")
		want := "Cian911"

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		got1 := viper.Get("token")
		want1 := "1234test@gh*token"

		if got1 != want1 {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("It extracts filename and ext from path", func(t *testing.T) {
		path := "./config_test.yaml"

		gotFilename, gotExt := parseConfigFile(path)
		wantFilename := "config_test"
		wantExt := "yaml"

		if gotFilename != wantFilename || gotExt != wantExt {
			t.Errorf(
				"got: %s, want: %s, got: %s, want: %s",
				gotFilename,
				wantFilename,
				gotExt,
				wantExt,
			)
		}
	})
}
