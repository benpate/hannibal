package streams

import (
	"slices"

	"github.com/benpate/derp"
)

// OptionsClient is a Client wrapper that adds a fixed set of options to every Load, and binds each
// document it returns to itself, so that loads made later from those documents carry the options too.
type OptionsClient struct {
	inner   Client
	options []any
}

// NewOptionsClient returns a Client that loads through inner, placing options ahead of each Load's own.
func NewOptionsClient(inner Client, options ...any) Client {

	// NOTE: Unlike a stack layer's constructor, this never calls inner.SetRootClient.  A binding
	// wraps the top of a live stack, and re-rooting that stack would reach every load it makes.

	// A variadic argument may be the caller's own slice, so keep a private copy
	return OptionsClient{
		inner:   inner,
		options: slices.Clone(options),
	}
}

// Load retrieves a document from the inner client with this client's options added, and binds the
// result to this client so that its implicit loads carry the same options.
func (client OptionsClient) Load(url string, options ...any) (Document, error) {

	const location = "hannibal.streams.OptionsClient.Load"

	if client.inner == nil {
		return NilDocument(), derp.Internal(location, "Inner client must not be nil", url)
	}

	// RULE: Bound options come first, so a Load's own options follow them and win wherever the
	// last value applies.  Concat always allocates; append could write into the shared slice.
	result, err := client.inner.Load(url, slices.Concat(client.options, options)...)

	if err != nil {
		return result, derp.Wrap(err, location, "Loading document", url)
	}

	// Every document reached from this one loads through this client, options and all
	return result.AddOptions(WithClient(client)), nil
}

// Save passes the document to the inner client.
func (client OptionsClient) Save(document Document) error {

	const location = "hannibal.streams.OptionsClient.Save"

	if client.inner == nil {
		return derp.Internal(location, "Inner client must not be nil")
	}

	return client.inner.Save(document)
}

// Delete passes the document ID to the inner client.
func (client OptionsClient) Delete(documentID string) error {

	const location = "hannibal.streams.OptionsClient.Delete"

	if client.inner == nil {
		return derp.Internal(location, "Inner client must not be nil", documentID)
	}

	return client.inner.Delete(documentID)
}

// SetRootClient passes the top-level client down to the inner client, as every stack layer does.
func (client OptionsClient) SetRootClient(rootClient Client) {

	// An OptionsClient placed inside a stack must pass this on, or the layers below it keep an
	// older root, and documents they return would load around everything above them
	if client.inner != nil {
		client.inner.SetRootClient(rootClient)
	}
}

// Verify that OptionsClient satisfies the Client interface.
var _ Client = OptionsClient{}
