package gitclient

import (
	"testing"
)

func TestDefaultCommitMsg(t *testing.T) {
	t.Run("It returns a default commit message", func(t *testing.T) {
		got := DefaultCommitMsg()
		want := "Merged by gomerge CLI."

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
