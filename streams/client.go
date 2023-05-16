package streams

type Client interface {
	Load(uri string) (Document, error)
}
