package types

import (
	"sync"
	"sync/atomic"
)

// Set represents a thread-safe set data structure that stores unique elements of type T.
// Backed by a regular map protected by sync.RWMutex, which is more efficient than
// sync.Map for most set workloads (frequent reads and writes on the same keys).
type Set[T comparable] struct {
	size atomic.Int64
	mu   sync.RWMutex
	m    map[T]struct{}
}

// NewSet creates a new instance of the Set data structure.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		m: make(map[T]struct{}),
	}
}

// Put adds an element to the set.
// It returns a boolean indicating whether the element was added successfully (true if added, false if already exists).
func (s *Set[T]) Put(v T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.m[v]; ok {
		return false
	}
	s.m[v] = struct{}{}
	s.size.Add(1)
	return true
}

// Pop removes an element from the set.
// It returns the removed element and a boolean indicating whether the element existed in the set.
func (s *Set[T]) Pop(v T) (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.m[v]; ok {
		delete(s.m, v)
		s.size.Add(-1)
		return v, true
	}
	return v, false
}

// Remove removes an element from the set.
// It returns a boolean indicating whether the element was successfully removed (true if removed, false if not found).
func (s *Set[T]) Remove(v T) bool {
	_, ok := s.Pop(v)
	return ok
}

// Exist checks if an element exists in the set.
// It returns a boolean indicating whether the element exists in the set (true if exists, false if not found).
func (s *Set[T]) Exist(v T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[v]
	return ok
}

// Size returns the current size of the set.
func (s *Set[T]) Size() int64 {
	return s.size.Load()
}

// Range iterates over all elements in the set and calls the provided function for each element.
// It stops iteration if the function returns false.
// Note: the callback is invoked under a read lock; avoid long-running or blocking operations inside it.
func (s *Set[T]) Range(f func(t T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k := range s.m {
		if !f(k) {
			return
		}
	}
}
