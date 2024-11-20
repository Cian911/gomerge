package gitclient

import (
	"testing"
)

func TestDefaultApproveMsg(t *testing.T) {
	t.Run("It returns a default approve message.", func(t *testing.T) {
		got := DefaultApproveMsg()
		want := "PR has been approved by [GoMerge](https://github.com/Cian911/gomerge) tool. :rocket:"

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
