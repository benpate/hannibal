package collections

import (
	"iter"
	"math"
	"time"

	"github.com/benpate/hannibal/streams"
)

func RangeDocumentsBefore(iterator iter.Seq[streams.Document], limit int64) iter.Seq[streams.Document] {

	return func(yield func(streams.Document) bool) {

		// Empty limit actually means "the end of time"
		if limit == 0 {
			limit = math.MaxInt64
		}

		// Convert to a time.Time object for comparison
		limitTime := time.Unix(limit, 0)

		// Scan through the range of documents
		for document := range iterator {

			// Filter out documents that are equal or after the "before" date
			if document.Published().After(limitTime) {
				continue
			}

			if !yield(document) {
				return
			}
		}
	}
}
