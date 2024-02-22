# Hannibal / inbox

The Inbox library gives you: 1) an ActivityPub inbox handler that parses and validates incoming ActivityPub messages, and 2) a router that identifies messages by their type and object, and routes them to the correct business logic in your application

``` golang


// Set up handlers for different kinds of activities and documents
activityHandler := inbox.NewRouter[CustomContextType]()

// Here's a handler to accept Create/Note messages
activityHandler.Add(vocab.ActivityTypeCreate, vocab.ObjectTypeNote, func(context CustomContextType, activity streams.Document) error {
	// do something with the activity
})

// You can do wildcards too.  Here's a handler to accept 
// Follow/Any messages
activityHandler.Add(vocab.ActivityTypeFollow, vocab.Any, func(context CustomContextType, activity streams.Document) error {
	// do something with the follow request.
	// remember to send an "Accept" message back to the sender
})

// Here's a catch-all handler that receives any uncaught messages
activityHandler.Add(vocab.Any, vocab.Any, func(context CustomContextType, activity streams.Document) error {
	// do something with this activity
})

// Add routes to your web server
myAppRouter.POST("/my/inbox",func (r *http.Request, w *http.Response) {

	// Parse and validate the posted activity
	activity, err := pub.ReceiveRequest(r)
	
	// Handle errors however you like
	if err != nil {
		...
	}

	context := // create custom "Context" value for this request
	
	// Pass the activity to the activityHandler, that will figure
	// out what kind of activity/object we have and pass it to the 
	// previously registered handler function
	if err := activityHandler.Handle(context, actitity); err != nil {
		// do something with the error
	}
	
	// Success!
}

```