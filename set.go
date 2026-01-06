package set

import (
	"fmt"
	"maps"
	"runtime"
	"strings"
	"sync"
)

var countWorkers int = runtime.NumCPU()

type set[T comparable] struct {
	s  map[T]struct{}
	mu sync.RWMutex
}

// NewEmptySet returns an empty set with the given capacity.
func NewEmptySet[T comparable](size int) *set[T] {
	return &set[T]{s: make(map[T]struct{}, size)}
}

// NewSet returns a set with the specified elements.
func NewSet[T comparable](e ...T) *set[T] {
	set := &set[T]{s: make(map[T]struct{}, len(e))}
	for _, val := range e {
		set.insert(val)
	}
	return set
}

func (s *set[T]) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	elements := make([]string, 0, len(s.s))
	for elem := range s.s {
		elements = append(elements, fmt.Sprintf("%v", elem))
	}

	return "Set[" + strings.Join(elements, " ") + "]"
}

// Empty returns an indication that the set has no elements, or that it is nil.
func (s *set[T]) Empty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.s) == 0
}

// Size Returns the number of elements in a set.
func (s *set[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.s)
}

// Insert inserts a new element into a set.
// If such an element was already in the set, then nothing will change.
func (s *set[T]) Insert(e ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.insert(e...)
}

// insert internal Insert method without a mutex
func (s *set[T]) insert(e ...T) {
	for _, val := range e {
		s.s[val] = struct{}{}
	}
}

// Exists checks if an element exists within a set.
func (s *set[T]) Exists(e T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.exists(e)
}

// exists internal Exists method without a mutex.
func (s *set[T]) exists(e T) bool {
	_, exists := s.s[e]
	return exists
}

// Delete removes an element from a set.
// If such an element is not present in the set, it will not cause panic.
func (s *set[T]) Delete(e T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.delete(e)
}

// delete internal Delete method without a mutex.
func (s *set[T]) delete(e T) {
	delete(s.s, e)
}

// Clear removes all elements from a set.
func (s *set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.s)
}

// clear internal Clear method without a mutex.
func (s *set[T]) clear() {
	clear(s.s)
}

// Pop removes and return a random item from the set.
func (s *set[T]) Pop() T {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := getZeroValue[T]()

	if s.Size() == 0 {
		return result
	}

	for val := range s.s {
		s.delete(val)
		result = val
		break
	}

	return result
}

// ExtendFromSlice expands the set with elements from the slice.
func (s *set[T]) ExtendFromSlice(slice []T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, val := range slice {
		s.insert(val)
	}
}

// ExtendFromMap adds keys from the map to the set
func (s *set[T]) ExtendFromMap(m map[T]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k := range m {
		s.insert(k)
	}
}

// UnionUpdate extends a set with elements from another set.
func (s *set[T]) UnionUpdate(src *set[T]) {
	src.mu.RLock()
	defer src.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	for val := range src.s {
		s.insert(val)
	}
}

// IntersectionUpdate removes elements if they do not exist in both sets.
func (s *set[T]) IntersectionUpdate(src *set[T]) {
	src.mu.RLock()
	defer src.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	for val := range s.s {
		if !src.Exists(val) {
			s.delete(val)
		}
	}
}

// DifferenceUpdate removes elements that are present in both sets.
func (s *set[T]) DifferenceUpdate(src *set[T]) {
	src.mu.RLock()
	defer src.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	for val := range s.s {
		if src.exists(val) {
			s.delete(val)
		}
	}
}

// ToSlice returns a slice consisting of the elements of the set.
func (s *set[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resultSlice := make([]T, 0, len(s.s))
	for val := range s.s {
		resultSlice = append(resultSlice, val)
	}
	return resultSlice
}

// ToMap returns a map where the key is an element of the set.
func (s *set[T]) ToMap() map[T]struct{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Clone(s.s)
}

// Do applies the passed function to each element of the set.
// The function is executed synchronously.
func (s *set[T]) Do(f func(T) T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for val := range s.s {
		s.delete(val)
		newVal := f(val)
		s.insert(newVal)
	}
}

// Do applies the passed function to each element of the set.
// The function starts several goroutines in parallel.
func (s *set[T]) DoAsync(f func(T) T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var wg sync.WaitGroup
	lenSet := len(s.s)
	countWorkers := getCountWorkers(lenSet)
	sem := make(chan struct{}, countWorkers)
	chValuesToInsert := make(chan T, lenSet)

	for val := range s.s {
		sem <- struct{}{}
		wg.Add(1)
		go func(v T) {
			defer wg.Done()
			chValuesToInsert <- f(v)
			<-sem
		}(val)
	}
	wg.Wait()
	close(chValuesToInsert)

	s.clear()
	for valToInsert := range chValuesToInsert {
		s.insert(valToInsert)
	}
}

// Filter leaves in the set only those elements that meet the given condition.
// The function is executed synchronously.
func (s *set[T]) Filter(f func(T) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for val := range s.s {
		if !f(val) {
			s.delete(val)
		}
	}
}

// FilterAsync leaves in the set only those elements that meet the given condition.
// The function starts several goroutines in parallel.
func (s *set[T]) FilterAsync(f func(T) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var wg sync.WaitGroup
	lenSet := len(s.s)
	countWorkers := getCountWorkers(lenSet)
	sem := make(chan struct{}, countWorkers)
	chValuesToDelete := make(chan T, lenSet)

	for val := range s.s {
		sem <- struct{}{}
		wg.Add(1)
		go func(v T) {
			defer wg.Done()
			if !f(v) {
				chValuesToDelete <- v
			}
			<-sem
		}(val)
	}
	wg.Wait()
	close(chValuesToDelete)

	for valToDelete := range chValuesToDelete {
		s.delete(valToDelete)
	}
}
