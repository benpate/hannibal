package router

import "github.com/benpate/hannibal/streams"

// RouteHandler is a function that handles a specific type of ActivityPub activity.
// RouteHandlers are registered with the Router object along with the names of the activity
// types that they correspond to.
type RouteHandler[T any] func(context T, activity streams.Document) error
