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
