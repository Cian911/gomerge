package printer

import (
	"fmt"
	"time"
)

func FormatID(id *int) string {
	return fmt.Sprintf("%d", *id)
}

func FormatString(str *string) string {
	return fmt.Sprintf("%s", *str)
}

func FormatTime(t *time.Time) string {
	return t.Format(time.RFC3339)
}
