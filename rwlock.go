package rwlock

import (
	"sync"
)

// A lock for a file. Can have multiple interested clients, represented in the counter.
// Will be freed when counter is zero.
type FileLock struct {
	counter uint32
	mutex   sync.RWMutex
}

// A file locking context. Usually there will be only 1 of these
type LockContext struct {
	mapLock sync.RWMutex
	locks   map[string]*FileLock
}

// Clean up the locks map after releasing for a path name
func (ctx *LockContext) cleanup(name string) {

	// Lock the mapping of file name to file lock
	ctx.mapLock.RLock()
	defer ctx.mapLock.RUnlock()

	lock, ok := ctx.locks[name]

	// Decrement the client counter, and remove if no one cares anymore
	if ok {
		lock.counter--
		if lock.counter == 0 {
			delete(ctx.locks, name)
		}
	}
}

// Get or create a new lock for `name`
func (ctx *LockContext) getOrCreateLock(name string) *FileLock {

	// Acquire the lock to read current file locks
	ctx.mapLock.RLock()
	defer ctx.mapLock.RUnlock()

	lock, ok := ctx.locks[name]

	// Inrement the counter if it exists, or create a new lock with 1 interest
	if ok {
		lock.counter++
	} else {
		// We need to write to the file map
		ctx.mapLock.RUnlock()
		ctx.mapLock.Lock()
		defer ctx.mapLock.Unlock()
		lock = &FileLock{counter: 1}
		ctx.locks[name] = lock
	}

	return lock
}

// Aquire the reader lock, and return a function that unlocks
func (ctx *LockContext) AcquireRead(name string) func() {
	lock := ctx.getOrCreateLock(name)
	lock.mutex.RLock()

	return func() {
		lock.mutex.RUnlock()
		ctx.cleanup(name)
	}
}

// Acquire the writer lock
func (ctx *LockContext) AcquireWrite(name string) func() {
	lock := ctx.getOrCreateLock(name)
	lock.mutex.Lock()

	return func() {
		lock.mutex.Unlock()
		ctx.cleanup(name)
	}
}

// Execute a function with the specified permissions.
// Permissions are a mapping of lock name to a boolean which is true for write permission, or false for read
func (ctx *LockContext) WithPermissions(permissions map[string]bool, f func()) {
	// Wait group for acquiring multiple permissions
	var wg sync.WaitGroup

	// This is to lock the releasers map
	var lock sync.Mutex

	// A mapping of file name to a function to release permissions
	var releasers map[string]func()

	// acquire permissions for a file and add its releaser
	acquire := func(name string, writer bool) {
		lock.Lock()
		defer lock.Unlock()
		defer wg.Done()
		if writer {
			release := ctx.AcquireWrite(name)
			releasers[name] = release
		} else {
			release := ctx.AcquireRead(name)
			releasers[name] = release
		}
	}

	defer func() {
		// Release all permissions
		for _, releaser := range releasers {
			releaser()
		}
	}()

	// Acquire permissions concurrently
	wg.Add(len(permissions))
	for name, writer := range permissions {
		go acquire(name, writer)
	}
	wg.Wait()

	// Execute the provided function
	f()
}
