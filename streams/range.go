package streams

import "iter"

func Range(document Document) iter.Seq[Document] {

	return func(yield func(Document) bool) {
		for ; document.NotNil(); document = document.Tail() {
			if !yield(document.Head()) {
				return
			}
		}
	}
}
