package xconv

import (
	"strconv"
	"time"
)

func I2S(i int) string {
	return strconv.Itoa(i)
}

func I642S(i int64) string {
	return strconv.FormatInt(i, 10)
}

func F642S(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func S2I(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func S2I64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func S2F64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func T2S(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func S2T(timeString string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeString)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
