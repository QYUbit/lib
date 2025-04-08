package sets

import "fmt"

// Implentation of the set data structure. See: https://en.wikipedia.org/wiki/Set_(abstract_data_type).
type Set[T comparable] struct {
	items map[T]bool
}

// String representation of a set.
func (s Set[T]) String() string {
	return fmt.Sprint(ToSlice(s))
}

// Returns a new set containing all the items of the slice items.
func NewSet[T comparable](items []T) Set[T] {
	s := Set[T]{items: make(map[T]bool)}
	for _, item := range items {
		s.items[item] = true
	}
	return s
}

// Adds items to the set s. If an item already exists the item will be ignored.
func Add[T comparable](s Set[T], items ...T) Set[T] {
	for _, item := range items {
		s.items[item] = true
	}
	return s
}

// Removes items from the set s if present.
func Remove[T comparable](s Set[T], items ...T) Set[T] {
	for _, item := range items {
		delete(s.items, item)
	}
	return s
}

// Reports whether the item search is present in the set s.
func Exists[T comparable](s Set[T], search T) bool {
	_, exists := s.items[search]
	return exists
}

// Returns all items of the set s as a slice.
func ToSlice[T comparable](s Set[T]) []T {
	var slice []T
	for item := range s.items {
		slice = append(slice, item)
	}
	return slice
}

// Calls the callback function fn for each item in the set s.
func ForEachFunc[T comparable](s Set[T], fn func(T)) {
	for item := range s.items {
		fn(item)
	}
}

// Reports wheter the sets s1 and s2 are eaqual. Ignores order.
func Equals[T comparable](s1 Set[T], s2 Set[T]) bool {
	if len(s1.items) != len(s2.items) {
		return false
	}
	for v1 := range s1.items {
		for v2 := range s2.items {
			if v1 == v2 {
				break
			}
			return false
		}
	}
	return true
}

// Return a new set for those items in the set s which satisfy the callback fn.
func FilterFunc[T comparable](s Set[T], fn func(any) bool) Set[T] {
	var newItems []T
	for item := range s.items {
		if fn(item) {
			newItems = append(newItems, item)
		}
	}
	return NewSet(newItems)
}
