package xconv

import (
	"strconv"
	"time"
)

func I2S(i int) string {
	return strconv.Itoa(i)
}

func S2I(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func I642S(i int64) string {
	return strconv.FormatInt(i, 10)
}

func S2I64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func F642S(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func S2F64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func B2S(v bool) string {
	return strconv.FormatBool(v)
}

func S2B(v string) bool {
	b, _ := strconv.ParseBool(v)
	return b
}

func T2S(ts int64) string {
	return time.Unix(ts, 0).Format(time.DateTime)
}

func S2T(str string) int64 {
	t, err := time.Parse(time.DateTime, str)
	if err != nil {
		return 0
	}
	return t.Unix()
}
