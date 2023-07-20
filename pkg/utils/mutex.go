package utils

import "sync"

func ReturningMutex[T any](mutex *sync.Mutex, f func() T) T {
	mutex.Lock()
	result := f()
	mutex.Unlock()
	return result
}

func UseMutex(mutex *sync.Mutex, f func()) {
	mutex.Lock()
	f()
	mutex.Unlock()
}
