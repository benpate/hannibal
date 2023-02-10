package jsonld

// Cache is a simple interace for caching JSON-LD documents.
// There is a default implementation in the "cache" package,
// but this can be replaced with any other implementation.
type Cache interface {
	Get(string) map[string]any
	Set(string, map[string]any)
	Delete(string)
}

var cache Cache

func UseCache(c Cache) {
	cache = c
}
