package streams

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// OrderedCollection is a subtype of Collection in which members of the logical collection are assumed to always be strictly ordered.
// https://www.w3.org/ns/activitystreams#OrderedCollection
type OrderedCollection struct {
	Context      Context `json:"@context,omitempty"     bson:"@context,omitempty"`
	ID           string  `json:"id,omitempty"           bson:"id,omitempty"`
	Type         string  `json:"type,omitempty"         bson:"type,omitempty"`
	Summary      string  `json:"summary,omitempty"      bson:"summary,omitempty"`      // A natural language summarization of the object encoded as HTML. Multiple language tagged summaries may be provided.
	TotalItems   int     `json:"totalItems,omitempty"   bson:"totalItems,omitempty"`   // A non-negative integer specifying the total number of objects contained by the logical view of the collection. This number might not reflect the actual number of items serialized within the Collection object instance.
	OrderedItems []any   `json:"orderedItems,omitempty" bson:"orderedItems,omitempty"` // Identifies the items contained in a collection. The items might be ordered or unordered.
	Current      string  `json:"current,omitempty"      bson:"current,omitempty"`      // In a paged Collection, indicates the page that contains the most recently updated member items.
	First        string  `json:"first,omitempty"        bson:"first,omitempty"`        // In a paged Collection, indicates the furthest preceding page of items in the collection.
	Last         string  `json:"last,omitempty"         bson:"last,omitempty"`         // In a paged Collection, indicates the furthest proceeding page of the collection.
}

func NewOrderedCollection(collectionID string) OrderedCollection {
	return OrderedCollection{
		Context:      DefaultContext(),
		Type:         vocab.CoreTypeOrderedCollection,
		ID:           collectionID,
		OrderedItems: make([]any, 0),
	}
}

func (c *OrderedCollection) UnmarshalJSON(data []byte) error {

	result := mapof.NewAny()

	if err := json.Unmarshal(data, &result); err != nil {
		return derp.Wrap(err, "activitystreams.OrderedCollection.UnmarshalJSON", "Error unmarshalling JSON", string(data))
	}

	return c.UnmarshalMap(result)
}

func (c *OrderedCollection) UnmarshalMap(data mapof.Any) error {

	if dataType := data.GetString("type"); dataType != vocab.CoreTypeOrderedCollection {
		return derp.InternalError("activitystreams.OrderedCollection.UnmarshalMap", "Invalid type", dataType)
	}

	c.Type = vocab.CoreTypeOrderedCollection
	c.Summary = data.GetString("summary")
	c.TotalItems = data.GetInt("totalItems")
	c.Current = data.GetString("current")
	c.First = data.GetString("first")
	c.Last = data.GetString("last")

	if dataItems, ok := data["items"]; ok {
		if items, ok := UnmarshalItems(dataItems); ok {
			c.OrderedItems = items
		}
	}

	return nil
}
