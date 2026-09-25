package clients

import (
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/require"
)

// TestCarpool_ConcurrentLoadsShareOneFetch confirms that 50 callers loading one URL at the same
// moment, each through its own client stack, cause one fetch and all receive the document.
func TestCarpool_ConcurrentLoadsShareOneFetch(t *testing.T) {

	synctest.Test(t, func(t *testing.T) {

		inner := newGateClient()
		carpool := NewCarpool()

		const riders = 50
		results := make([]streams.Document, riders)
		errs := make([]error, riders)

		var done sync.WaitGroup

		for index := range riders {
			done.Go(func() {
				results[index], errs[index] = carpool.Client(inner, testSigner).Load("https://example.com/actor")
			})
		}

		// Every rider is now either fetching or waiting on the one who is
		synctest.Wait()
		close(inner.gate)
		done.Wait()

		require.Equal(t, int64(1), inner.loads.Load())

		for index := range riders {
			require.NoError(t, errs[index])
			require.Equal(t, "https://example.com/actor", results[index].ID())
		}
	})
}

// TestCarpool_RidersGetTheirOwnDocument confirms that callers sharing one fetch each receive an
// independent copy, bound to their own client stack.
func TestCarpool_RidersGetTheirOwnDocument(t *testing.T) {

	synctest.Test(t, func(t *testing.T) {

		inner := newGateClient()
		carpool := NewCarpool()

		// Two stacks that sign as the same actor, standing in for two viewers on one domain
		aliceRoot := &mockInnerClient{}
		bobRoot := &mockInnerClient{}

		alice := carpool.Client(inner, testSigner)
		alice.SetRootClient(aliceRoot)

		bob := carpool.Client(inner, testSigner)
		bob.SetRootClient(bobRoot)

		var aliceResult, bobResult streams.Document
		var done sync.WaitGroup

		done.Go(func() { aliceResult, _ = alice.Load("https://example.com/actor") })
		done.Go(func() { bobResult, _ = bob.Load("https://example.com/actor") })

		synctest.Wait()
		close(inner.gate)
		done.Wait()

		require.Equal(t, int64(1), inner.loads.Load())

		// RULE: A follow-up load from a rider's document goes through that rider's own stack
		require.Same(t, aliceRoot, aliceResult.Client())
		require.Same(t, bobRoot, bobResult.Client())

		// RULE: Changing one rider's copy leaves the other's untouched
		aliceResult.Map()[vocab.PropertyName] = "Changed by Alice"
		aliceResult.HTTPHeader().Set("X-Alice", "yes")

		require.Equal(t, "Original", bobResult.Name())
		require.Empty(t, bobResult.HTTPHeader().Get("X-Alice"))
	})
}

// TestCarpool_OptionsLoadAlone confirms that a Load with options is passed straight to the inner
// client, options intact, and never shares a fetch.
func TestCarpool_OptionsLoadAlone(t *testing.T) {

	synctest.Test(t, func(t *testing.T) {

		inner := newGateClient()
		carpool := NewCarpool()

		var done sync.WaitGroup

		for range 3 {
			done.Go(func() {
				_, _ = carpool.Client(inner, testSigner).Load("https://example.com/actor", "refresh")
			})
		}

		synctest.Wait()
		close(inner.gate)
		done.Wait()

		require.Equal(t, int64(3), inner.loads.Load())
		require.Equal(t, []any{"refresh"}, inner.lastOptions())
	})
}

// TestCarpool_DifferentURLsLoadSeparately confirms that concurrent Loads of different URLs never
// wait on each other.
func TestCarpool_DifferentURLsLoadSeparately(t *testing.T) {

	synctest.Test(t, func(t *testing.T) {

		inner := newGateClient()
		carpool := NewCarpool()

		var done sync.WaitGroup

		done.Go(func() { _, _ = carpool.Client(inner, testSigner).Load("https://example.com/alice") })
		done.Go(func() { _, _ = carpool.Client(inner, testSigner).Load("https://example.com/bob") })

		synctest.Wait()
		close(inner.gate)
		done.Wait()

		require.Equal(t, int64(2), inner.loads.Load())
	})
}

// TestCarpool_DifferentSignersLoadSeparately confirms that concurrent Loads of one URL by different
// signers never share a fetch, so no caller receives a result made with another's credentials.
func TestCarpool_DifferentSignersLoadSeparately(t *testing.T) {

	synctest.Test(t, func(t *testing.T) {

		inner := newGateClient()
		carpool := NewCarpool()

		var done sync.WaitGroup

		done.Go(func() { _, _ = carpool.Client(inner, "User:alice").Load("https://example.com/actor") })
		done.Go(func() { _, _ = carpool.Client(inner, "User:bob").Load("https://example.com/actor") })

		synctest.Wait()
		close(inner.gate)
		done.Wait()

		require.Equal(t, int64(2), inner.loads.Load())
	})
}

