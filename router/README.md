# Hannibal / router

The router package gives you two things: 1) `ReceiveRequest`, which parses and validates an incoming ActivityPub HTTP request, and 2) a `Router` that identifies an activity by its type and object type, then dispatches it to the matching handler in your application.

```go
// Create a Router whose handlers receive your own context type
activityRouter := router.New[CustomContextType]()

// Handle Create/Note messages
activityRouter.Add(vocab.ActivityTypeCreate, vocab.ObjectTypeNote, func(context CustomContextType, activity streams.Document) error {
	// do something with the activity
	return nil
})

// Wildcards work too. Handle Follow/Any messages
activityRouter.Add(vocab.ActivityTypeFollow, vocab.Any, func(context CustomContextType, activity streams.Document) error {
	// handle the follow request (remember to send an "Accept" back)
	return nil
})

// A catch-all handler for anything not matched above
activityRouter.Add(vocab.Any, vocab.Any, func(context CustomContextType, activity streams.Document) error {
	// do something with this activity
	return nil
})

// Wire it into your web server
myAppRouter.POST("/my/inbox", func(w http.ResponseWriter, r *http.Request) {

	// Parse and validate the inbound activity (verifies the HTTP Signature, etc.)
	activity, err := router.ReceiveRequest(r, myClient)
	if err != nil {
		// handle the error however you like
		return
	}

	// Build whatever context value this request needs
	context := makeContext(r)

	// Dispatch to the handler registered for this activity/object type
	if err := activityRouter.Handle(context, activity); err != nil {
		// handle the error
	}
})
```

## Matching

`Add(activityType, objectType, handler)` registers a handler for a specific activity type and object type. `Handle` looks for the most specific match first, falling back to wildcards (`vocab.Any`) for the object type, the activity type, or both — so a single `(vocab.Any, vocab.Any)` handler acts as a catch-all.

## Receiving Requests

`ReceiveRequest(request, client, options...)` reads the request body, parses it into a `streams.Document`, and runs the validator chain before returning. Options let you tune it:

- `WithValidators(...)` — replace the validator chain (defaults to HTTP Signature verification). See [validator](../validator/) for the available checks.
- `WithPublicKeyFinder(...)` — supply the key finder used to verify signatures.
- `WithMaxBodySize(bytes)` — cap the request body size.
