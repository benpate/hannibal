# sigs

Common-Sense HTTP Signatures

<img src="https://github.com/benpate/hannibal/raw/main/meta/sigs.jpg" style="width:100%; display:block; margin-bottom:20px;" alt="Oil painting titled: Signers of the Constitution, by Thomas Pritchard Rossiter (1817-1871)"/>


[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/hannibal/sigs.svg)](https://pkg.go.dev/github.com/benpate/hannibal/sigs)
[![Build Status](https://img.shields.io/github/actions/workflow/status/EmissarySocial/emissary/go.yml?branch=main)](https://github.com/EmissarySocial/emissary/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/hannibal?style=flat-square)](https://goreportcard.com/report/github.com/benpate/hannibal)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/hannibal/sigs.svg?style=flat-square)](https://codecov.io/gh/benpate/hannibal/tree/main/sigs)

This library is a simple-yet-thorough implementation the IETF HTTP Signatures specification.  It aims to be extensively tested and documented, with extensions for you to test and troubleshoot your own implementations.

## Project Status (DO NOT USE)

This code is still being developed and is not ready to use.  

## Key Features

* Simple API with sensible defaults
* Complete implementation of all commonly used cipher and digest algorithms
* Extensive use of [zerolog](https://github.com/rs/zerolog) logging to simplify troubleshooting

## Signing Outbound Requests

The `sigs` library makes it easy to sign an outbound http.Request.  It includes sensible defaults (shown below) so that most uses should "just work" with minimal configuration.

```go
// Generate a private key to sign. Real code would likely
// retrieve this from a database.
privateKey, err := rsa.GeneratePrivateKey(rand.Reader, 2048)

// Generate an http.Request to sign. Real code would likely
// set additional values in the outbound request.
request := http.NewRequest("POST", "https://example.com", nil)

// Sign the Request with the Private Key. Yes, that's it.
err = sigs.Sign(request, privateKey)

```

### Signing Options

In the event that you need to customize the way you sign a Request, you can pass one or more optional functions into the `Sign` function.

| Option | Description | Default |
|--------|-------------|---------|
| SignDigests | Sets the algorithm(s) to use when creating the "Digest" header. | `sha-256` |
| SignFields | Sets the field(s) to use when creating the "Signature" header . | `(request-target) host date digest` |

```go
// How to sign a request using additional options
err := sigs.Sign(
	request, 
	privateKey,
	sigs.SignatureFields("(request-target)", "(created)", "(expires)"),
	sigs.SignatureDigest("sha-512"),
)
```

### Object Notation

In most cases, the above syntax is the simplest way to use `sigs`.  However, the library also publishes the underlying objects used to sign http.Requests, which you can also access directly.  For instance, you may need to do this if you need to use complex logic to determine what options to set.

```go
signer := sigs.NewSigner()
signer.Use(SignFields("content-type", "date"))

if err := signer.Sign(request, privateKey); err != nil {
	// handle error
}

```

## Verifying Inbound Requests

```go
// Define an http.Request.  Real code would receive this request
// from a remote server in an http.Handler function
request := http.NewRequest("POST", https://example.com", nil)

// Define the PEM-encoded certificate.  Real code would
// retrieve this from the remote user's profile.
publicKeyPEM := `-----BEGIN PRIVATE KEY----- ... -----END PRIVATE KEY-----`

// Verify the request has a valid signature from the certificate.
// Yes, that's it.
err := sigs.Verify(request, publicKeyPEM)
```

### Verification Options

In the event that you need to customize the way you verify a Request, you can pass one or more optional functions into the `Verify` function.

| Option | Description | Default |
|--------|-------------|---------|
| VerifyDigests | Identifies one ore more digest algorithms to use when verifying the "Digest" header. AT LEAST ONE algorithm must match for verification to succeed. If present, additional algorithms will be ignored. | `sha-256` |
| VerifyFields | Identifies a list of fields that MUST ALL be present in the signature.  If present, additional fields will be ignored. | `(request-target) host date digest` |

```go
// How to verify a request using additional options
err := sigs.Verify(
	request, 
	publicKeyPEM,
	sigs.VerifyFields("(request-target)", "(created)", "(expires)"),
	sigs.VerifyDigest("sha-512"),
)
```

### Object Notation
In most cases, the above syntax is the simplest way to use `sigs`.  However, the library also publishes the underlying objects used to verfy http.Requests, which you can also access directly.  For instance, you may need to do this if you need to use complex logic to determine what options to set.

```go
verifier := sigs.NewVerifier()
verifier.Use(VerifyFields("content-type","date"))

if err := verifier.Verify(request, publicKeyPEM); err != nil {
	// handle error
}

```

## Troubleshooting

The `sigs` library generates fine grained debugging information with zerolog structured logging library.  By default, it sets the logging level to `Disabled` so that no logging information is written.  If you need to see deeper into `sigs` add the following into your application code:

```go
func main() {

	// zerolog.SetGlobalLevel(zerolog.Trace) // for a step-by-step trace of every sigs action.
	// zerolog.SetGlobalLevel(zerolog.Debug) // for higher-level debugging of signatures and verification
	
	// your code here...
}

```


## Supported Algorithms

### Signatures

* hs2019
* rsa-sha256
* rsa-shaA512
* hmac-sha256 (In Progress)
* ecdsa-sha256 (In Progress)

### Digests

* sha256
* sha512
* unixsum (In Progress)
* unixcksum (In Progress)
* crc32c  (In Progress)
* id-sha-256 (In Progress) 
* id-sha-512 (In Progress)

## References

* IETF Standard: https://datatracker.ietf.org/doc/html/draft-ietf-httpbis-message-signatures
* Mastodon Security Documentation: https://docs.joinmastodon.org/spec/security/