package xlock

import (
	"sync"
	"time"
)

type locker struct {
	owner     string
	timestamp int64
}

var (
	lockers   = make(map[string]locker)
	lockersMu = sync.Mutex{}

	timeout int64 = 3600
)

// SetTimeout 设置锁超时时间，单位：秒
func SetTimeout(d int64) {
	timeout = d
}

func Lock(key string) bool {
	return OwnerLock("", key)
}

func Unlock(key string) {
	OwnerUnlock("", key)
}

func OwnerLock(owner string, key string) bool {
	lockersMu.Lock()
	defer lockersMu.Unlock()

	lk, ok := lockers[key]
	if ok || now()-lk.timestamp < timeout {
		return false
	}

	lockers[key] = locker{
		owner:     owner,
		timestamp: now(),
	}

	return true
}

func now() int64 {
	return time.Now().Unix()
}

func OwnerUnlock(owner string, key string) bool {
	lockersMu.Lock()
	defer lockersMu.Unlock()

	lk, ok := lockers[key]
	if !ok {
		return true
	}

	if lk.owner != owner {
		return false
	}

	delete(lockers, key)

	return true
}
