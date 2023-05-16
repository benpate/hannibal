package streams

import (
	"time"

	"github.com/benpate/hannibal/vocab"
)

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

func (document Document) Actor() Document {
	return document.Get(vocab.PropertyActor)
}

func (document Document) ActorID() string {
	return document.Actor().ID()
}

func (document Document) Activity() Document {
	return document.Get(vocab.CoreTypeActivity)
}

func (document Document) Attachment() Document {
	return document.Get(vocab.PropertyAttachment)
}

func (document Document) AttributedTo() Document {
	return document.Get(vocab.PropertyAttributedTo)
}

func (document Document) Audience() Document {
	return document.Get(vocab.PropertyAudience)
}

func (document Document) Bcc() Document {
	return document.Get(vocab.PropertyBCC)
}

func (document Document) Bto() Document {
	return document.Get(vocab.PropertyBTo)
}

func (document Document) Cc() Document {
	return document.Get(vocab.PropertyCC)
}

func (document Document) Context() Document {
	return document.Get(vocab.PropertyContext)
}

func (document Document) Current() Document {
	return document.Get(vocab.PropertyCurrent)
}

func (document Document) Deleted() time.Time {
	return document.Get(vocab.PropertyDeleted).Time()
}

func (document Document) Describes() Document {
	return document.Get(vocab.PropertyDescribes)
}

func (document Document) First() Document {
	return document.Get(vocab.PropertyFirst)
}

func (document Document) FormerType() string {
	return document.Get(vocab.PropertyFormerType).String()
}

func (document Document) Generator() Document {
	return document.Get(vocab.PropertyGenerator)
}

func (document Document) Icon() Document {
	return document.Get(vocab.PropertyIcon)
}

func (document Document) IconURL() string {
	icon := document.Icon().Head()

	if icon.IsObject() {
		return icon.URL()
	}

	return icon.String()
}

func (document Document) Image() Document {
	return document.Get(vocab.PropertyImage)
}

func (document Document) ImageURL() string {
	image := document.Image().Head()

	if image.IsObject() {
		return image.URL()
	}

	return image.String()
}

func (document Document) InReplyTo() Document {
	return document.Get(vocab.PropertyInReplyTo)
}

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

func (document Document) Last() Document {
	return document.Get(vocab.PropertyLast)
}

func (document Document) Location() Document {
	return document.Get(vocab.PropertyLocation)
}

func (document Document) OneOf() Document {
	return document.Get(vocab.PropertyOneOf)
}

func (document Document) AnyOf() Document {
	return document.Get(vocab.PropertyAnyOf)
}

func (document Document) Closed() Document {
	return document.Get(vocab.PropertyClosed)
}

func (document Document) Origin() Document {
	return document.Get(vocab.PropertyOrigin)
}

func (document Document) Next() Document {
	return document.Get(vocab.PropertyNext)
}

func (document Document) Object() Document {
	return document.Get(vocab.PropertyObject)
}

func (document Document) ObjectID() string {
	return document.Get(vocab.PropertyObject).ID()
}

func (document Document) Prev() Document {
	return document.Get(vocab.PropertyPrev)
}

func (document Document) Preview() Document {
	return document.Get(vocab.PropertyPreview)
}

func (document Document) PublicKey() Document {
	return document.Get(vocab.PropertyPublicKey)
}

func (document Document) PublicKeyPEM() string {
	return document.Get(vocab.PropertyPublicKeyPEM).String()
}

func (document Document) Result() Document {
	return document.Get(vocab.PropertyResult)
}

func (document Document) Replies() Document {
	return document.Get(vocab.PropertyReplies)
}

func (document Document) Tag() Document {
	return document.Get(vocab.PropertyTag)
}

func (document Document) Target() Document {
	return document.Get(vocab.PropertyTarget)
}

func (document Document) TargetID() string {
	return document.Get(vocab.PropertyTarget).ID()
}

func (document Document) To() Document {
	return document.Get(vocab.PropertyTo)
}

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

func (document Document) Url() Document {

	value := document.Get(vocab.PropertyURL)

	if value.IsString() {
		return document.sub(map[string]any{vocab.PropertyHref: value})
	}

	return value
}

func (document Document) Accuracy() float64 {
	return document.Get(vocab.PropertyAccuracy).Float()
}

func (document Document) Altitude() float64 {
	return document.Get(vocab.PropertyAltitude).Float()
}

func (document Document) Content() string {
	return document.Get(vocab.PropertyContent).String()
}

// TODO: Re-Implement Language Maps
func (document Document) Name() string {
	return document.Get(vocab.PropertyName).String()
}

// TODO: Implement Durations per
// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-duration
func (document Document) Duration() string {
	return document.Get(vocab.PropertyDuration).String()
}

func (document Document) Height() int {
	return document.Get(vocab.PropertyHeight).Int()
}

func (document Document) Href() string {
	return document.Get(vocab.PropertyHref).String()
}

func (document Document) Hreflang() string {
	return document.Get(vocab.PropertyHrefLang).String()
}

func (document Document) PartOf() Document {
	return document.Get(vocab.PropertyPartOf)
}

func (document Document) Latitude() float64 {
	return document.Get(vocab.PropertyLatitude).Float()
}

func (document Document) Longitude() float64 {
	return document.Get(vocab.PropertyLongitude).Float()
}

func (document Document) MediaType() string {
	return document.Get(vocab.PropertyMediaType).String()
}

func (document Document) EndTime() time.Time {
	return document.Get(vocab.PropertyEndTime).Time()
}

func (document Document) Published() time.Time {
	return document.Get(vocab.PropertyPublished).Time()
}

func (document Document) StartTime() time.Time {
	return document.Get(vocab.PropertyStartTime).Time()
}

func (document Document) Radius() float64 {
	return document.Get(vocab.PropertyRadius).Float()
}

// Rel is expected to be a string, but this function
// returns a document because it may contain multiple values (rel:["canonical", "preview"])
func (document Document) Rel() Document {
	return document.Get(vocab.PropertyRel)
}

func (document Document) Relationship() string {
	return document.Get(vocab.PropertyRelationship).String()
}

func (document Document) StartIndex() int {
	return document.Get(vocab.PropertyStartIndex).Int()
}

func (document Document) Subject() Document {
	return document.Get(vocab.PropertySubject)
}

// // TODO: Implement Language Maps
func (document Document) Summary() string {
	return document.Get(vocab.PropertySummary).String()
}

func (document Document) TotalItems() int {
	return document.Get(vocab.PropertyTotalItems).Int()
}

func (document Document) URL() string {
	return document.Get(vocab.PropertyURL).String()
}

func (document Document) Units() string {
	return document.Get(vocab.PropertyUnits).String()
}

func (document Document) Updated() time.Time {
	return document.Get(vocab.PropertyUpdated).Time()
}

func (document Document) Width() int {
	return document.Get(vocab.PropertyWidth).Int()
}
