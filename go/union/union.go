package union

import (
	"errors"
	"fmt"
	"slices"
)

var ErrValueNotInOptions = errors.New("value not in options")

// Implemetation of the union data structure.
type Union[T comparable] struct {
	options []T
	index   int
}

// String representation of an union.
func (u Union[T]) String() string {
	v, defined := Get(u)
	if defined {
		return fmt.Sprintf("%v", *v)
	}
	return "undefined"
}

// Returns a new undefined union.
func NewUnion[T comparable](options []T) Union[T] {
	return Union[T]{
		options: options,
		index:   -1,
	}
}

// Returns a pointer to the current value of the union u and reports whether it is defined.
func Get[T comparable](u Union[T]) (*T, bool) {
	if Defined(u) {
		return &u.options[u.index], true
	}
	return nil, false
}

// Sets the union u to the value v. Returns an error when v is not in options.
func Set[T comparable](u *Union[T], v T) error {
	if slices.Contains(u.options, v) {
		u.index = slices.Index(u.options, v)
		return nil
	}
	return ErrValueNotInOptions
}

// Returns possible values for the union u.
func Options[T comparable](u Union[T]) []T {
	return u.options
}

// Resets the union u to an undefined state.
func Reset[T comparable](u *Union[T]) {
	u.index = -1
}

// Reports whether the union u is defined.
func Defined[T comparable](u Union[T]) bool {
	return u.index >= 0 && u.index < len(u.options)
}

// Reports whether two unions are equal.
func Equals[T comparable](u1, u2 Union[T]) bool {
	if len(u1.options) != len(u2.options) {
		return false
	}
	for i := range u1.options {
		if u1.options[i] != u2.options[i] {
			return false
		}
	}
	return u1.index == u2.index
}
