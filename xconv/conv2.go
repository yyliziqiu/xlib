package xconv

import (
	"strconv"
	"time"
)

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func Float64ToString(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func StringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

func StringToBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func TimestampToString(t int64) string {
	return time.Unix(t, 0).Format(time.DateTime)
}

func StringToTimestamp(s string) int64 {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return 0
	}
	return t.Unix()
}
