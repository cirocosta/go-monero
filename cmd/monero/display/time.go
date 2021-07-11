package display

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// HoursInMonth represents the approximate number of hours in a month.
	//
	// Note that hours per month is approximate due to fluctuations in
	// month length.
	//
	HoursInMonth = 730

	// HoursInDay represents the number of hours in a day.
	//
	HoursInDay = 24

	// HoursInWeek represents the number of hours in a week.
	//
	HoursInWeek = HoursInDay * 7

	// HoursInYear represents the number of hours in a year.
	//
	HoursInYear = HoursInDay * 365
)

// Duration returns a string with a plain English description of the length of
// time that the time.Duration t contains.
//
func Duration(t time.Duration) string {
	hours := int(t.Hours())
	minutes := int(t.Minutes())
	seconds := int(t.Seconds())

	if hours >= HoursInYear {
		years := hours / HoursInYear

		if years == 1 {
			return "a year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}

	if hours >= HoursInMonth {
		months := hours / HoursInMonth

		if months == 1 {
			return "about a month ago"
		}
		return fmt.Sprintf("about %d months ago", months)
	}

	if hours >= HoursInWeek {
		weeks := hours / HoursInWeek

		if weeks == 1 {
			return "a week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	}

	if hours >= HoursInDay {
		days := hours / HoursInDay

		if days == 1 {
			return "a day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}

	if hours > 0 {
		if hours == 1 {
			return "an hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	if minutes > 0 {
		if minutes == 1 {
			return "a minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}

	if seconds > 2 {
		return strconv.Itoa(seconds) + " seconds ago"
	}

	return "just now"
}

// Since provides an easy way to call Duration without having to call
// time.Since on a time.Time first.
func Since(t time.Time) string {
	return Duration(time.Since(t))
}
