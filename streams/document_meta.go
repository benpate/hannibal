package streams

import "github.com/benpate/rosetta/mapof"

/******************************************
 * Metadata Getters
 ******************************************/

// Meta returns a pointer to the metadata associated with this document.
func (document Document) Meta() *mapof.Any {
	return &document.metadata
}

// MetaBool returns a bool from the document's metadata.
func (document Document) MetaBool(name string) bool {
	return document.metadata.GetBool(name)
}

// MetaFloat returns a float64 from the document's metadata.
func (document Document) MetaFloat(name string) float64 {
	return document.metadata.GetFloat(name)
}

// MetaInt returns an int from the document's metadata.
func (document Document) MetaInt(name string) int {
	return document.metadata.GetInt(name)
}

// MetaInt64 returns an int64 from the document's metadata.
func (document Document) MetaInt64(name string) int64 {
	return document.metadata.GetInt64(name)
}

// MetaString returns a string from the document's metadata.
func (document Document) MetaString(name string) string {
	return document.metadata.GetString(name)
}

/******************************************
 * Metadata Setters
 ******************************************/

// MetaSet replaces the metadata associated with this document.
func (document *Document) MetaSet(value mapof.Any) {
	document.metadata = value
}

// MetaAdd adds the specified metadata to the document.
func (document *Document) MetaAdd(value mapof.Any) {
	for key, value := range value {
		document.metadata[key] = value
	}
}

// MetaSetBool sets a bool value in the document's metadata.
func (document *Document) MetaSetBool(name string, value bool) {
	document.metadata.SetBool(name, value)
}

// MetaSetFloat sets a float64 value in the document's metadata.
func (document *Document) MetaSetFloat(name string, value float64) {
	document.metadata.SetFloat(name, value)
}

// MetaSetInt sets an int value in the document's metadata.
func (document *Document) MetaSetInt(name string, value int) {
	document.metadata.SetInt(name, value)
}

// MetaSetInt64 sets an int64 value in the document's metadata.
func (document *Document) MetaSetInt64(name string, value int64) {
	document.metadata.SetInt64(name, value)
}

// MetaSetString sets a string value in the document's metadata.
func (document *Document) MetaSetString(name string, value string) {
	document.metadata.SetString(name, value)
}
