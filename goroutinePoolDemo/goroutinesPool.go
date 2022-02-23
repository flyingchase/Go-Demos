package goroutinepooldemo

import (
	"errors"
	"math"
)

var (
	ErrInvalidPoolSize      = errors.New("invalid size for pool")
	ErrInvalidPoolExpiry    = errors.New("invalid expiry for pool")
	ErrPoolClosed           = errors.New("this pool has been already Closed")
	defaultGoroutinePool, _ = NewPool(DefaultGoroutinePoolSize)
)

const (
	// 每隔DefaultCleanIntervalTime就清理
	DefaultCleanIntervalTime = 1
	Closed                   = 1
	DefaultGoroutinePoolSize = math.MaxInt32
)

// package pool and worker
func Submit(task func()) error {
	return defaultGoroutinePool.Submit(task)
}
func Running() int {
	return defaultGoroutinePool.Runing()

}
func Cap() int {
	return defaultGoroutinePool.Cap()
}
func Free() int {
	return defaultGoroutinePool.Free()
}

func Release() {
	_ = defaultGoroutinePool.Release()
}
