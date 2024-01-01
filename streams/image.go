package streams

import (
	"github.com/benpate/hannibal/unit"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/convert"
)

// https://www.w3.org/ns/activitystreams#Image
type Image struct {
	value any
}

// NewImage creates a new Image object from a JSON-LD value (string, map[string]any, or []any)
func NewImage(value any) Image {

	switch typed := value.(type) {

	case Document:
		return NewImage(typed.value.Raw())

	case unit.Value:
		return NewImage(typed.Raw())

	case Image:
		return typed

	case string:
		return Image{value: typed}

	case map[string]any:
		return Image{value: typed}

	case []any:
		return Image{value: typed}
	}

	return Image{""}
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-href
func (image Image) URL() string {

	switch typed := image.value.(type) {

	case string:
		return typed

	case map[string]any:

		if url := convert.String(typed[vocab.PropertyURL]); url != "" {
			return url
		}

		if href := convert.String(typed[vocab.PropertyHref]); href != "" {
			return href
		}

	case []any:
		if len(typed) > 0 {
			return NewImage(typed[0]).URL()
		}
	}

	return ""
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-summary
func (image Image) Summary() string {

	switch typed := image.value.(type) {

	case map[string]any:
		return convert.String(typed[vocab.PropertySummary])

	case []any:
		if len(typed) > 0 {
			return NewImage(typed[0]).Summary()
		}
	}

	return ""
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-mediatype
func (image Image) MediaType() string {

	switch typed := image.value.(type) {

	case map[string]any:
		return convert.String(typed[vocab.PropertyMediaType])

	case []any:
		if len(typed) > 0 {
			return NewImage(typed[0]).MediaType()
		}
	}

	return ""
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-height
func (image Image) Height() int {

	switch typed := image.value.(type) {

	case map[string]any:
		return convert.Int(typed[vocab.PropertyHeight])

	case []any:
		if len(typed) > 0 {
			return NewImage(typed[0]).Height()
		}
	}

	return 0
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-width
func (image Image) Width() int {

	switch typed := image.value.(type) {

	case map[string]any:
		return convert.Int(typed[vocab.PropertyWidth])

	case []any:
		if len(typed) > 0 {
			return NewImage(typed[0]).Width()
		}
	}

	return 0
}

// IsNil returns TRUE if this image is nil (having no URL)
func (image Image) IsNil() bool {
	return image.URL() == ""
}

// NotNil returns TRUE if this image has a URL
func (image Image) NotNil() bool {
	return !image.IsNil()
}

// HasHeight returns TRUE if this image has a height defined
func (image Image) HasHeight() bool {
	return image.Height() > 0
}

// HasWidth returns TRUE if this image has a width defined
func (image Image) HasWidth() bool {
	return image.Width() > 0
}

// HasDimensions returns TRUE if this image has both a height and width defined
func (image Image) HasDimensions() bool {
	return image.HasHeight() && image.HasWidth()
}

// AspectRatio calculates the aspect ratio of the image (width / height)
// If height and width are not available, then 0 is returned
func (image Image) AspectRatio() float64 {
	if image.HasDimensions() {
		return float64(image.Width()) / float64(image.Height())
	}

	return 0
}
