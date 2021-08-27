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

func FormatMergeability(str *string) string {
	switch *str {
	case "behind", "blocked", "dirty", "draft", "unknown", "unstable":
		return "FALSE"
	case "clean", "has_hooks":
		return "TRUE"
	default:
		return "UNKNOWN"
	}
}
