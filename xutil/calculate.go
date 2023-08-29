package xutil

func DivFloat64Int(a float64, b int) float64 {
	if b == 0 {
		return 0
	}
	return a / float64(b)
}

func DivFloat64Int64(a float64, b int64) float64 {
	if b == 0 {
		return 0
	}
	return a / float64(b)
}

func DivFloat64Float64(a float64, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func RateInt(a int, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a-b) / float64(b)
}

func RateInt64(a int64, b int64) float64 {
	if b == 0 {
		return 0
	}
	return float64(a-b) / float64(b)
}

func RateFloat64(a float64, b float64) float64 {
	if b == 0 {
		return 0
	}
	return (a - b) / b
}

// RateIntInt deprecated
func RateIntInt(a int, b int) float64 {
	return RateInt(a, b)
}

// RateInt64Int64 deprecated
func RateInt64Int64(a int64, b int64) float64 {
	return RateInt64(a, b)
}

// RateFloat64Float64 deprecated
func RateFloat64Float64(a float64, b float64) float64 {
	return RateFloat64(a, b)
}
