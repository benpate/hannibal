package streams

// CollectionHeader is an opinionated format for generating/parsing the header information
// of a Collection.  It does not include any actual items, just the `totalItems` and `first` page URL.
type CollectionHeader struct {
	Context    Context `json:"@context,omitempty"    bson:"context,omitempty"`    // JSON-LD context to use
	ID         string  `json:"id,omitempty"          bson:"id,omitempty"`         // ID/URL of the Collection
	Type       string  `json:"type,omitempty"        bson:"type,omitempty"`       // Type of the Collection ("Collection" or "OrderedCollection")
	TotalItems int     `json:"totalItems,omitempty"  bson:"totalItems,omitempty"` // A non-negative integer specifying the total number of objects contained by the logical view of the collection. This number might not reflect the actual number of items serialized within the Collection object instance.
	First      string  `json:"first,omitempty"       bson:"first,omitempty"`      // In a paged Collection, indicates the furthest preceding page of items in the collection.
}
