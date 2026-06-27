package streams

// UnmarshalItems normalizes a JSON-LD items value into a slice, returning TRUE if successful.
func UnmarshalItems(data any) ([]any, bool) {
	result, ok := data.([]any)
	return result, ok
}
