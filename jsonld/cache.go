package jsonld

type Cache interface {
	Get(string) map[string]any
	Set(string, map[string]any)
}
