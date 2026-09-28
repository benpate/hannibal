package streams

import "github.com/microcosm-cc/bluemonday"

// strictPolicy is the shared policy that strips all HTML. Callers must never modify it.
var strictPolicy = bluemonday.StrictPolicy()

// ugcPolicy is the shared policy that allows user-generated HTML. Callers must never modify it.
var ugcPolicy = bluemonday.UGCPolicy()
