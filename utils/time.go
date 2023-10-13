package utils

import (
	"time"
)

func DayDate(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
}

func ISOWeekDate(v time.Time) time.Time {
	weekday := int(v.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	return DayDate(v).AddDate(0, 0, -1*(weekday-1))
}

func MonthDate(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), 1, 0, 0, 0, 0, v.Location())
}

func YearDate(v time.Time) time.Time {
	return time.Date(v.Year(), 1, 1, 0, 0, 0, 0, v.Location())
}
