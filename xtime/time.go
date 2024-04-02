package xtime

import (
	"time"
)

func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

var daysOfMonth = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func DaysOfMonth(year int, month time.Month) int {
	if month < 1 || month > 12 {
		return 0
	}
	if IsLeapYear(year) && month == time.February {
		return 29
	}
	return daysOfMonth[month-1]
}

func DayBegin(t time.Time, loc *time.Location) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func DayEnd(t time.Time, loc *time.Location) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, loc)
}

func DayRange(t time.Time, loc *time.Location) (time.Time, time.Time) {
	return DayBegin(t, loc), DayEnd(t, loc)
}

func WeekBegin(t time.Time, loc *time.Location) time.Time {
	n := int(t.Weekday())
	if n == 0 {
		n = 7
	}
	return DayBegin(t.AddDate(0, 0, 1-n), loc)
}

func WeekEnd(t time.Time, loc *time.Location) time.Time {
	n := int(t.Weekday())
	if n == 0 {
		n = 7
	}
	return DayEnd(t.AddDate(0, 0, 7-n), loc)
}

func WeekRange(t time.Time, loc *time.Location) (time.Time, time.Time) {
	return WeekBegin(t, loc), WeekEnd(t, loc)
}

func MonthBegin(t time.Time, loc *time.Location) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, loc)
}

func MonthEnd(t time.Time, loc *time.Location) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, DaysOfMonth(year, month), 23, 59, 59, 0, loc)
}

func MonthRange(t time.Time, loc *time.Location) (time.Time, time.Time) {
	return MonthBegin(t, loc), MonthEnd(t, loc)
}

func YearBegin(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, loc)
}

func YearEnd(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 0, loc)
}

func YearRange(t time.Time, loc *time.Location) (time.Time, time.Time) {
	return YearBegin(t, loc), YearEnd(t, loc)
}
