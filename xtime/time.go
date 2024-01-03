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
