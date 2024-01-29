package types

import (
	"sync"
	"sync/atomic"
)

// Set represents a thread-safe set data structure that stores unique elements of type T.
type Set[T any] struct {
	size int64
	m    sync.Map
}

// NewSet creates a new instance of the Set data structure.
func NewSet[T any]() *Set[T] {
	return &Set[T]{}
}

// Put adds an element to the set.
// It returns a boolean indicating whether the element was added successfully (true if added, false if already exists).
func (s *Set[T]) Put(v T) bool {
	_, ok := s.m.LoadOrStore(v, struct{}{})
	if !ok {
		atomic.AddInt64(&s.size, 1)
	}
	return !ok
}

// Pop removes an element from the set.
// It returns the removed element and a boolean indicating whether the element existed in the set.
func (s *Set[T]) Pop(v T) (T, bool) {
	_, ok := s.m.LoadAndDelete(v)
	if ok {
		atomic.AddInt64(&s.size, -1)
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
	_, ok := s.m.Load(v)
	return ok
}

// Size returns the current size of the set.
func (s *Set[T]) Size() int64 {
	return atomic.LoadInt64(&s.size)
}

// Range iterates over all elements in the set and calls the provided function for each element.
// It stops iteration if the function returns false.
func (s *Set[T]) Range(f func(t T) bool) {
	s.m.Range(func(key, _ any) bool {
		return f(key.(T))
	})
}
