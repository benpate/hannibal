# Hannibal / sender

Sender delivers outbound ActivityPub activities from an outbox to each recipient's inbox URL.
It signs every activity with the sending actor's private key and delivers asynchronously through a
[turbine](https://github.com/benpate/turbine) queue, so a slow or failing recipient never blocks the caller.

```go
// A Locator tells the Sender how to find Actors and their followers/recipients
sender := sender.New(myLocator, myQueue)

// Send an activity. The Sender looks up the actor, fans out to every recipient
// in `to`, `cc`, `bto`, and `bcc`, signs each request, and enqueues delivery.
activity := mapof.Any{
	"@context": vocab.ContextTypeActivityStreams,
	"type":     vocab.ActivityTypeCreate,
	"actor":    "https://example.com/@me",
	"object":   mapof.Any{"type": "Note", "content": "Hello, Fediverse"},
	"to":       []string{"https://remote.example/@you"},
}

if err := sender.Send(activity); err != nil {
	derp.Report(err)
}
```

## Interfaces

The Sender depends on two interfaces that you implement for your application:

- **`Locator`** resolves an address into an `Actor` (the signer) and expands a recipient URL into the
  individual inbox URLs that should receive a copy.
- **`Actor`** exposes the actor's ID and its private key, used to sign each outbound request.

## Queue Consumer

`Consumer(sender)` returns a `queue.Consumer` that processes the delivery tasks the Sender enqueues.
It MUST be connected to a live turbine queue so that `SendToAllRecipients` and `SendToSingleRecipient`
are actually executed.

> **Note:** Outbound deliveries refuse to connect to private/loopback IP addresses (SSRF protection
> provided by [remote](https://github.com/benpate/remote)). Production keeps this guard active.
