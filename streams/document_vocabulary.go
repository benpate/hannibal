package streams

import (
	"time"

	"github.com/benpate/hannibal/vocab"
)

/******************************************
 * ActivityStreams 2.0 Properties
 * https://www.w3.org/TR/activitystreams-vocabulary
 ******************************************/

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-id
func (document Document) ID() string {

	// Try the ActivityPub standard "id" property first
	if id := document.Get(vocab.PropertyID); !id.IsNil() {
		return id.String()
	}

	// Try the JSON-LD standard "@id" property second
	if id := document.Get(vocab.PropertyID_Alternate); !id.IsNil() {
		return document.Get(vocab.PropertyID).String()
	}

	// LOL, no.
	return ""
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-actor
func (document Document) Actor() Document {
	return document.Get(vocab.PropertyActor).MustLoad()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-attachment
func (document Document) Attachment() Document {
	return document.Get(vocab.PropertyAttachment)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-attributedto
func (document Document) AttributedTo() Document {
	return document.Get(vocab.PropertyAttributedTo).MustLoad()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-audience
func (document Document) Audience() Document {
	return document.Get(vocab.PropertyAudience)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-bcc
func (document Document) BCC() Document {
	return document.Get(vocab.PropertyBCC).MustLoad()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-bto
func (document Document) BTo() Document {
	return document.Get(vocab.PropertyBTo).MustLoad()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-cc
func (document Document) CC() Document {
	return document.Get(vocab.PropertyCC).MustLoad()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-context
// IMPORTANT: THIS IS NOT THE @context PROPERTY REQUIRED FOR EVERY JSON-LD DOCUMENT
func (document Document) Context() string {
	return document.Get(vocab.PropertyContext).String()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-current
func (document Document) Current() Document {
	return document.Get(vocab.PropertyCurrent)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-deleted
func (document Document) Deleted() time.Time {
	return document.Get(vocab.PropertyDeleted).Time()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-describes
func (document Document) Describes() Document {
	return document.Get(vocab.PropertyDescribes)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-first
func (document Document) First() Document {
	return document.Get(vocab.PropertyFirst)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-formertype
func (document Document) FormerType() string {
	return document.Get(vocab.PropertyFormerType).String()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-generator
func (document Document) Generator() Document {
	return document.Get(vocab.PropertyGenerator)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-icon
func (document Document) Icon() Image {
	return NewImage(document.Get(vocab.PropertyIcon))
}

// IconOrImage is a hybrid accessor that returns the "icon" property (if not nil),
// otherwise it returns the "image" property.  This is useful for working with different
// ActivityPub objects, which may use either property.
func (document Document) IconOrImage() Image {

	if icon := document.Icon(); icon.NotNil() {
		return icon
	}

	return document.Image()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-image
func (document Document) Image() Image {
	return NewImage(document.Get(vocab.PropertyImage))
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-inreplyto
func (document Document) InReplyTo() Document {
	return document.Get(vocab.PropertyInReplyTo)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-instrument
func (document Document) Instrument() Document {
	return document.Get(vocab.PropertyInstrument)
}

// Items returns the items collection for this Document.  If the
// document contains an "orderedItems" collection, then it is
// returned instead.
func (document Document) Items() Document {

	// Search the "orderedItems" property first (guessing this will be more common)
	if result := document.Get(vocab.PropertyOrderedItems); !result.IsNil() {
		return result
	}

	// Search the "items" property second (guessing this will be less common)
	if result := document.Get(vocab.PropertyItems); !result.IsNil() {
		return result
	}

	// Value not found :(
	return NilDocument()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-last
func (document Document) Last() Document {
	return document.Get(vocab.PropertyLast)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-location
func (document Document) Location() Document {
	return document.Get(vocab.PropertyLocation)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-oneof
func (document Document) OneOf() Document {
	return document.Get(vocab.PropertyOneOf)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-anyof
func (document Document) AnyOf() Document {
	return document.Get(vocab.PropertyAnyOf)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-closed
func (document Document) Closed() Document {
	return document.Get(vocab.PropertyClosed)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-origin
func (document Document) Origin() Document {
	return document.Get(vocab.PropertyOrigin)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-next
func (document Document) Next() Document {
	return document.Get(vocab.PropertyNext)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-object
func (document Document) Object() Document {
	return document.Get(vocab.PropertyObject)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-prev
func (document Document) Prev() Document {
	return document.Get(vocab.PropertyPrev)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-preview
func (document Document) Preview() Document {
	return document.Get(vocab.PropertyPreview)
}

func (document Document) PublicKey() Document {
	return document.Get(vocab.PropertyPublicKey)
}

func (document Document) PublicKeyPEM() string {
	return document.Get(vocab.PropertyPublicKeyPEM).String()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-result
func (document Document) Result() Document {
	return document.Get(vocab.PropertyResult)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-replies
func (document Document) Replies() Document {
	return document.Get(vocab.PropertyReplies)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-tag
func (document Document) Tag() Document {
	return document.Get(vocab.PropertyTag)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-target
func (document Document) Target() Document {
	return document.Get(vocab.PropertyTarget).MustLoad()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-to
func (document Document) To() Document {
	return document.Get(vocab.PropertyTo).MustLoad()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-type
func (document Document) Type() string {

	// Try the ActivityPub standard "type" property first
	if value := document.Get(vocab.PropertyType); !value.IsNil() {
		return value.String()
	}

	// Try the JSON-LD standard "@type" property second
	if value := document.Get(vocab.PropertyType_Alternate); !value.IsNil() {
		return value.String()
	}

	// LOL, Fail
	return vocab.Unknown
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-url
func (document Document) Url() Document {
	return document.Get(vocab.PropertyURL)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-accuracy
func (document Document) Accuracy() float64 {
	return document.Get(vocab.PropertyAccuracy).Float()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-altitude
func (document Document) Altitude() float64 {
	return document.Get(vocab.PropertyAltitude).Float()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-content
func (document Document) Content() string {
	return document.Get(vocab.PropertyContent).HTMLString()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-name
func (document Document) Name() string {
	// TODO: Re-Implement Language Maps
	return document.Get(vocab.PropertyName).String()
}

// TODO: Implement Durations per
// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-duration
func (document Document) Duration() string {

	return document.Get(vocab.PropertyDuration).String()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-height
func (document Document) Height() int {
	return document.Get(vocab.PropertyHeight).Int()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-href
func (document Document) Href() string {
	return document.Get(vocab.PropertyHref).String()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-hreflang
func (document Document) Hreflang() string {
	return document.Get(vocab.PropertyHrefLang).String()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-partof
func (document Document) PartOf() Document {
	return document.Get(vocab.PropertyPartOf)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-latitude
func (document Document) Latitude() float64 {
	return document.Get(vocab.PropertyLatitude).Float()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-longitude
func (document Document) Longitude() float64 {
	return document.Get(vocab.PropertyLongitude).Float()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-mediatype
func (document Document) MediaType() string {
	return document.Get(vocab.PropertyMediaType).String()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-endtime
func (document Document) EndTime() time.Time {
	return document.Get(vocab.PropertyEndTime).Time()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-published
func (document Document) Published() time.Time {
	return document.Get(vocab.PropertyPublished).Time()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-starttime
func (document Document) StartTime() time.Time {
	return document.Get(vocab.PropertyStartTime).Time()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-radius
func (document Document) Radius() float64 {
	return document.Get(vocab.PropertyRadius).Float()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-rel
// Rel is expected to be a string, but this function
// returns a document because it may contain multiple values (rel:["canonical", "preview"])
func (document Document) Rel() Document {
	return document.Get(vocab.PropertyRel)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-relationship
func (document Document) Relationship() string {
	return document.Get(vocab.PropertyRelationship).String()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-startindex
func (document Document) StartIndex() int {
	return document.Get(vocab.PropertyStartIndex).Int()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-subject
func (document Document) Subject() Document {
	return document.Get(vocab.PropertySubject)
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-summary
// TODO: Implement Language Maps
func (document Document) Summary() string {
	return document.Get(vocab.PropertySummary).String()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-totalitems
func (document Document) TotalItems() int {
	return document.Get(vocab.PropertyTotalItems).Int()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-width
func (document Document) URL() string {
	return document.Get(vocab.PropertyURL).String()
}

// URLOrID returns the URL of the document, if it exists, otherwise it returns the ID.
func (document Document) URLOrID() string {
	if url := document.URL(); url != "" {
		return url
	}
	return document.ID()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-units
func (document Document) Units() string {
	return document.Get(vocab.PropertyUnits).String()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-updated
func (document Document) Updated() time.Time {
	return document.Get(vocab.PropertyUpdated).Time()
}

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-width
func (document Document) Width() int {
	return document.Get(vocab.PropertyWidth).Int()
}
