package streams

import (
	"strings"
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// benchmarkSink keeps benchmark results alive so the compiler cannot discard the work.
var benchmarkSink string

// benchmarkCreateNote returns a Mastodon-shaped Create/Note of roughly 3.5 KB.
func benchmarkCreateNote() mapof.Any {

	paragraph := "<p>Lorem ipsum dolor sit amet, <a href=\"https://example.com/tags/go\" class=\"mention hashtag\" rel=\"tag\">#<span>go</span></a> consectetur adipiscing elit.</p>"

	return mapof.Any{
		vocab.AtContext:     vocab.NamespaceActivityStreams,
		vocab.PropertyID:    "https://example.com/users/alice/statuses/1/activity",
		vocab.PropertyType:  vocab.ActivityTypeCreate,
		vocab.PropertyActor: "https://example.com/users/alice",
		vocab.PropertyTo:    []any{vocab.NamespaceActivityStreamsPublic},
		vocab.PropertyCC:    []any{"https://example.com/users/alice/followers"},
		vocab.PropertyObject: mapof.Any{
			vocab.PropertyID:           "https://example.com/users/alice/statuses/1",
			vocab.PropertyType:         vocab.ObjectTypeNote,
			vocab.PropertyAttributedTo: "https://example.com/users/alice",
			vocab.PropertySummary:      "Content warning: <b>benchmarks</b>",
			vocab.PropertyContent:      strings.Repeat(paragraph, 20),
			vocab.PropertyURL:          "https://example.com/@alice/1",
			vocab.PropertyMediaType:    "text/html",
		},
	}
}

// BenchmarkDocumentID measures reading an identifier through the sanitizing accessor.
func BenchmarkDocumentID(b *testing.B) {
	for document := NewDocument(benchmarkCreateNote()); b.Loop(); {
		benchmarkSink = document.ID()
	}
}

// BenchmarkDocumentContent measures reading HTML content through the UGC policy.
func BenchmarkDocumentContent(b *testing.B) {
	for object := NewDocument(benchmarkCreateNote()).Object(); b.Loop(); {
		benchmarkSink = object.Content()
	}
}

// BenchmarkDocumentAccessorSet measures the accessors an inbound Create typically reads.
func BenchmarkDocumentAccessorSet(b *testing.B) {
	for document := NewDocument(benchmarkCreateNote()); b.Loop(); {
		object := document.Object()
		benchmarkSink = document.ID()
		benchmarkSink = document.Type()
		benchmarkSink = document.ActorID()
		benchmarkSink = object.ID()
		benchmarkSink = object.Type()
		benchmarkSink = object.AttributedTo().ID()
		benchmarkSink = object.Summary()
		benchmarkSink = object.Content()
		benchmarkSink = object.URL()
		benchmarkSink = object.MediaType()
	}
}
