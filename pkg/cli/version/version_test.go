package version

import (
	"bytes"
	"testing"
)

func TestPrintVersionInformation(t *testing.T) {
	t.Run("It prints the version information", func(t *testing.T) {
		Version = "0.1"
		Build = "alpha"
		BuildDate = "2021-05-06T23:09:49Z"

		var b bytes.Buffer
		printVersionInformation(&b)

		got := b.String()
		want := `
Gomerge: 
version: 0.1
build: alpha
build date: 2021-05-06T23:09:49Z
`

		if got != want {
			t.Errorf("got %s, want: %s", got, want)
		}
	})
}
