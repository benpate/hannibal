package streams

import (
	"time"

	"github.com/benpate/hannibal/vocab"
)

func (document Document) ID() string {
	return document.Get(vocab.PropertyID).AsString()
}

func (document Document) Actor() Document {
	return document.Get("actor")
}

func (document Document) ActorID() string {
	return document.Actor().ID()
}

func (document Document) Activity() Document {
	return document.Get("activity")
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
	return document.Get(vocab.PropertyDeleted).AsTime()
}

func (document Document) Describes() Document {
	return document.Get(vocab.PropertyDescribes)
}

func (document Document) First() Document {
	return document.Get(vocab.PropertyFirst)
}

func (document Document) FormerType() string {
	return document.Get(vocab.PropertyFormerType).AsString()
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

	return icon.AsString()
}

func (document Document) Image() Document {
	return document.Get(vocab.PropertyImage)
}

func (document Document) ImageURL() string {
	image := document.Image().Head()

	if image.IsObject() {
		return image.URL()
	}

	return image.AsString()
}

func (document Document) InReplyTo() Document {
	return document.Get(vocab.PropertyInReplyTo)
}

func (document Document) Instrument() Document {
	return document.Get(vocab.PropertyInstrument)
}

func (document Document) Last() Document {
	return document.Get(vocab.PropertyLast)
}

func (document Document) Location() Document {
	return document.Get(vocab.PropertyLocation)
}

func (document Document) Items() Document {
	return document.Get(vocab.PropertyItems)
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
	if result := document.Get(vocab.PropertyType).AsString(); result != "" {
		return result
	}
	return vocab.Unknown
}

// TODO: HIGH: Special handling for URL properties
// strings => {href: <string>}
// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-url
func (document Document) Url() Document {
	return document.Get(vocab.PropertyURL)
}

func (document Document) Accuracy() float64 {
	return document.Get(vocab.PropertyAccuracy).AsFloat()
}

func (document Document) Altitude() float64 {
	return document.Get(vocab.PropertyAltitude).AsFloat()
}

func (document Document) Content() string {
	return document.Get(vocab.PropertyContent).AsString()
}

// TODO: Re-Implement Language Maps
func (document Document) Name() string {
	return document.Get(vocab.PropertyName).AsString()
}

// TODO: Implement Durations per
// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-duration
func (document Document) Duration() string {
	return document.Get(vocab.PropertyDuration).AsString()
}

func (document Document) Height() int {
	return document.Get(vocab.PropertyHeight).AsInt()
}

func (document Document) Href() string {
	return document.Get(vocab.PropertyHref).AsString()
}

func (document Document) Hreflang() string {
	return document.Get(vocab.PropertyHrefLang).AsString()
}

func (document Document) PartOf() Document {
	return document.Get(vocab.PropertyPartOf)
}

func (document Document) Latitude() float64 {
	return document.Get(vocab.PropertyLatitude).AsFloat()
}

func (document Document) Longitude() float64 {
	return document.Get(vocab.PropertyLongitude).AsFloat()
}

func (document Document) MediaType() string {
	return document.Get(vocab.PropertyMediaType).AsString()
}

func (document Document) EndTime() time.Time {
	return document.Get(vocab.PropertyEndTime).AsTime()
}

func (document Document) Published() time.Time {
	return document.Get(vocab.PropertyPublished).AsTime()
}

func (document Document) StartTime() time.Time {
	return document.Get(vocab.PropertyStartTime).AsTime()
}

func (document Document) Radius() float64 {
	return document.Get(vocab.PropertyRadius).AsFloat()
}

// Rel is expected to be a string, but this function
// returns a document because it may contain multiple values (rel:["canonical", "preview"])
func (document Document) Rel() Document {
	return document.Get(vocab.PropertyRel)
}

func (document Document) Relationship() string {
	return document.Get(vocab.PropertyRelationship).AsString()
}

func (document Document) StartIndex() int {
	return document.Get(vocab.PropertyStartIndex).AsInt()
}

func (document Document) Subject() Document {
	return document.Get(vocab.PropertySubject)
}

// // TODO: Implement Language Maps
func (document Document) Summary() string {
	return document.Get(vocab.PropertySummary).AsString()
}

func (document Document) TotalItems() int {
	return document.Get(vocab.PropertyTotalItems).AsInt()
}

func (document Document) URL() string {
	return document.Get(vocab.PropertyURL).AsString()
}

func (document Document) Units() string {
	return document.Get(vocab.PropertyUnits).AsString()
}

func (document Document) Updated() time.Time {
	return document.Get(vocab.PropertyUpdated).AsTime()
}

func (document Document) Width() int {
	return document.Get(vocab.PropertyWidth).AsInt()
}
