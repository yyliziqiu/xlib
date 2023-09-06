package xlock

import (
	"sync"
)

func NewSyncLocker() *SyncLocker {
	return &SyncLocker{
		locks:   make(map[string]metadata),
		timeout: 3600,
	}
}

func NewSyncLockerWithTimeout(timeout int64) *SyncLocker {
	return &SyncLocker{
		locks:   make(map[string]metadata),
		timeout: timeout,
	}
}

type SyncLocker struct {
	locks   map[string]metadata
	locksMu sync.Mutex

	timeout int64
}

func (l *SyncLocker) Lock(key string) bool {
	return l.OwnerLock("", key)
}

func (l *SyncLocker) Unlock(key string) {
	l.OwnerUnlock("", key)
}

func (l *SyncLocker) OwnerLock(owner string, key string) bool {
	l.locksMu.Lock()
	defer l.locksMu.Unlock()

	lk, ok := l.locks[key]
	if ok || now()-lk.timestamp < l.timeout {
		return false
	}

	l.locks[key] = metadata{
		owner:     owner,
		timestamp: now(),
	}

	return true
}

func (l *SyncLocker) OwnerUnlock(owner string, key string) bool {
	l.locksMu.Lock()
	defer l.locksMu.Unlock()

	lk, ok := l.locks[key]
	if !ok {
		return true
	}

	if lk.owner != owner {
		return false
	}

	delete(l.locks, key)

	return true
}
