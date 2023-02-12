package jsonld

import (
	"encoding/json"
	"strconv"
	"time"
)

type Reader struct {
	value  any
	client *Client
}

func NilReader() Reader {
	return Reader{
		value:  nil,
		client: nil,
	}
}

// Get returns a sub-property of the current object
func (r Reader) Get(key string) Reader {

	switch typed := r.value.(type) {

	case string:
		if key == "id" {
			return r.client.NewReader(typed)
		} else {
			return r.Load().Get(key)
		}

	case int:
		if key == "id" {
			return r.client.NewReader(typed)
		}

	case int64:
		if key == "id" {
			return r.client.NewReader(typed)
		}

	case map[string]any:
		return r.client.NewReader(typed[key])

	case []any:
		if len(typed) > 0 {
			return r.client.NewReader(typed[0])
		}
	}

	return NilReader()
}

// AsBool returns the current object as a floating-point value
func (r Reader) AsBool() bool {

	switch typed := r.value.(type) {

	case string:
		return typed == "true"

	case int:
		return typed != 0

	case int64:
		return typed != 0

	case float64:
		return typed != 0

	case bool:
		return typed

	case map[string]any:
		return r.Get("id").AsBool()

	case []any:
		return r.Get("id").AsBool()
	}

	return false
}

// AsFloat returns the current object as an integer value
func (r Reader) AsFloat() float64 {

	switch typed := r.value.(type) {

	case string:
		if result, err := strconv.ParseFloat(typed, 64); err != nil {
			return result
		}

	case int:
		return float64(typed)

	case int64:
		return float64(typed)

	case float64:
		return typed

	case bool:
		if typed {
			return 1
		}
		return 0

	case map[string]any:
		return r.Get("id").AsFloat()

	case []any:
		return r.Get("id").AsFloat()
	}

	return 0
}

// AsInt returns the current object as an integer value
func (r Reader) AsInt() int {

	switch typed := r.value.(type) {

	case string:
		if result, err := strconv.Atoi(typed); err != nil {
			return result
		}

	case int:
		return typed

	case int64:
		return int(typed)

	case float64:
		return int(typed)

	case bool:
		if typed {
			return 1
		}
		return 0

	case map[string]any:
		return r.Get("id").AsInt()

	case []any:
		return r.Get("id").AsInt()
	}

	return 0
}

// AsString returns the current object as a string value
func (r Reader) AsString() string {

	switch typed := r.value.(type) {

	case string:
		return typed

	case int:
		return strconv.Itoa(typed)

	case int64:
		return strconv.FormatInt(typed, 10)

	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)

	case bool:
		return strconv.FormatBool(typed)

	case map[string]any:
		return r.Get("id").AsString()

	case []any:
		return r.Get("id").AsString()
	}

	return ""
}

// AsTime returns the current object as a time value
func (r Reader) AsTime() time.Time {

	switch typed := r.value.(type) {

	case string:
		if result, err := time.Parse(time.RFC3339, typed); err == nil {
			return result
		}

	case int:
		return time.Unix(int64(typed), 0)

	case int64:
		return time.Unix(typed, 0)

	case float64:
		return time.Unix(int64(typed), 0)
	}

	return time.Time{}
}

// Load retrieves a remote object if an ID is available
func (r Reader) Load() Reader {
	result, _ := r.client.Load(r.AsString())
	return result
}

// IsEmpty return TRUE if the current object is empty
func (r Reader) IsEmptyTail() bool {

	if slice, ok := r.value.([]any); ok {
		return len(slice) < 2
	}

	return true
}

// Tail returns a slice of all records after the first.
func (r Reader) Tail() Reader {

	if slice, ok := r.value.([]any); ok {

		if len(slice) > 1 {

			return Reader{
				value:  slice[1:],
				client: r.client,
			}
		}
	}

	return NilReader()
}

func (r Reader) Value() any {
	return r.value
}

func (r Reader) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.value)
}
