package inbox

import (
	"encoding/json"
	"fmt"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/property"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/rs/zerolog/log"
)

// Router is a simple object that routes incoming ActivityPub activities to the appropriate handler
type Router[T any] struct {
	routes map[string]RouteHandler[T]
	config Config
}

// RouteHandler is a function that handles a specific type of ActivityPub activity.
// RouteHandlers are registered with the Router object along with the names of the activity
// types that they correspond to.
type RouteHandler[T any] func(context T, activity streams.Document) error

// NewRouter creates a new Router object
func NewRouter[T any](options ...Option) Router[T] {
	result := Router[T]{
		routes: make(map[string]RouteHandler[T]),
		config: NewConfig(),
	}

	for _, option := range options {
		option(&result.config)
	}

	return result
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

	const location = "hannibal.inbox.Router.Handle"

	activityType := activity.Type()

	// If this is a Document (not an Activity) then wrap it in
	// an implicit "Create" activity before routing.
	if vocab.ValidateActivityType(activityType) == vocab.Unknown {

		newValue := property.Map{
			vocab.AtContext:      activity.AtContext(),
			vocab.PropertyID:     activity.ID(),
			vocab.PropertyActor:  activity.Actor(),
			vocab.PropertyType:   vocab.ActivityTypeCreate,
			vocab.PropertyObject: activity.Value(),
		}

		activity.SetValue(newValue)
		activityType = vocab.ActivityTypeCreate
	}

	objectType := activity.Object().Type()

	if canDebug() {
		log.Debug().Str("type", activityType+"/"+objectType).Msg("Hannibal Router: Received Message")

		if canTrace() {
			marshalled, _ := json.MarshalIndent(activity.Value(), "", "  ")
			fmt.Println(string(marshalled))
		}
	}

	if routeHandler, ok := router.routes[activityType+"/"+objectType]; ok {
		log.Trace().Str("type", activityType+"/"+objectType).Msg("Hannibal Router: route matched.")
		return routeHandler(context, activity)
	}

	if routeHandler, ok := router.routes[activityType+"/"+vocab.Any]; ok {
		log.Trace().Str("type", activityType+"/*").Msg("Hannibal Router: route matched.")
		return routeHandler(context, activity)
	}

	if routeHandler, ok := router.routes[vocab.Any+"/"+objectType]; ok {
		log.Trace().Str("type", "*/"+objectType).Msg("Hannibal Router: route matched.")
		return routeHandler(context, activity)
	}

	if routeHandler, ok := router.routes[vocab.Any+"/"+vocab.Any]; ok {
		log.Trace().Str("type", "*/*").Msg("Hannibal Router: route matched.")
		return routeHandler(context, activity)
	}

	return derp.NewBadRequestError(location, "No route found for activity", activityType, objectType, activity.Value())
}
