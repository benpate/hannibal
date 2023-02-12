package jsonld

import (
	"time"

	"github.com/benpate/hannibal/vocab"
)

/******************************************
 * Document Getters
 ******************************************/

func (reader Reader) Actor() Reader {
	return reader.Get("actor")
}

func (reader Reader) ActorID() string {
	return reader.Actor().ID()
}

func (reader Reader) Activity() Reader {
	return reader.Get("activity")
}

/******************************************
 * Property Getters
 ******************************************/

func (reader Reader) ID() string {
	return reader.Get(vocab.PropertyID).AsString()
}

func (reader Reader) Attachment() Reader {
	return reader.Get(vocab.PropertyAttachment)
}

func (reader Reader) AttributedTo() Reader {
	return reader.Get(vocab.PropertyAttributedTo)
}

func (reader Reader) Audience() Reader {
	return reader.Get(vocab.PropertyAudience)
}

func (reader Reader) Bcc() Reader {
	return reader.Get(vocab.PropertyBCC)
}

func (reader Reader) Bto() Reader {
	return reader.Get(vocab.PropertyBTo)
}

func (reader Reader) Cc() Reader {
	return reader.Get(vocab.PropertyCC)
}

func (reader Reader) Context() Reader {
	return reader.Get(vocab.PropertyContext)
}

func (reader Reader) Current() Reader {
	return reader.Get(vocab.PropertyCurrent)
}

func (reader Reader) First() Reader {
	return reader.Get(vocab.PropertyFirst)
}

func (reader Reader) Generator() Reader {
	return reader.Get(vocab.PropertyGenerator)
}

func (reader Reader) Icon() Reader {
	return reader.Get(vocab.PropertyIcon)
}

func (reader Reader) Image() Reader {
	return reader.Get(vocab.PropertyImage)
}

func (reader Reader) InReplyTo() Reader {
	return reader.Get(vocab.PropertyInReplyTo)
}

func (reader Reader) Instrument() Reader {
	return reader.Get(vocab.PropertyInstrument)
}

func (reader Reader) Last() Reader {
	return reader.Get(vocab.PropertyLast)
}

func (reader Reader) Location() Reader {
	return reader.Get(vocab.PropertyLocation)
}

func (reader Reader) Items() Reader {
	return reader.Get(vocab.PropertyItems)
}

func (reader Reader) OneOf() Reader {
	return reader.Get(vocab.PropertyOneOf)
}

func (reader Reader) AnyOf() Reader {
	return reader.Get(vocab.PropertyAnyOf)
}

func (reader Reader) Closed() Reader {
	return reader.Get(vocab.PropertyClosed)
}

func (reader Reader) Origin() Reader {
	return reader.Get(vocab.PropertyOrigin)
}

func (reader Reader) Next() Reader {
	return reader.Get(vocab.PropertyNext)
}

func (reader Reader) Object() Reader {
	return reader.Get(vocab.PropertyObject)
}

func (reader Reader) Prev() Reader {
	return reader.Get(vocab.PropertyPrev)
}

func (reader Reader) Preview() Reader {
	return reader.Get(vocab.PropertyPreview)
}

func (reader Reader) Result() Reader {
	return reader.Get(vocab.PropertyResult)
}

func (reader Reader) Replies() Reader {
	return reader.Get(vocab.PropertyReplies)
}

func (reader Reader) Tag() Reader {
	return reader.Get(vocab.PropertyTag)
}

func (reader Reader) Target() Reader {
	return reader.Get(vocab.PropertyTarget)
}

func (reader Reader) To() Reader {
	return reader.Get(vocab.PropertyTo)
}

func (reader Reader) Type() string {
	return reader.Get(vocab.PropertyType).AsString()
}

// TODO: HIGH: Special handling for URL properties
// strings => {href: <string>}
// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-url
func (reader Reader) Url() Reader {
	return reader.Get(vocab.PropertyURL)
}

func (reader Reader) Accuracy() float64 {
	return reader.Get(vocab.PropertyAccuracy).AsFloat()
}

func (reader Reader) Altitude() float64 {
	return reader.Get(vocab.PropertyAltitude).AsFloat()
}

func (reader Reader) Content() string {
	return reader.Get(vocab.PropertyContent).AsString()
}

// TODO: Re-Implement Language Maps
func (reader Reader) Name() string {
	return reader.Get(vocab.PropertyName).AsString()
}

// TODO: Implement Durations per
// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-duration
func (reader Reader) Duration() string {
	return reader.Get(vocab.PropertyDuration).AsString()
}

func (reader Reader) Height() int {
	return reader.Get(vocab.PropertyHeight).AsInt()
}

func (reader Reader) Href() string {
	return reader.Get(vocab.PropertyHref).AsString()
}

func (reader Reader) Hreflang() string {
	return reader.Get(vocab.PropertyHrefLang).AsString()
}

func (reader Reader) PartOf() Reader {
	return reader.Get(vocab.PropertyPartOf)
}

func (reader Reader) Latitude() float64 {
	return reader.Get(vocab.PropertyLatitude).AsFloat()
}

func (reader Reader) Longitude() float64 {
	return reader.Get(vocab.PropertyLongitude).AsFloat()
}

func (reader Reader) MediaType() string {
	return reader.Get(vocab.PropertyMediaType).AsString()
}

func (reader Reader) EndTime() time.Time {
	return reader.Get(vocab.PropertyEndTime).AsTime()
}

func (reader Reader) Published() time.Time {
	return reader.Get(vocab.PropertyPublished).AsTime()
}

func (reader Reader) StartTime() time.Time {
	return reader.Get(vocab.PropertyStartTime).AsTime()
}

func (reader Reader) Radius() float64 {
	return reader.Get(vocab.PropertyRadius).AsFloat()
}

func (reader Reader) Rel() Reader {
	return reader.Get(vocab.PropertyRel)
}

func (reader Reader) StartIndex() int {
	return reader.Get(vocab.PropertyStartIndex).AsInt()
}

// // TODO: Implement Language Maps
func (reader Reader) Summary() string {
	return reader.Get(vocab.PropertySummary).AsString()
}

func (reader Reader) TotalItems() int {
	return reader.Get(vocab.PropertyTotalItems).AsInt()
}

func (reader Reader) Units() string {
	return reader.Get(vocab.PropertyUnits).AsString()
}

func (reader Reader) Updated() time.Time {
	return reader.Get(vocab.PropertyUpdated).AsTime()
}

func (reader Reader) Width() int {
	return reader.Get(vocab.PropertyWidth).AsInt()
}

func (reader Reader) Subject() Reader {
	return reader.Get(vocab.PropertySubject)
}

func (reader Reader) Relationship() string {
	return reader.Get(vocab.PropertyRelationship).AsString()
}

func (reader Reader) Describes() Reader {
	return reader.Get(vocab.PropertyDescribes)
}

func (reader Reader) FormerType() string {
	return reader.Get(vocab.PropertyFormerType).AsString()
}

func (reader Reader) Deleted() time.Time {
	return reader.Get(vocab.PropertyDeleted).AsTime()
}
