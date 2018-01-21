package gonyan

import (
	"sync"
)

// mutex wraps a standard sync.Mutex basic functionalities together with a bool
// for easy deactivation and ease possible feature extensions.
type mutex struct {
	mtx      sync.Mutex
	disabled bool
}

// Lock proceeds with internal mutex lock if the mutex has not been disabled.
func (m *mutex) Lock() {
	if !m.disabled {
		m.mtx.Lock()
	}
}

// Unlock proceeds with internal mutex unlock if the mutex has not been
// disabled.
func (m *mutex) Unlock() {
	if !m.disabled {
		m.mtx.Unlock()
	}
}

// Disable sets the mutex state to disabled thus preventing the Lock and Unlock
// functions to actually make a difference.
// Beware of disabling the mutex while in actual use, especially during a
// critical section execution because, since Unlock won't work it will cause
// deadlocks.
func (m *mutex) Disable() {
	m.disabled = true
}

// Enable sets the mutex state to enabled reactivating the Lock and Unlock
// functions.
func (m *mutex) Enable() {
	m.disabled = false
}
