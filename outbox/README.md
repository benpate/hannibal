# Hannibal / outbox

Outbox mimics an ActivityPub outbox.  Passing a document to the outbox will
use an outbound retry-queue to deliver it to all recipients in `to`, `cc`, 
and `bto` fields.

```go

// Get an Actor's outbox
actor := outbox.New(myActor, outbox.WithClient(myClient), outbox.WithQueue(myQueue))

// The document is the ActivityPub document you're sending
document := map[string]any{
	"@context": vocab.ContextTypeDefault.
	"type": vocab.ActivityTypeCreate,
	"actor": actor.ActorID,
	"object": map[string]any{
		"type": "Note",
		"name": "A new note",	
	},
}

// Send a document via the Actor's outbox
if err := actor.Send(document); err != nil {
	derp.Report(err)
}

```
