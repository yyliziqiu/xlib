package xconv

func I2S(i int) string {
	return IntToString(i)
}

func S2I(s string) int {
	return StringToInt(s)
}

func I642S(i int64) string {
	return Int64ToString(i)
}

func S2I64(s string) int64 {
	return StringToInt64(s)
}

func F642S(f float64, prec int) string {
	return Float64ToString(f, prec)
}

func S2F64(s string) float64 {
	return StringToFloat64(s)
}

func B2S(b bool) string {
	return BoolToString(b)
}

func S2B(s string) bool {
	return StringToBool(s)
}

func T2S(t int64) string {
	return TimestampToString(t)
}

func S2T(s string) int64 {
	return StringToTimestamp(s)
}
