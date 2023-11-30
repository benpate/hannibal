# Hannibal / outbox

Outbox mimics an ActivityPub outbox.  Passing a document to the outbox will
use an outbound retry-queue to deliver it to all recipients in `to`, `cc`, 
and `bto` fields.

```go

// Get an Actor's outbox
actor := outbox.New(myActor, outbox.WithClient(myClient), outbox.WithQueue(myQueue))

// Send a document via the Actor's outbox
if err := actor.Send(document); err != nil {
	derp.Report(err)
}

```