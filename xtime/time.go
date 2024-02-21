package xtime

import (
	"time"
)

func GetTimestampOfHour(ts int64) int64 {
	t := time.Unix(ts, 0)
	year, month, day := t.Date()
	return time.Date(year, month, day, t.Hour(), 0, 0, 0, time.Local).Unix()
}

func GetTimestampOfDay(ts int64) int64 {
	t := time.Unix(ts, 0)
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
}

func GetTimestampOfMonth(ts int64) int64 {
	t := time.Unix(ts, 0)
	year, month, _ := t.Date()
	return time.Date(year, month, 0, 0, 0, 0, 0, time.Local).Unix()
}

func GetTimestampOfYear(ts int64) int64 {
	t := time.Unix(ts, 0)
	return time.Date(t.Year(), 0, 0, 0, 0, 0, 0, time.Local).Unix()
}

func TodayRange() (time.Time, time.Time) {
	year, month, day := time.Now().Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	stops := time.Now()
	return start, stops
}

func TodayStartAt() time.Time {
	start, _ := TodayRange()
	return start
}

func YesterdayRange() (time.Time, time.Time) {
	year, month, day := time.Now().AddDate(0, 0, -1).Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	stops := time.Date(year, month, day, 23, 59, 59, 0, time.Local)
	return start, stops
}

func YesterdayStartAt() time.Time {
	start, _ := YesterdayRange()
	return start
}

func CurrentMonthRange() (time.Time, time.Time) {
	year, month, _ := time.Now().Date()
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	stops := time.Now()
	return start, stops
}

func CurrentMonthStartAt() time.Time {
	start, _ := CurrentMonthRange()
	return start
}
