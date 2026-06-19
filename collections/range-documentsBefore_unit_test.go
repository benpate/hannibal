package collections

import (
	"iter"
	"testing"
	"time"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
)

// docWithPublished builds a streams.Document carrying a "published" timestamp.
func docWithPublished(id string, published time.Time) streams.Document {
	return streams.NewDocument(map[string]any{
		vocab.PropertyID:        id,
		vocab.PropertyPublished: published.UTC().Format(time.RFC3339),
	})
}

// sliceIterator yields the given documents as an iter.Seq.
func sliceIterator(documents ...streams.Document) iter.Seq[streams.Document] {
	return func(yield func(streams.Document) bool) {
		for _, document := range documents {
			if !yield(document) {
				return
			}
		}
	}
}

// collectIDs drains an iterator into a slice of IDs.
func collectIDs(iterator iter.Seq[streams.Document]) []string {
	result := make([]string, 0)
	for document := range iterator {
		result = append(result, document.ID())
	}
	return result
}

// TestRangeDocumentsBefore confirms only documents published on or before the
// limit are yielded; later documents are filtered out.
func TestRangeDocumentsBefore(t *testing.T) {

	base := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	limit := base.Unix()

	before := docWithPublished("https://example.com/before", base.Add(-24*time.Hour))
	exactlyAt := docWithPublished("https://example.com/at", base)
	after := docWithPublished("https://example.com/after", base.Add(24*time.Hour))

	input := sliceIterator(before, exactlyAt, after)
	got := collectIDs(RangeDocumentsBefore(input, limit))

	// "before" and "at" (not strictly after) pass; "after" is filtered.
	assert.Contains(t, got, "https://example.com/before")
	assert.Contains(t, got, "https://example.com/at")
	assert.NotContains(t, got, "https://example.com/after")
}

// TestRangeDocumentsBefore_ZeroLimit confirms a zero limit means "no limit" --
// every document is yielded, even far-future ones.
func TestRangeDocumentsBefore_ZeroLimit(t *testing.T) {

	past := docWithPublished("https://example.com/past", time.Now().Add(-48*time.Hour))
	future := docWithPublished("https://example.com/future", time.Now().Add(48*time.Hour))

	got := collectIDs(RangeDocumentsBefore(sliceIterator(past, future), 0))

	assert.Equal(t, []string{"https://example.com/past", "https://example.com/future"}, got)
}

// TestRangeDocumentsBefore_EarlyStop confirms the iterator honors an early break
// from the consumer.
func TestRangeDocumentsBefore_EarlyStop(t *testing.T) {

	old := time.Now().Add(-100 * time.Hour)
	a := docWithPublished("https://example.com/a", old)
	b := docWithPublished("https://example.com/b", old)
	c := docWithPublished("https://example.com/c", old)

	var collected []string
	for document := range RangeDocumentsBefore(sliceIterator(a, b, c), 0) {
		collected = append(collected, document.ID())
		if document.ID() == "https://example.com/b" {
			break
		}
	}

	assert.Equal(t, []string{"https://example.com/a", "https://example.com/b"}, collected)
}

// TestRangeDocumentsBefore_Empty confirms an empty input yields nothing.
func TestRangeDocumentsBefore_Empty(t *testing.T) {
	got := collectIDs(RangeDocumentsBefore(sliceIterator(), 0))
	assert.Empty(t, got)
}
