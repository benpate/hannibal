package streams

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// CollectionPage is used to represent distinct subsets of items from a Collection. Refer to the Activity Streams 2.0 Core for a complete description of the CollectionPage object.
// https://www.w3.org/ns/activitystreams#CollectionPage
type CollectionPage struct {
	Context    Context `json:"@context,omitempty"    bson:"context,omitempty"`
	Type       string  `json:"type,omitempty"        bson:"type,omitempty"`       // Identifies the Object or Link type. (CollectionPage, OrderedCollectionPage)
	ID         string  `json:"id,omitempty"          bson:"id,omitempty"`         // Provides the globally unique identifier for an Object or Link.
	Summary    string  `json:"summary,omitempty"     bson:"summary,omitempty"`    // A natural language summarization of the object encoded as HTML. Multiple language tagged summaries may be provided.
	TotalItems int     `json:"totalItems,omitempty"  bson:"totalItems,omitempty"` // A non-negative integer specifying the total number of objects contained by the logical view of the collection. This number might not reflect the actual number of items serialized within the Collection object instance.
	Current    string  `json:"current,omitempty"     bson:"current,omitempty"`    // In a paged Collection, indicates the page that contains the most recently updated member items.
	First      string  `json:"first,omitempty"       bson:"first,omitempty"`      // In a paged Collection, indicates the furthest preceding page of items in the collection.
	Last       string  `json:"last,omitempty"        bson:"last,omitempty"`       // In a paged Collection, indicates the furthest proceeding page of the collection.
	PartOf     string  `json:"partOf,omitempty"      bson:"partOf,omitempty"`     // dentifies the Collection to which a CollectionPage objects items belong.
	Prev       string  `json:"prev,omitempty"        bson:"prev,omitempty"`       // In a paged Collection, identifies the previous page of items.
	Next       string  `json:"next,omitempty"        bson:"next,omitempty"`       // In a paged Collection, indicates the next page of items.
	Items      []any   `json:"items,omitempty"       bson:"items,omitempty"`      // Identifies the items contained in a collection. The items might be ordered or unordered.
}

func NewCollectionPage(pageID string) CollectionPage {
	return CollectionPage{
		Context: DefaultContext(),
		ID:      pageID,
		Type:    vocab.CoreTypeCollectionPage,
	}
}

/******************************************
 * JSON Marshalling
 ******************************************/
func (c *CollectionPage) UnmarshalJSON(data []byte) error {

	result := mapof.NewAny()

	if err := json.Unmarshal(data, &result); err != nil {
		return derp.Wrap(err, "activitystreams.CollectionPage.UnmarshalJSON", "Error unmarshalling JSON", string(data))
	}

	return c.UnmarshalMap(result)
}

func (c *CollectionPage) UnmarshalMap(data mapof.Any) error {

	if dataType := data.GetString("type"); dataType != vocab.CoreTypeCollectionPage {
		return derp.InternalError("activitystreams.CollectionPage.UnmarshalMap", "Invalid type", dataType)
	}

	c.Type = vocab.CoreTypeCollectionPage
	c.ID = data.GetString("id")
	c.Summary = data.GetString("summary")
	c.TotalItems = data.GetInt("totalItems")
	c.Current = data.GetString("current")
	c.First = data.GetString("first")
	c.Last = data.GetString("last")
	c.PartOf = data.GetString("partOf")
	c.Prev = data.GetString("prev")
	c.Next = data.GetString("next")

	if dataItems, ok := data["items"]; ok {
		if items, ok := UnmarshalItems(dataItems); ok {
			c.Items = items
		}
	}

	return nil
}