// TestCarpool_NothingIsKept confirms that a Carpool is not a cache: a Load that starts after the
// previous one finished fetches again.
func TestCarpool_NothingIsKept(t *testing.T) {

	inner := newGateClient()
	close(inner.gate)

	client := NewCarpool().Client(inner, testSigner)

	_, err := client.Load("https://example.com/actor")
	require.NoError(t, err)

	_, err = client.Load("https://example.com/actor")
	require.NoError(t, err)

	require.Equal(t, int64(2), inner.loads.Load())
}

// TestCarpool_ErrorReachesEveryRider confirms that a failed fetch fails every caller sharing it,
// with the original error code intact.
func TestCarpool_ErrorReachesEveryRider(t *testing.T) {

	synctest.Test(t, func(t *testing.T) {

		inner := newGateClient()
		inner.err = derp.NotFound("test", "Actor not found")
		carpool := NewCarpool()

		errs := make([]error, 5)
		var done sync.WaitGroup

		for index := range errs {
			done.Go(func() {
				_, errs[index] = carpool.Client(inner, testSigner).Load("https://example.com/actor")
			})
		}

		synctest.Wait()
		close(inner.gate)
		done.Wait()

		require.Equal(t, int64(1), inner.loads.Load())

		for _, err := range errs {
			require.True(t, derp.IsNotFound(err))
		}
	})
}

// TestCarpool_SoloCallerKeepsItsDocument confirms that a caller with no riders receives the inner
// client's document unchanged, bound to whatever client the inner stack gave it.
func TestCarpool_SoloCallerKeepsItsDocument(t *testing.T) {

	inner := newGateClient()
	close(inner.gate)

	client := NewCarpool().Client(inner, testSigner)
	client.SetRootClient(&mockInnerClient{})

	result, err := client.Load("https://example.com/actor")
	require.NoError(t, err)

	require.Same(t, inner, result.Client())
}

// TestCarpool_RebindWithoutRoot confirms that a shared document reaching a stack whose root was
// never set is bound to the CarpoolClient itself, never to a default HTTP client.
func TestCarpool_RebindWithoutRoot(t *testing.T) {

	client := NewCarpool().Client(&mockInnerClient{}, testSigner)

	result := client.rebind(streams.NewDocument(map[string]any{vocab.PropertyID: "https://example.com/actor"}))

	require.Same(t, client, result.Client())
}

// TestCarpoolClient_Delegates confirms that Save, Delete, and SetRootClient reach the inner client.
func TestCarpoolClient_Delegates(t *testing.T) {

	inner := &mockInnerClient{}
	client := NewCarpool().Client(inner, testSigner)

	document := streams.NewDocument(map[string]any{vocab.PropertyID: "https://example.com/actor"})
	require.NoError(t, client.Save(document))
	require.Equal(t, "https://example.com/actor", inner.savedDocument.ID())

	require.NoError(t, client.Delete("https://example.com/actor"))
	require.Equal(t, "https://example.com/actor", inner.deletedID)

	client.SetRootClient(inner)
	require.True(t, inner.rootClientSet)
}

/******************************************
 * Helpers
 ******************************************/

// testSigner is the signer shared by every stack in a test that does not compare signers
const testSigner = "Application:test"

// gateClient is a streams.Client whose Loads block until its gate is closed, counting every Load
// that reaches it.
type gateClient struct {
	gate    chan struct{}
	err     error
	loads   atomic.Int64
	mutex   sync.Mutex
	options []any
}

// newGateClient returns a gateClient with its gate open for the test to close
func newGateClient() *gateClient {
	return &gateClient{gate: make(chan struct{})}
}

// Load waits for the gate, then returns a fresh document for the URL (or the scripted error)
func (client *gateClient) Load(url string, options ...any) (streams.Document, error) {

	client.loads.Add(1)

	client.mutex.Lock()
	client.options = options
	client.mutex.Unlock()

	<-client.gate

	if client.err != nil {
		return streams.NilDocument(), client.err
	}

	return streams.NewDocument(
		map[string]any{vocab.PropertyID: url, vocab.PropertyName: "Original"},
		streams.WithClient(client),
		streams.WithHTTPHeader(http.Header{"Content-Type": {"application/activity+json"}}),
	), nil
}

// lastOptions returns the options passed to the most recent Load
func (client *gateClient) lastOptions() []any {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return client.options
}

// Save does nothing
func (client *gateClient) Save(streams.Document) error {
	return nil
}

// Delete does nothing
func (client *gateClient) Delete(string) error {
	return nil
}

// SetRootClient does nothing
func (client *gateClient) SetRootClient(streams.Client) {}
