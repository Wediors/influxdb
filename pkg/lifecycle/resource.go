package lifecycle

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// Resource keeps track of references and has some compile time debug hooks
// to help diagnose leaks. It keeps track of if it is open or not and allows
// blocking until all references are released.
type Resource struct {
	sem sync.RWMutex  // counts outstanding references
	mu  sync.Mutex    // held for changes to ch
	ch  chan struct{} // non-nil when opened. closed when closing.
}

// Open Marks the reference counter as open. Should not be called concurrently
// with Close, Opened.
func (res *Resource) Open() {
	res.mu.Lock()
	res.ch = make(chan struct{})
	res.mu.Unlock()
}

// Close waits for any outstanding references and marks the reference counter
// as closed, so that Acquire returns an error. Should not be called
// concurrently with Open or Opened.
func (res *Resource) Close() {
	res.mu.Lock()
	if res.ch != nil {
		close(res.ch) // signal any references.
		res.ch = nil
	}
	res.mu.Unlock()

	res.sem.Lock() // wait for any oustanding references and block any calls to Acquire.
	res.sem.Unlock()
}

// Opened returns true if the resource is currently open. Should not be called
// concurrently with Open or Close.
func (res *Resource) Opened() bool {
	res.mu.Lock()
	opened := res.ch != nil
	res.mu.Unlock()
	return opened
}

// Acquire returns a Reference used to keep alive some resource.
func (res *Resource) Acquire() (*Reference, error) {
	res.mu.Lock()
	res.sem.RLock()
	if res.ch == nil {
		res.sem.RUnlock()
		res.mu.Unlock()
		return nil, resourceClosed()
	}
	// RLock intentionally left open.
	res.mu.Unlock()

	return live.track(&Reference{res: res}), nil
}

// Reference is an open reference for some resource.
type Reference struct {
	res *Resource
	id  uint64
}

// Release causes the Reference to be freed. It is safe to call multiple times.
func (r *Reference) Release() {
	// Inline a sync.Once using the res pointer as the flag for if we have
	// called unlock or not. This reduces the size of a ref.
	addr := (*unsafe.Pointer)(unsafe.Pointer(&r.res))
	old := atomic.LoadPointer(addr)
	if old != nil && atomic.CompareAndSwapPointer(addr, old, nil) {
		live.untrack(r)
		(*Resource)(old).sem.RUnlock()
	}
}

// Closing returns a channel that is closed when the associated resource
// is closing. If the reference is released, it returns nil. It should
// not be called concurrently with Close or Release.
func (r *Reference) Closing() <-chan struct{} { return r.res.ch }

// Close makes a Reference an io.Closer. It is safe to call multiple times.
func (r *Reference) Close() error {
	r.Release()
	return nil
}

// References is a helper to aggregate a group of references.
type References []*Reference

// Release releases all of the references. It is safe to call multiple times.
func (rs References) Release() {
	for _, r := range rs {
		r.Release()
	}
}

// Close makes references an io.Closer. It is safe to call multiple times.
func (rs References) Close() error {
	rs.Release()
	return nil
}
