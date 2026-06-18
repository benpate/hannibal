package property

import (
	"testing"
	"time"

	"github.com/benpate/hannibal"
	"github.com/stretchr/testify/assert"
)

// TestTime exercises the full Value contract for the Time type.
func TestTime(t *testing.T) {

	when := time.Date(2024, time.January, 22, 15, 4, 5, 0, time.UTC)
	value := Time(when)

	assert.True(t, value.IsTime())
	assert.True(t, when.Equal(value.Time()))
	assert.True(t, IsTime(value))

	// Scalar Value contract.
	assert.Equal(t, Nil{}, value.Get("anything"))
	assert.Equal(t, Map{"key": "v"}, value.Set("key", "v"))
	assert.Equal(t, value, value.Head())
	assert.Equal(t, Nil{}, value.Tail())
	assert.Equal(t, 1, value.Len())
	assert.Equal(t, map[string]any{}, value.Map())
	assert.Equal(t, value, value.Clone())

	// Raw round-trips to the underlying time.Time.
	assert.True(t, when.Equal(value.Raw().(time.Time)))

	// String uses the shared W3C/HTTP time format -- compare against the same
	// helper the type itself uses so the two can never drift apart.
	assert.Equal(t, hannibal.TimeFormat(when), value.String())
}

// TestTime_IsNil confirms Time reports nil only for the zero time.
func TestTime_IsNil(t *testing.T) {
	assert.True(t, Time(time.Time{}).IsNil())
	assert.False(t, Time(time.Now()).IsNil())
}
