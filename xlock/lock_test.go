package xlock

import "testing"

func TestLock(t *testing.T) {
	t.Log(Lock("key1"))
	t.Log(Lock("key1"))
	Unlock("key1")
	t.Log(Lock("key1"))
	t.Log(Lock("key2"))
}
