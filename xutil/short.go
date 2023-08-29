package xutil

func If(exp bool, a, b interface{}) interface{} {
	if exp {
		return a
	}
	return b
}

func IfInt(exp bool, a, b int) int {
	if exp {
		return a
	}
	return b
}

func IfInt64(exp bool, a, b int64) int64 {
	if exp {
		return a
	}
	return b
}

func IfFloat64(exp bool, a, b float64) float64 {
	if exp {
		return a
	}
	return b
}

func IfString(exp bool, a, b string) string {
	if exp {
		return a
	}
	return b
}

func IfNil(a, b interface{}) interface{} {
	if a == nil {
		return b
	}
	return a
}

func IfZeroInt(a, b int) int {
	if a == 0 {
		return b
	}
	return a
}

func IfZeroInt64(a, b int64) int64 {
	if a == 0 {
		return b
	}
	return a
}

func IfZeroFloat64(a, b float64) float64 {
	if a == 0 {
		return b
	}
	return a
}

func IfEmptyString(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

func IES(a, b string) string {
	return IfEmptyString(a, b)
}
