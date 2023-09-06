package xlock

import (
	"time"
)

type Locker interface {
	Lock(key string) bool
	Unlock(key string)
	OwnerLock(owner string, key string) bool
	OwnerUnlock(owner string, key string) bool
}

type metadata struct {
	owner     string
	timestamp int64
}

func now() int64 {
	return time.Now().Unix()
}

var defaultLocker Locker = NewSyncLocker()

func SetDefaultLocker(locker Locker) {
	defaultLocker = locker
}

func Lock(key string) bool {
	return defaultLocker.Lock(key)
}

func Unlock(key string) {
	defaultLocker.Unlock(key)
}

func OwnerLock(owner string, key string) bool {
	return defaultLocker.OwnerLock(owner, key)
}

func OwnerUnlock(owner string, key string) bool {
	return defaultLocker.OwnerUnlock(owner, key)
}
