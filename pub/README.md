# pub

This is the topmost library of Hannibal, providing a simple interface for sending and receiving ActivityPub messages.

## How to Send Messages

Hannibal includes a number of convenience functions for sending activities of various types.

```golang
// You'll need to populate a pub.Actor struct, which contains
// an individual's public IRI and cryptographic keys
actor := LoadMyActor()

// The document is the ActivityPub document you're sending
document := map[string]any{
	"@context": vocab.ContextTypeDefault.
	"type": vocab.ActivityTypeCreate,
	"actor": actor.ActorID,
	"object": map[string]any{
	
	},
}

// The target is the profile URL of the message recipient
target := "https://emissary.social/@username"

if err := pub.SendActivity(actor, document, target); err != nil {
	// handle errors
}

// Yes.  That's all there is.
```


## How to Receive Messages

Hannibal includes an ActivityPub receiver that parses and validates messages, and a router that makes it easy to pass new messages back into your application

``` golang


// Set up handlers for different kinds of activities and documents
activityHandler := pub.NewRouter()

// Here's a handler to accept Create/Note messages
activityHandler.Add(vocab.ActivityTypeCreate, vocab.ObjectTypeNote, func(activity streams.Document) error {
	// do something with the activity
})

// You can do wildcards too.  Here's a handler to accept 
// Follow/Any messages
activityHandler.Add(vocab.ActivityTypeFollow, vocab.Any, func(activity streams.Document) error {
	// do something with the follow request.
	// remember to send an "Accept" message back to the sender
})

// Here's a catch-all handler that receives any uncaught messages
activityHandler.Add(vocab.Any, vocab.Any, func(activity streams.Document) error {
	// do something with this activity
})

// Add routes to your web server
myRouter.POST("/my/inbox",func (r *http.Request, w *http.Response) {

	// Parse and validate the posted activity
	activity, err := pub.ReceiveInboxRequest(r)
	
	// Handle errors however you like
	if err != nil {
		...
	}
	
	// Pass the activity to the activityHandler, that will figure
	// out what kind of activity/object we have and pass it to the 
	// previously registered handler function
	if err := activityHandler.Handle(actitity); err != nil {
		// do something with the error
	}
	
	// Success!
}

```