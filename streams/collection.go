package streams

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// Collection is a subtype of Object that represents ordered or unordered sets of Object or Link instances.
// https://www.w3.org/ns/activitystreams#Collection
type Collection struct {
	Context    Context `json:"@context,omitempty"    bson:"context,omitempty"`
	ID         string  `json:"id,omitempty"          bson:"id,omitempty"`
	Type       string  `json:"type,omitempty"        bson:"type,omitempty"`
	Summary    string  `json:"summary,omitempty"     bson:"summary,omitempty"`    // A natural language summarization of the object encoded as HTML. Multiple language tagged summaries may be provided.
	TotalItems int     `json:"totalItems,omitempty"  bson:"totalItems,omitempty"` // A non-negative integer specifying the total number of objects contained by the logical view of the collection. This number might not reflect the actual number of items serialized within the Collection object instance.
	Current    string  `json:"current,omitempty"     bson:"current,omitempty"`    // In a paged Collection, indicates the page that contains the most recently updated member items.
	First      string  `json:"first,omitempty"       bson:"first,omitempty"`      // In a paged Collection, indicates the furthest preceeding page of items in the collection.
	Last       string  `json:"last,omitempty"        bson:"last,omitempty"`       // In a paged Collection, indicates the furthest proceeding page of the collection.
	Items      []any   `json:"items,omitempty"       bson:"items,omitempty"`      // Identifies the items contained in a collection. The items might be ordered or unordered.
}

func NewCollection() Collection {
	return Collection{
		Context: DefaultContext(),
		Type:    vocab.CoreTypeCollection,
	}
}

/******************************************
 * JSON Marshalling
 ******************************************/

func (c *Collection) UnmarshalJSON(data []byte) error {

	result := make(map[string]any)

	if err := json.Unmarshal(data, &result); err != nil {
		return derp.Wrap(err, "activitystreams.Collection.UnmarshalJSON", "Error unmarshalling JSON", string(data))
	}

	return c.UnmarshalMap(result)
}

func (c *Collection) UnmarshalMap(data mapof.Any) error {

	if dataType := data.GetString("type"); dataType != vocab.CoreTypeCollection {
		return derp.NewInternalError("activitystreams.Collection.UnmarshalMap", "Invalid type", dataType)
	}

	c.Type = vocab.CoreTypeCollection
	c.Summary = data.GetString("summary")
	c.TotalItems = data.GetInt("totalItems")
	c.Current = data.GetString("current")
	c.First = data.GetString("first")
	c.Last = data.GetString("last")

	if dataItems, ok := data["items"]; ok {
		if items, ok := UnmarshalItems(dataItems); ok {
			c.Items = items
		}
	}

	return nil
}
