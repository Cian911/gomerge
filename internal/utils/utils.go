package utils

import (
	"math"
	"strconv"
	"time"
)

const (
	ApproxDaysInYear  = 365
	ApproxDaysInMonth = 28
	DaysInWeek        = 7
)

func TimeElapsed(then time.Time) string {
	var parts []string
	var text string

	now := time.Now()
	diff := now.Sub(then)
	day := math.Floor(diff.Hours() / 24)
	year := math.Floor(day / ApproxDaysInYear)
	month := math.Floor(day / ApproxDaysInMonth)
	week := math.Floor(day / DaysInWeek)
	hour := math.Floor(math.Abs(diff.Hours()))
	minute := math.Floor(math.Abs(diff.Minutes()))
	second := math.Floor(math.Abs(diff.Seconds()))

	if year > 0 {
		parts = append(parts, strconv.Itoa(int(year))+"y")
	}

	if month > 0 {
		parts = append(parts, strconv.Itoa(int(month))+"mo")
	}

	if week > 0 {
		parts = append(parts, strconv.Itoa(int(week))+"w")
	}

	if day > 0 {
		parts = append(parts, strconv.Itoa(int(day))+"d")
	}

	if hour > 0 {
		parts = append(parts, strconv.Itoa(int(hour))+"h")
	}

	if minute > 0 {
		parts = append(parts, strconv.Itoa(int(minute))+"m")
	}

	if second > 0 {
		parts = append(parts, strconv.Itoa(int(second))+"s")
	}

	if now.After(then) {
		text = " ago"
	} else {
		text = " after"
	}

	if len(parts) == 0 {
		return "just now"
	}

	return parts[0] + text
}
