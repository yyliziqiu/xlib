package xutil

import (
	"strconv"
	"time"
)

func NewTimer() Timer {
	return Timer{
		start: time.Now(),
		pause: time.Now(),
	}
}

type Timer struct {
	start time.Time
	pause time.Time
}

func (t Timer) Pause() time.Duration {
	d := time.Now().Sub(t.pause)
	t.pause = time.Now()
	return d
}

func (t Timer) Pauses() string {
	return t.manualDuration(t.Pause())
}

var timeUnit = []string{"ns", "us", "ms", "s"}

func (t Timer) manualDuration(du time.Duration) string {
	d := float64(du)

	i := 0
	for d > 1000 && i < len(timeUnit)-1 {
		d = d / 1000
		i++
	}

	return strconv.FormatFloat(d, 'f', 2, 64) + timeUnit[i]
}

func (t Timer) Stop() time.Duration {
	return time.Now().Sub(t.start)
}

func (t Timer) Stops() string {
	return t.manualDuration(t.Stop())
}

func Timestamp2String(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func String2Timestamp(timeString string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeString)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func GetStartTimestampOfDay(ts int64) int64 {
	t := time.Unix(ts, 0)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
}
