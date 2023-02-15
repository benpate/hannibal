package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

// Router is a simple object that routes incoming ActivityPub activities to the appropriate handler
type Router struct {
	routes map[string]RouteHandler
}

// RouteHandler is a function that handles a specific type of ActivityPub activity.
// RouteHandlers are registered with the Router object along with the names of the activity
// types that they correspond to.
type RouteHandler func(activity streams.Document) error

// NewRouter creates a new Router object
func NewRouter() Router {
	return Router{
		routes: make(map[string]RouteHandler),
	}
}

// Add puts a new route to the router.  You can use "*" as a wildcard for
// either the activityType or objectType. The Handler method tries to match
// handlers from most specific to least specific.
// activity/object
// activity/*
// */object
// */*
//
// For performance reasons, this function is not thread-safe.
// So, you should add all routes before starting the server, for
// instance, in your app's `init` functions.
func (router *Router) Add(activityType string, objectType string, routeHandler RouteHandler) {
	router.routes[activityType+"/"+objectType] = routeHandler
}

// Handle takes an ActivityPub activity and routes it to the appropriate handler
func (router *Router) Handle(activity streams.Document) error {

	activityType := activity.Type()
	objectType := activity.Object().Type()

	if routeHandler, ok := router.routes[activityType+"/"+objectType]; ok {
		return routeHandler(activity)
	}

	if routeHandler, ok := router.routes[activityType+"/*"]; ok {
		return routeHandler(activity)
	}

	if routeHandler, ok := router.routes["*/"+objectType]; ok {
		return routeHandler(activity)
	}

	if routeHandler, ok := router.routes["*/*"]; ok {
		return routeHandler(activity)
	}

	return derp.NewBadRequestError("pub.Router.Handle", "No route found for activity", activity.Value())
}
