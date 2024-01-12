package xrand

import (
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func String(length int) string {
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))

	var stringBuilder strings.Builder
	stringBuilder.Grow(length)

	for i := 0; i < length; i++ {
		char := charset[rander.Intn(len(charset))]
		stringBuilder.WriteByte(char)
	}

	return stringBuilder.String()
}
