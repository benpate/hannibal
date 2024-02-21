package inbox

import (
	"encoding/json"
	"fmt"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

// Router is a simple object that routes incoming ActivityPub activities to the appropriate handler
type Router[T any] struct {
	routes map[string]RouteHandler[T]
}

// RouteHandler is a function that handles a specific type of ActivityPub activity.
// RouteHandlers are registered with the Router object along with the names of the activity
// types that they correspond to.
type RouteHandler[T any] func(context T, activity streams.Document) error

// NewRouter creates a new Router object
func NewRouter[T any]() Router[T] {
	return Router[T]{
		routes: make(map[string]RouteHandler[T]),
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
func (router *Router[T]) Add(activityType string, objectType string, routeHandler RouteHandler[T]) {
	router.routes[activityType+"/"+objectType] = routeHandler
}

// Handle takes an ActivityPub activity and routes it to the appropriate handler
func (router *Router[T]) Handle(context T, activity streams.Document) error {

	activityType := activity.Type()
	objectType := activity.Object().Type()

	if packageDebugLevel >= DebugLevelTerse {
		if packageDebugLevel >= DebugLevelVerbose {
			fmt.Println("------------------------------------------")
		}
		fmt.Println("HANNIBAL: Received Message: " + activityType + "/" + objectType)
		if packageDebugLevel >= DebugLevelVerbose {
			marshalled, _ := json.MarshalIndent(activity.Value(), "", "  ")
			fmt.Println(string(marshalled))
		}
	}

	if routeHandler, ok := router.routes[activityType+"/"+objectType]; ok {
		if packageDebugLevel >= DebugLevelVerbose {
			fmt.Println("HANNIBAL: Found Route: " + activityType + "/" + objectType)
		}
		return routeHandler(context, activity)
	}

	if routeHandler, ok := router.routes[activityType+"/*"]; ok {
		if packageDebugLevel >= DebugLevelVerbose {
			fmt.Println("HANNIBAL: Found Route: " + activityType + "/*")
		}
		return routeHandler(context, activity)
	}

	if routeHandler, ok := router.routes["*/"+objectType]; ok {
		if packageDebugLevel >= DebugLevelVerbose {
			fmt.Println("HANNIBAL: Found Route: " + "*/" + objectType)
		}
		return routeHandler(context, activity)
	}

	if routeHandler, ok := router.routes["*/*"]; ok {
		if packageDebugLevel >= DebugLevelVerbose {
			fmt.Println("HANNIBAL: Found Route: */*")
		}
		return routeHandler(context, activity)
	}

	return derp.NewBadRequestError("hannibal.pub.Router.Handle", "No route found for activity", activity.Value())
}
