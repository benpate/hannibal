# Hannibal / validator

Validator provides pluggable checks that confirm an inbound ActivityPub activity is authentic before
the [router](../router) hands it to your application. The router runs a chain of validators; each one
returns a `Result` of `ResultValid`, `ResultInvalid`, or `ResultUnknown`.

```go
// Build a chain of validators (this is what router uses by default)
validators := []router.Validator{
	validator.NewHTTPSig(nil),       // verify the HTTP Signature
	validator.NewDeletedObject(),    // confirm Delete activities really are deleted
}
```

## Results

A validator returns one of three results, and the router walks the chain accordingly:

- **`ResultValid`** — the activity is authentic; stop and accept it.
- **`ResultInvalid`** — the activity is forged or malformed; stop and reject it.
- **`ResultUnknown`** — this validator can't decide; continue to the next validator. If every validator
  returns `ResultUnknown`, validation fails closed.

## Included Validators

- **`HTTPSig`** verifies the request's HTTP Signature against the actor's public key.
- **`MatchActor`** confirms the activity's actor matches an expected actor ID.
- **`DeletedObject`** confirms a `Delete` activity refers to an object that is actually gone.
- **`HTTPLookup`** confirms an activity exists by fetching it from its origin server.
- **`None`** performs no validation (always `ResultUnknown`); useful as a placeholder in tests.
