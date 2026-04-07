package transform

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// supportedFormats lists date formats that ParseFlexible will try.
// BUG: missing ISO 8601 with timezone offset.
var supportedFormats = []string{
	"2006-01-02",
	"01/02/2006",
	"02-Jan-2006",
	"January 2, 2006",
	"Jan 2, 2006",
	"2006-01-02T15:04:05Z",
	"2006-01-02 15:04:05",
	"Mon, 02 Jan 2006 15:04:05",
}

// ParseFlexible tries multiple date formats and returns the first successful parse.
func ParseFlexible(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}

	for _, layout := range supportedFormats {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %q", s)
}

// FormatISO formats a time.Time as an ISO 8601 string.
func FormatISO(t time.Time) string {
	return t.Format(time.RFC3339)
}

// RelativeTime returns a human-readable relative time string.
func RelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	future := diff < 0
	if future {
		diff = -diff
	}

	seconds := int(math.Abs(diff.Seconds()))
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	months := days / 30
	years := days / 365

	var unit string
	var count int

	switch {
	case seconds < 60:
		return "just now"
	case minutes < 60:
		unit = "minute"
		count = minutes
	case hours < 24:
		unit = "hour"
		count = hours
	case days < 30:
		unit = "day"
		count = days
	case months < 12:
		unit = "month"
		count = months
	default:
		unit = "year"
		count = years
	}

	if count != 1 {
		unit += "s"
	}

	if future {
		return fmt.Sprintf("in %d %s", count, unit)
	}
	return fmt.Sprintf("%d %s ago", count, unit)
}
