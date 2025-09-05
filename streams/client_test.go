package streams

import (
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
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

func (client testClient) Save(document Document) error {
	return nil
}

func (client testClient) Delete(documentID string) error {
	return nil
}

func TestTestClient(t *testing.T) {

	// this is just a hack to make the "unused" linting messages go away
	client := testClient{}

	document := NewDocument(nil, WithClient(client))

	require.Nil(t, client.Save(document))
	require.Nil(t, client.Delete(document.ID()))
}
