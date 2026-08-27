package collections

// Option is a functional option that configures a collection traversal.
type Option func(*RangeConfig)

// WithMaxPages returns an Option that overrides the maximum number of pages that a
// traversal will follow. A value less than one yields no pages at all.
func WithMaxPages(maxPages int) Option {
	return func(config *RangeConfig) {
		config.maxPages = maxPages
	}
}

// WithMaxDocuments returns an Option that overrides the maximum number of documents
// that RangeDocuments will yield. A value less than one yields no documents at all.
func WithMaxDocuments(maxDocuments int) Option {
	return func(config *RangeConfig) {
		config.maxDocuments = maxDocuments
	}
}
