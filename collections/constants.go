package collections

// defaultMaxPages is the default bound on how many pages a traversal will follow
// (override per-call with WithMaxPages).  The "next" chain is remote-controlled data,
// so a cycle of non-empty pages (or an endlessly generated chain) would otherwise
// iterate forever.  The empty-page check in RangePages only catches WriteFreely-style
// loops of EMPTY pages; this cap is the backstop for everything else.  1024 pages is
// far beyond any legitimate collection a crawler should walk in one pass.
const defaultMaxPages = 1024

// defaultMaxDocuments is the default bound on how many documents RangeDocuments will
// yield across all pages (override per-call with WithMaxDocuments).  It complements
// defaultMaxPages: the page cap alone still admits a hostile server stuffing each
// page with an enormous "items" array, so the document cap bounds the total work a
// traversal can hand its caller.
const defaultMaxDocuments = 10_000
