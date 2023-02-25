# pub

This is the topmost library of Hannibal, providing a simple interface for sending and receiving ActivityPub messages.

## How to Send Messages

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

``` golang
// coming soon..

```