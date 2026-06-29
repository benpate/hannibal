# Hannibal / outbox

Outbox mimics an ActivityPub outbox. You create an `Actor` from an ID and a private key, then call one of its `Send*` methods to deliver an activity. The Actor signs each outbound request with its private key and delivers it to every recipient addressed in the activity's `to`, `cc`, `bto`, and `bcc` fields.

```go
// Create an Actor from its ID and private key, plus optional settings
actor := outbox.NewActor(
	"https://example.com/@me",
	privateKey,
	outbox.WithPublicKey("https://example.com/@me#main-key"),
	outbox.WithClient(myClient),
	outbox.WithFollowers(myFollowersIterator),
)

// Send a Create activity for a new Note. The Actor signs the request and
// delivers it to the recipients addressed in the document.
note := streams.NewDocument(mapof.Any{
	"@context": vocab.ContextTypeActivityStreams,
	"type":     vocab.ObjectTypeNote,
	"attributedTo": actor.ActorID(),
	"content":  "A new note",
	"to":       vocab.NamespaceActivityStreamsPublic,
})

actor.SendCreate(note)
```

## Sending Activities

The typed helpers build the wrapping activity for you: `SendCreate`, `SendUpdate`, `SendDelete`, `SendFollow`, `SendAccept`, `SendLike`, `SendDislike`, `SendAnnounce`, and `SendUndo`. For anything they don't cover, `Send(message, recipients...)` delivers a raw activity, and `SendOne(recipientID, message)` delivers to a single recipient.

## Options

`NewActor` takes the actor ID and private key as required arguments, plus optional `ActorOption` settings:

- `WithPublicKey(id)` — the public-key ID advertised in signatures.
- `WithClient(client)` — the `streams.Client` used to resolve recipients (defaults to a standard client).
- `WithFollowers(iterator)` — an iterator over the actor's followers, used to expand the special "followers" recipient.

> **Note:** Outbound deliveries refuse to connect to private/loopback IP addresses (SSRF protection inherited from [remote](https://github.com/benpate/remote)). Production keeps this guard active; only tests that deliver to a loopback server opt out.
