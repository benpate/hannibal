package streams

import "iter"

// Uniquer is a utility class that helps to identify unique values
type Uniquer[T comparable] struct {
	seen map[T]struct{}
}

// NewUniquer returns a fully initialized Uniquer object
func NewUniquer[T comparable]() *Uniquer[T] {
	return &Uniquer[T]{
		seen: make(map[T]struct{}),
	}
}

// IsUnique returns TRUE if the value has not been seen before.
// Subsequent calls to IsUnique() with the same value will return FALSE.
func (u *Uniquer[T]) IsUnique(id T) bool {

	if _, seen := u.seen[id]; seen {
		return false
	}

	u.seen[id] = struct{}{}
	return true
}

// IsDuplicate returns TRUE if the value has been seen before.
func (u *Uniquer[T]) IsDuplicate(id T) bool {
	return !u.IsUnique(id)
}

func (u *Uniquer[T]) Range(iterator iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {

		for item := range iterator {
			if u.IsUnique(item) {
				if !yield(item) {
					return // Stop iterating if the yield function returns false
				}
			}
		}
	}
}
