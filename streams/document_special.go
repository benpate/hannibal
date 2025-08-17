package streams

import (
	"strconv"
)

/******************************************
 * Type Detection
 ******************************************/

func (document Document) DocumentCategory() string {
	return DocumentCategory(document.Type())
}

// IsActivity returns TRUE if this document represents an Activity
func (document Document) IsActivity() bool {
	return IsActivity(document.Type())
}

// NotActivity returns TRUE if this document does NOT represent an Activity
func (document Document) NotActivity() bool {
	return !document.IsActivity()
}

// IsActor returns TRUE if this document represents an Actor
func (document Document) IsActor() bool {
	return IsActor(document.Type())
}

// NotActor returns TRUE if this document does NOT represent an Actor
func (document Document) NotActor() bool {
	return !document.IsActor()
}

// IsCollection returns TRUE if this document represents a Collection or CollectionPage
func (document Document) IsCollection() bool {
	return IsCollection(document.Type())
}

// NotCollection returns TRUE if the document does NOT represent a Collection or CollectionPage
func (document Document) NotCollection() bool {
	return !document.IsCollection()
}

// IsObject returns TRUE if this document represents an Object type (Article, Note, etc)
func (document Document) IsObject() bool {
	return IsObject(document.Type())
}

// NotObject returns TRUE if this document does NOT represent an Object type (Article, Note, etc)
func (document Document) NotObject() bool {
	return !document.IsObject()
}

// Metadata returns counts for various interactions: Announces, Replies, Likes, and Dislikes
func (document Document) Metadata() *Metadata {
	return &document.metadata
}

// HasIcon returns TRUE if this document has a valid Icon property
func (document Document) HasIcon() bool {
	return document.Icon().NotNil()
}

// HasImage returns TRUE if this document has a valid Image property
func (document Document) HasImage() bool {
	return document.Image().NotNil()
}

// HasContent returns TRUE if this document has a valid Content property
func (document Document) HasContent() bool {
	return document.Content() != ""
}

// HasSummary returns TRUE if this document has a valid Summary property
func (document Document) HasSummary() bool {
	return document.Summary() != ""
}

func (document Document) HasDimensions() bool {
	return document.Width() > 0 && document.Height() > 0
}

func (document Document) SummaryWithTagLinks() string {

	summary := document.Summary()

	if summary == "" {
		return ""
	}

	for tag := range document.Tag().Range() {
		href := tag.Href()

		if href == "" {
			continue
		}

		tagName := tag.Name()
		tagNameLength := len(tagName)

		if tagNameLength == 0 {
			continue
		}

		startPosition := 0
		for {

			index := indexOfNoCase(summary, tagName, startPosition)

			if index < 0 {
				break
			}

			tagLink := `<a href="` + href + `" target="_blank">` + tagName + `</a>`
			tagLinkLength := len(tagLink)

			summary = summary[:index] + tagLink + summary[index+tagNameLength:]

			startPosition = index + tagLinkLength
		}
	}

	return summary
}

func (document Document) AspectRatio() string {

	width := document.Width()
	height := document.Height()

	if width == 0 || height == 0 {
		return "auto"
	}

	ratio := float64(width) / float64(height)
	return strconv.FormatFloat(ratio, 'f', -1, 64)
}

// If this document is an activity (create, update, delete, etc), then
// this method returns the activity's Object.  Otherwise, it returns
// the document itself.
func (document Document) UnwrapActivity() Document {

	// If this is an "Activity" type, the dig deeper into the object
	// to find the actual document.
	// This is recursive because it's possible to have a deep tree
	// such as Announce > Create > Document. Looking at you, Lemmy...
	if document.IsActivity() {
		return document.Object().UnwrapActivity()
	}

	return document
}
