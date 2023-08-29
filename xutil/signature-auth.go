package xutil

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrTokenFormat    = errors.New("token format error")
	ErrTokenLength    = errors.New("token length error")
	ErrTokenExpired   = errors.New("token expired error")
	ErrTokenSignature = errors.New("token signature error")
)

func GenerateSignature(salt string, key string) string {
	return fmt.Sprintf("%s%x", salt, sha256.Sum256([]byte(salt+key)))
}

func VerifySignature(token string, saltSize int, key string) error {
	if len(token) != saltSize+64 {
		return ErrTokenLength
	}

	salt := token[:saltSize]

	if token != GenerateSignature(salt, key) {
		return ErrTokenSignature
	}

	return nil
}

func GenerateTimeSignature(key string) string {
	salt := strconv.FormatInt(time.Now().Unix(), 10)

	return GenerateSignature(salt, key)
}

func VerifyTimeSignature(token string, ttl time.Duration, key string) error {
	if len(token) != 74 {
		return ErrTokenLength
	}

	saltSize := 10

	timestamp, err := strconv.Atoi(token[:saltSize])
	if err != nil {
		return ErrTokenFormat
	}

	if time.Now().Add(-ttl).Unix() > int64(timestamp) {
		return ErrTokenExpired
	}

	return VerifySignature(token, saltSize, key)
}
