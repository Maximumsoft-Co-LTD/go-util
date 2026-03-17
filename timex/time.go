// Package timex provides time truncation and formatting helpers.
package timex

import (
	"fmt"
	"time"

	"github.com/nleeper/goment"
)

// LayoutISO is the standard ISO datetime format for time.Format / time.Parse.
const LayoutISO = "2006-01-02 15:04:05"

// StartOfDay returns the start of the current day in UTC using goment.
func StartOfDay() (time.Time, error) {
	goDate, err := goment.New()
	if err != nil {
		return time.Time{}, fmt.Errorf("creating goment instance: %w", err)
	}
	return goDate.StartOf("day").UTC().ToTime(), nil
}

// EndOfDay returns the end of the current day in UTC using goment.
func EndOfDay() (time.Time, error) {
	goDate, err := goment.New()
	if err != nil {
		return time.Time{}, fmt.Errorf("creating goment instance: %w", err)
	}
	return goDate.EndOf("day").UTC().ToTime(), nil
}

// TruncateDate returns 00:00 of the specified day-of-month relative to the given time.
// Falls back to the previous month if the date has not yet occurred.
// Returns zero time.Time if date is out of range (0–31).
func TruncateDate(now time.Time, date int) time.Time {
	if date < 0 || date > 31 {
		return time.Time{}
	}
	rounded := time.Date(now.Year(), now.Month(), date, 0, 0, 0, 0, time.Local)
	if rounded.After(now) {
		lastDateOfMonth := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.UTC).Day()
		if date > lastDateOfMonth {
			date = lastDateOfMonth
		}
		rounded = time.Date(now.Year(), now.Month()-1, date, 0, 0, 0, 0, time.Local)
	}
	return rounded
}

// TruncateWeek returns 00:00 of the most recent occurrence of the specified weekday
// relative to the given time. day: 0 = Sunday, 6 = Saturday.
// Returns zero time.Time if day is out of range.
func TruncateWeek(now time.Time, day int) time.Time {
	if day < 0 || day > 6 {
		return time.Time{}
	}
	daysToSubtract := (int(now.Weekday()) - day + 7) % 7
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -daysToSubtract)
}

// FormatDateTimeString formats a time.Time to LayoutISO string.
func FormatDateTimeString(t time.Time) (string, error) {
	dateNow, err := goment.New(t.Local().String(), "YYYY-MM-DD HH:mm:ss")
	if err != nil {
		return "", fmt.Errorf("parsing time with goment: %w", err)
	}
	return dateNow.ToTime().Format(LayoutISO), nil
}

// Bod returns the beginning of day in UTC (17:00 UTC = 00:00 ICT).
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 17, 0, 0, 0, t.UTC().Location())
}

// BodLocal returns the beginning of day at 00:00 in local timezone.
func BodLocal(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// Eod returns the exclusive end-of-day boundary in UTC.
// Returns 17:00 UTC of the next day (= 00:00 ICT of the next day).
// Use with half-open interval: [Bod, Eod)
func Eod(t time.Time) time.Time {
	year, month, day := t.Date()
	nextDay := time.Date(year, month, day, 1, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	return Bod(nextDay)
}

// EodLocal returns the exclusive end-of-day boundary in local timezone.
// Returns 00:00 of the next day. Use with half-open interval: [BodLocal, EodLocal)
func EodLocal(t time.Time) time.Time {
	year, month, day := t.Date()
	nextDay := time.Date(year, month, day, 1, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	return BodLocal(nextDay)
}

// TimeFormat returns the given time as a concatenated string at the specified granularity.
// Accepts "day", "hour", "minute", "second". Returns "" for unknown filters.
func TimeFormat(t time.Time, filter string) string {
	y, m, d := t.Date()
	switch filter {
	case "day":
		return fmt.Sprint(y) + fmt.Sprint(m) + fmt.Sprint(d)
	case "hour":
		return fmt.Sprint(y) + fmt.Sprint(m) + fmt.Sprint(d) + fmt.Sprint(t.Hour())
	case "minute":
		return fmt.Sprint(y) + fmt.Sprint(m) + fmt.Sprint(d) + fmt.Sprint(t.Hour()) + fmt.Sprint(t.Minute())
	case "second":
		return fmt.Sprint(y) + fmt.Sprint(m) + fmt.Sprint(d) + fmt.Sprint(t.Hour()) + fmt.Sprint(t.Minute()) + fmt.Sprint(t.Second())
	default:
		return ""
	}
}

// TimeNowOnly extracts the hour/minute/second from a Unix timestamp
// and applies them to the reference time's date. Returns as Unix timestamp.
func TimeNowOnly(ref time.Time, t int64) int64 {
	dt := time.Unix(t, 0)
	result := time.Date(ref.Year(), ref.Month(), ref.Day(), dt.Hour(), dt.Minute(), dt.Second(), 0, time.Local)
	return result.Unix()
}

// TimeInRange checks if tTime falls within the half-open interval [startTime, endTime).
func TimeInRange(tTime, startTime, endTime time.Time) bool {
	return !tTime.Before(startTime) && tTime.Before(endTime)
}
