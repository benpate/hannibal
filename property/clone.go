package property

import (
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// cloneChild returns a deep copy of a value stored inside a Map or a Slice, so that no container
// in the result is shared with the original.
//
// Only containers are cloned. Everything else is returned by assignment: JSON scalars are
// immutable, and unrecognized values must NOT be routed through NewValue, whose final fallback is
// Nil{} -- a round-trip would silently erase them.
//
// Concrete types are preserved wherever they can be: a mapof.Any child clones to a mapof.Any, not
// to a bare map[string]any. The one exception is the reflection fallback at the bottom, which
// cannot reconstruct an arbitrary named type and so normalizes to map[string]any / []any.
func cloneChild(value any) any {

	switch typed := value.(type) {

	// Byte slices are data, not containers. They are checked before the reflection fallback
	// below, which would otherwise see a slice and shred them into []any.
	case []byte:
		return typed

	// Values clone themselves, and keep their own concrete type.
	case Value:
		return typed.Clone()

	case map[string]any:
		return cloneMap(typed)

	case mapof.Any:
		return mapof.Any(cloneMap(typed))

	case []any:
		return cloneSlice(typed)

	case sliceof.Any:
		return sliceof.Any(cloneSlice(typed))
	}

	// Wayward containers (primitive.A from Mongo, map[string]string, and friends) are named or
	// differently-typed containers that the switch above cannot match, so they need the same
	// reflection check NewValue uses. These normalize to map[string]any / []any.
	if convert.IsMap(value) {
		return cloneMap(convert.MapOfAny(value))
	}

	if convert.IsSlice(value) {
		return cloneSlice(convert.SliceOfAny(value))
	}

	// Scalars and anything else: copied by assignment.
	return value
}

// cloneMap deep-copies every entry of a string-keyed map into a new map.
func cloneMap(value map[string]any) map[string]any {

	result := make(map[string]any, len(value))

	for key, child := range value {
		result[key] = cloneChild(child)
	}

	return result
}

// cloneSlice deep-copies every element of a slice into a new slice.
func cloneSlice(value []any) []any {

	result := make([]any, len(value))

	for index, child := range value {
		result[index] = cloneChild(child)
	}

	return result
}
