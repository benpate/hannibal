# Hannibal / clients

This package provides `streams.Client` implementations that load ActivityStreams documents from remote
servers. Each client wraps an inner client, so you can stack them (caching, lookup, transport) and pass
the result wherever a `streams.Client` is expected.

## HashLookup

`HashLookup` resolves URLs that contain a `#fragment`. Many ActivityStreams objects (public keys,
attachments, tags) are not published at their own URL but are embedded inside a parent document and
identified by a fragment. `HashLookup` loads the base document, then searches its top-level properties
for an object whose `id` matches the full URL.

```go
// Wrap any inner client with HashLookup
client := clients.NewHashLookup(myInnerClient)

// A plain URL passes straight through to the inner client
actor, _ := client.Load("https://example.com/@me")

// A URL with a fragment loads the base document, then returns the embedded
// object whose "id" equals the full URL (e.g. the actor's public key)
key, _ := client.Load("https://example.com/@me#main-key")
```

`Save`, `Delete`, and `SetRootClient` delegate to the inner client unchanged.

## Carpool

`Carpool` merges concurrent Loads of the same URL by the same signer. When many callers ask for one document at the same moment, the first caller runs the Load and the others wait and receive its result, so a burst of 200 requests for one actor becomes one fetch. It is not a cache: once the Load finishes, nothing is kept, and the next caller loads again.

Create one `Carpool` per process and wrap each client stack with it. Stacks built from the same `Carpool` share their in-progress Loads.

```go
// Once, at startup
carpool := clients.NewCarpool()

// Every time a client stack is built, naming the actor that stack signs as
client := carpool.Client(myInnerClient, "Application:example.com")
```

The signer is part of the grouping key, because a remote server may answer each signer differently: one signer can be blocked where another is not. Two stacks share a Load only when they would have signed it the same way. Any string works as a signer as long as it identifies the signing actor and contains no NUL character.

Every caller that shares a Load receives its own deep copy of the document, bound to its own client stack, so follow-up loads are signed and filtered as that caller. A Load that passes any options runs alone, because the Carpool cannot tell whether two sets of options ask for the same result. A Carpool covers one process; servers behind a load balancer each have their own.

`Save`, `Delete`, and `SetRootClient` delegate to the inner client.
