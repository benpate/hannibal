package streams

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
)

// TestImage_FromMap exercises every accessor against a fully-populated map value.
func TestImage_FromMap(t *testing.T) {

	image := NewImage(map[string]any{
		vocab.PropertyHref:      "https://example.com/pic.png",
		vocab.PropertySummary:   "A picture",
		vocab.PropertyMediaType: "image/png",
		vocab.PropertyWidth:     640,
		vocab.PropertyHeight:    480,
	})

	assert.Equal(t, "https://example.com/pic.png", image.URL())
	assert.Equal(t, "https://example.com/pic.png", image.Href())
	assert.Equal(t, "A picture", image.Summary())
	assert.Equal(t, "image/png", image.MediaType())
	assert.Equal(t, 640, image.Width())
	assert.Equal(t, 480, image.Height())
	assert.True(t, image.HasWidth())
	assert.True(t, image.HasHeight())
	assert.True(t, image.HasDimensions())
	assert.True(t, image.NotNil())
	assert.InDelta(t, 640.0/480.0, image.AspectRatio(), 0.0001)
}

// TestImage_HrefFallsBackToURL confirms Href reads the "url" property when "href"
// is absent.
func TestImage_HrefFallsBackToURL(t *testing.T) {

	image := NewImage(map[string]any{vocab.PropertyURL: "https://example.com/pic.png"})
	assert.Equal(t, "https://example.com/pic.png", image.Href())
}

// TestImage_FromString confirms a bare string value is treated as the URL.
func TestImage_FromString(t *testing.T) {

	image := NewImage("https://example.com/pic.png")
	assert.Equal(t, "https://example.com/pic.png", image.URL())
	assert.Equal(t, "", image.MediaType())
	assert.Equal(t, 0, image.Width())
}

// TestImage_FromSlice confirms a slice value resolves through its first element.
func TestImage_FromSlice(t *testing.T) {

	image := NewImage([]any{
		map[string]any{
			vocab.PropertyHref:      "https://example.com/first.png",
			vocab.PropertyMediaType: "image/jpeg",
			vocab.PropertySummary:   "First",
			vocab.PropertyWidth:     100,
			vocab.PropertyHeight:    200,
		},
	})

	assert.Equal(t, "https://example.com/first.png", image.URL())
	assert.Equal(t, "image/jpeg", image.MediaType())
	assert.Equal(t, "First", image.Summary())
	assert.Equal(t, 100, image.Width())
	assert.Equal(t, 200, image.Height())
}

// TestImage_Nil confirms nil/empty values produce an empty image with zero values.
func TestImage_Nil(t *testing.T) {

	image := NewImage(nil)
	assert.True(t, image.IsNil())
	assert.False(t, image.NotNil())
	assert.False(t, image.HasDimensions())
	assert.Equal(t, float64(0), image.AspectRatio())

	// An empty slice has no first element, so every accessor returns its zero value.
	empty := NewImage([]any{})
	assert.Equal(t, "", empty.URL())
	assert.Equal(t, "", empty.Summary())
	assert.Equal(t, 0, empty.Width())
}
