package filedir

import (
	"sync"
)

type fileDirectorySync struct {
	mu    sync.Mutex
	locks map[string]*fileLock
}

type fileLock struct {
	mu      sync.RWMutex
	counter int
}

func NewFileDirectorySync() (*fileDirectorySync, error) {
	return &fileDirectorySync{
		mu:    sync.Mutex{},
		locks: make(map[string]*fileLock),
	}, nil
}

func (dir *fileDirectorySync) getFileLock(path string) *fileLock {
	dir.mu.Lock()
	defer dir.mu.Unlock()

	lock, exists := dir.locks[path]
	if !exists {
		lock = &fileLock{}
		dir.locks[path] = lock
	}
	lock.counter++
	return lock
}

func (dir *fileDirectorySync) releaseFileLock(path string) {
	dir.mu.Lock()
	defer dir.mu.Unlock()

	lock, exists := dir.locks[path]
	if !exists {
		return
	}

	lock.counter--
	if lock.counter == 0 {
		delete(dir.locks, path)
	}
}

func (dir *fileDirectorySync) RLock(path string) {
	lock := dir.getFileLock(path)
	lock.mu.RLock()
}

func (dir *fileDirectorySync) RUnlock(path string) {
	dir.mu.Lock()
	lock, exists := dir.locks[path]
	dir.mu.Unlock()

	if !exists {
		return
	}

	lock.mu.RUnlock()
	dir.releaseFileLock(path)
}

func (dir *fileDirectorySync) Lock(path string) {
	lock := dir.getFileLock(path)
	lock.mu.Lock()
}

func (dir *fileDirectorySync) Unlock(path string) {
	dir.mu.Lock()
	lock, exists := dir.locks[path]
	dir.mu.Unlock()

	if !exists {
		return
	}

	lock.mu.Unlock()
	dir.releaseFileLock(path)
}
