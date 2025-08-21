package streams

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
)

// nolint:unused
type testClient struct {
	data mapof.Any
}

// nolint:unused
func (client testClient) SetRootClient(rootClient Client) {}

// nolint:unused
func (client testClient) Load(uri string, options ...any) (Document, error) {

	if value, ok := client.data[uri]; ok {
		return NewDocument(value, WithClient(client)), nil
	}

	return NilDocument(), derp.InternalError("hannibal.streams.testClient.Load", "Unknown URI", uri)
}
