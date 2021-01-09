package list

import "testing"

func TestParseOrgRepo(t *testing.T) {
	t.Run("It returns a valid tuple when no config is present", func(t *testing.T) {
		repo := "Cian911/syncwave"
		configPresent := false

		want1 := "Cian911"
		want2 := "syncwave"

		got1, got2 := parseOrgRepo(repo, configPresent)

		if got1 != want1 || got2 != want2 {
			t.Errorf("got1: %s, got2: %s, want1: %s, want2: %s", got1, got2, want1, want2)
		}
	})
}
