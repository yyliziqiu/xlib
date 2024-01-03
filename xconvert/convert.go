package xconvert

import (
	"strconv"
	"time"
)

func Int2String(i int) string {
	return strconv.Itoa(i)
}

func Int642String(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Float642String(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func String2Int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func String2Int64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func String2Float64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
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
