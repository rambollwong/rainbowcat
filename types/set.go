package types

import (
	"sync"
	"sync/atomic"
)

type Set[T any] struct {
	size int64
	m    sync.Map
}

func NewSet[T any]() *Set[T] {
	return &Set[T]{}
}

func (s *Set[T]) Put(v T) bool {
	_, ok := s.m.LoadOrStore(v, struct{}{})
	if !ok {
		atomic.AddInt64(&s.size, 1)
	}
	return !ok
}

func (s *Set[T]) Pop(v T) (T, bool) {
	t, ok := s.m.LoadAndDelete(v)
	if ok {
		atomic.AddInt64(&s.size, -1)
	}
	return t.(T), ok
}

func (s *Set[T]) Remove(v T) bool {
	_, ok := s.Pop(v)
	return ok
}

func (s *Set[T]) Exist(v T) bool {
	_, ok := s.m.Load(v)
	return ok
}

func (s *Set[T]) Size() int64 {
	return atomic.LoadInt64(&s.size)
}
