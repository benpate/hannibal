package streams

type Client interface {
	Load(uri string, defaultValue map[string]any) (Document, error)
}
