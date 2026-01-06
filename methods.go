package set

import (
	"maps"
	"reflect"
)

// Clone returns a shallowly copied set.
// Reference types will refer to the old object.
func Clone[T comparable](s *set[T]) *set[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	newSet := &set[T]{s: maps.Clone(s.s)}
	return newSet
}

// Difference returns the set of elements that are in the first set but not in the second.
func Difference[T comparable](original *set[T], different *set[T]) *set[T] {
	original.mu.RLock()
	defer original.mu.RUnlock()

	different.mu.RLock()
	defer different.mu.RUnlock()

	diffSet := NewEmptySet[T](0)

	for val := range original.s {
		if !different.exists(val) {
			diffSet.insert(val)
		}
	}

	return diffSet
}

// Intersection returns a set consisting of the common elements of both sets.
func Intersection[T comparable](s1 *set[T], s2 *set[T]) *set[T] {
	s1.mu.RLock()
	defer s1.mu.RUnlock()

	s2.mu.RLock()
	defer s2.mu.RUnlock()

	l1, l2 := len(s1.s), len(s2.s)
	if l1 > l2 {
		s1, s2 = s2, s1 // iterate over the smaller set
		l1, l2 = l2, l1
	}

	intersectionSet := NewEmptySet[T](0)

	for val := range s1.s {
		if s2.exists(val) {
			intersectionSet.insert(val)
		}
	}

	return intersectionSet
}

// Union returns a new set consisting of the elements of the given sets.
func Union[T comparable](s1 *set[T], s2 *set[T]) *set[T] {
	s1.mu.RLock()
	defer s1.mu.RUnlock()

	s2.mu.RLock()
	defer s2.mu.RUnlock()

	l1, l2 := len(s1.s), len(s2.s)
	if l1 < l2 {
		l1, l2 = l2, l1
	}

	unionSet := NewEmptySet[T](l1)

	for val := range s1.s {
		unionSet.insert(val)
	}
	for val := range s2.s {
		unionSet.insert(val)
	}

	return unionSet
}

// SymmetricDifference returns elements that are not common to the sets.
func SymmetricDifference[T comparable](s1 *set[T], s2 *set[T]) *set[T] {
	s1.mu.RLock()
	defer s1.mu.RUnlock()

	s2.mu.RLock()
	defer s2.mu.RUnlock()

	newSet := NewEmptySet[T](0)

	for val := range s1.s {
		if !s2.exists(val) {
			newSet.insert(val)
		}
	}
	for val := range s2.s {
		if !s1.exists(val) {
			newSet.insert(val)
		}
	}

	return newSet
}

// IsSubSet returns true if the first set is a subset of the second.
func IsSubSet[T comparable](s1 *set[T], s2 *set[T]) bool {
	s1.mu.RLock()
	defer s1.mu.RUnlock()

	s2.mu.RLock()
	defer s2.mu.RUnlock()

	if len(s1.s) < len(s2.s) {
		return false
	}

	for val := range s1.s {
		if !s2.Exists(val) {
			return false
		}
	}

	return true
}

// Equal comparison of two sets.
func Equal[T comparable](s1 *set[T], s2 *set[T]) bool {
	s1.mu.RLock()
	defer s1.mu.RUnlock()

	s2.mu.RLock()
	defer s2.mu.RUnlock()

	if len(s1.s) != len(s2.s) {
		return false
	}
	if s1 == nil || s2 == nil {
		return s1 == nil && s2 == nil
	}
	diffSet := Difference(s1, s2)
	return len(diffSet.s) == 0
}

// DeepEqual deep comparison of two sets.
func DeepEqual[T comparable](s1 *set[T], s2 *set[T]) bool {
	s1.mu.RLock()
	defer s1.mu.RUnlock()

	s2.mu.RLock()
	defer s2.mu.RUnlock()

	return reflect.DeepEqual(s1.s, s2.s)
}
