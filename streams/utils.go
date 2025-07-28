package streams

import (
	"iter"
	"strings"
)

// joinIterators combines multiple rangeFunc iterators into a single iterator.
func joinIterators[T any](iterators ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, iterator := range iterators {
			for value := range iterator {
				if !yield(value) {
					return
				}
			}
		}
	}
}

func indexOfNoCase(str string, substring string, startPosition int) int {

	if startPosition < 0 {
		startPosition = 0
	}

	if startPosition >= len(str) {
		return -1
	}

	if len(substring) == 0 {
		return -1
	}

	if startPosition > 0 {
		str = str[startPosition:]
	}

	str = strings.ToLower(str)
	substring = strings.ToLower(substring)

	return strings.Index(str, substring)
}
