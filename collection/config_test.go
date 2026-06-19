package collection

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/assert"
)

// TestNewConfig confirms each option sets the matching Config field, and that a
// config with no options is empty.
func TestNewConfig(t *testing.T) {

	empty := NewConfig()
	assert.Equal(t, Config{}, empty)

	full := NewConfig(
		WithSSEEndpoint("https://example.com/sse"),
		WithAttributedTo("https://example.com/actor"),
		WithAudience("https://example.com/followers"),
	)

	assert.Equal(t, "https://example.com/sse", full.SSEEndpoint)
	assert.Equal(t, "https://example.com/actor", full.AttributedTo)
	assert.Equal(t, "https://example.com/followers", full.Audience)
}

// TestConfig_Apply confirms Apply writes only the populated fields into the
// result map, leaving empty fields out.
func TestConfig_Apply(t *testing.T) {

	t.Run("all fields", func(t *testing.T) {
		config := Config{
			SSEEndpoint:  "https://example.com/sse",
			AttributedTo: "https://example.com/actor",
			Audience:     "https://example.com/followers",
		}

		result := mapof.NewAny()
		config.Apply(&result)

		assert.Equal(t, "https://example.com/actor", result[vocab.PropertyAttributedTo])
		assert.Equal(t, "https://example.com/followers", result[vocab.PropertyAudience])
		assert.Equal(t, "https://example.com/sse", result[vocab.PropertyEventStream])
	})

	t.Run("empty config writes nothing", func(t *testing.T) {
		result := mapof.NewAny()
		Config{}.Apply(&result)
		assert.Empty(t, result, "an empty config must not add any keys")
	})

	t.Run("partial config writes only populated fields", func(t *testing.T) {
		result := mapof.NewAny()
		Config{AttributedTo: "https://example.com/actor"}.Apply(&result)

		assert.Contains(t, result, vocab.PropertyAttributedTo)
		assert.NotContains(t, result, vocab.PropertyAudience)
		assert.NotContains(t, result, vocab.PropertyEventStream)
	})
}
