package streams

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
)

type testClient struct {
	data mapof.Any
}

func newTestClient() testClient {
	return testClient{
		data: testStreamData(),
	}
}

func (client testClient) LoadDocument(uri string, defaultValue map[string]any) (Document, error) {

	if value, ok := client.data[uri]; ok {
		return NewDocument(value, WithClient(client)), nil
	}

	return NilDocument(), derp.NewInternalError("hannibal.streams.testClient.Load", "Unknown URI", uri)
}

func (client testClient) LoadActor(uri string) (Document, error) {
	return NilDocument(), derp.NewInternalError("hannibal.streams.testClient.Load", "Not implemented", uri)
}

// testStreamData returns a collection of documents for the test client to return
func testStreamData() mapof.Any {

	rawData := mapof.String{
		"https://demo/collection": `{
			"@context":"https://w3.org/ns/activitystreams",
			"@id":"https://demo/collection",
			"@type":"Collection",
			"totalItems":3,
			"orderedItems":[
				"https://example/collection-url-1",
				"https://example/collection-url-2",
				"https://example/collection-url-3"
			]
		}`,
		"https://demo/ordered": `{
			"@context":"https://w3.org/ns/activitystreams",
			"type":"OrderedCollection",
			"totalItems":3,
			"orderedItems":[
				"https://example/url-1",
				"https://example/url-2",
				"https://example/url-3"
			]
		}`,
		"https://demo/ordered-page": `{
			"@context":"https://w3.org/ns/activitystreams",
			"type":"OrderedCollection",
			"totalItems":9,
			"first":"https://demo/ordered-page-1"
		}`,
		"https://demo/ordered-page-1": `{
			"@context":"https://w3.org/ns/activitystreams",
			"type":"OrderedCollectionPage",
			"totalItems":9,
			"next":"https://demo/ordered-page-2",
			"orderedItems":[
				"https://example/url-1",
				"https://example/url-2",
				"https://example/url-3"
			]
		}`,
		"https://demo/ordered-page-2": `{
			"@context":"https://w3.org/ns/activitystreams",
			"type":"OrderedCollectionPage",
			"totalItems":9,
			"next":"https://demo/ordered-page-3",
			"orderedItems":[
				"https://example/url-4",
				"https://example/url-5",
				"https://example/url-6"
			]
		}`,
		"https://demo/ordered-page-3": `{
			"@context":"https://w3.org/ns/activitystreams",
			"type":"OrderedCollectionPage",
			"totalItems":9,
			"orderedItems":[
				"https://example/url-7",
				"https://example/url-8",
				"https://example/url-9"
			]
		}`,
		"https://demo/interminus": `{
			"@context":"https://w3.org/ns/activitystreams",
			"type":"OrderedCollection",
			"totalItems":9,
			"first":"https://demo/interminus-1"
		}`,
		"https://demo/interminus-1": `{
			"@context":"https://w3.org/ns/activitystreams",
			"type":"OrderedCollectionPage",
			"totalItems":9,
			"next":"https://demo/interminus-2",
			"orderedItems":[
				"https://example/url-1",
				"https://example/url-2",
				"https://example/url-3"
			]
		}`,
	}

	result := mapof.NewAny()
	for key, value := range rawData {

		item := mapof.NewAny()
		if err := json.Unmarshal([]byte(value), &item); err != nil {
			panic(err)
		}

		result[key] = item
	}

	return result
}
