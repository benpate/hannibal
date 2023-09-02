package sigs

import (
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#section-2.1
type Signature struct {
	KeyID     string   // ID (URL) of the key used to create this signature
	Algorithm string   // Algorithm used to create this signature (should be ignored per IEFT spec)
	Headers   []string // List of headers that were signed
	Signature []byte   // Base64 encoded signature
	Created   int64    // Unix epoch (in seconds) when this signature was created
	Expires   int64    // Unix epoch (in seconds) when this signature expires
}

// NewSignature returns a fully initialized Signature object
func NewSignature() Signature {
	return Signature{
		Headers:   make([]string, 0),
		Signature: make([]byte, 0),
	}
}

// IsExpired returns TRUE if the current date is
// less than its expiration date, OR if the
// createDate + duration is less than the current date.
// Calculations are skipped if the duration, created,
// or expires values are zero.
func (signature Signature) IsExpired(duration int64) bool {

	now := time.Now().Unix()

	// If the "expires" value is valid and in the past, then the signature is expired
	if signature.Expires > 0 {
		if signature.Expires < now {
			return true
		}
	}

	// If the "created" and "duration" values are valid,
	// and their sum is in the past, then the signature is expired
	if (signature.Created > 0) && (duration > 0) {
		if (signature.Created + duration) < now {
			return true
		}
	}

	// Otherwise, the signature is not expired
	return false
}

// ParseSignature parses a string into an HTTP Signature
func ParseSignature(value string) (Signature, error) {

	result := NewSignature()

	// Split the signature into a list of key=value pairs
	items := strings.Split(value, ",")

	for _, item := range items {
		item = strings.TrimSpace(item)       // remove extra whitespace
		name, value := list.Split(item, '=') // split into key=value
		value = strings.Trim(value, `"`)     // remove quotes from value

		// Assemble key/value pairs into the Signature structure
		switch name {

		case "keyId":
			result.KeyID = value

		case "algorithm":
			result.Algorithm = value

		case "headers":
			result.Headers = strings.Split(value, " ")

		case "signature":
			if value, err := base64.StdEncoding.DecodeString(value); err != nil {
				return Signature{}, derp.Wrap(err, "sigs.ParseSignature", "Error decoding signature", value)
			} else {
				result.Signature = value
			}

		case "created":
			result.Created, _ = strconv.ParseInt(value, 10, 64)

		case "expires":
			result.Expires, _ = strconv.ParseInt(value, 10, 64)
		}
	}

	// RULE: Required Fields
	if result.KeyID == "" {
		return Signature{}, derp.New(derp.CodeBadRequestError, "sigs.ParseSignature", "Field 'keyId' is required.")
	}

	if len(result.Headers) == 0 {
		return Signature{}, derp.New(derp.CodeBadRequestError, "sigs.ParseSignature", "Field 'headers' is required.")
	}

	if len(result.Signature) == 0 {
		return Signature{}, derp.New(derp.CodeBadRequestError, "sigs.ParseSignature", "Field 'signature' is required.")
	}

	return result, nil
}

// String returns the Signature as a string
func (signature Signature) String() string {

	var buffer strings.Builder

	buffer.WriteString(`keyId="`)
	buffer.WriteString(signature.KeyID)
	buffer.WriteString(`"`)

	buffer.WriteString(`,algorithm="`)
	buffer.WriteString(signature.Algorithm)
	buffer.WriteString(`"`)

	buffer.WriteString(`,headers="`)
	buffer.WriteString(strings.Join(signature.Headers, " "))
	buffer.WriteString(`"`)

	buffer.WriteString(`,signature="`)
	buffer.WriteString(signature.Base64())
	buffer.WriteString(`"`)

	if signature.Created > 0 {
		buffer.WriteString(`,created="`)
		buffer.WriteString(strconv.FormatInt(signature.Created, 10))
		buffer.WriteString(`"`)
	}

	if signature.Expires > 0 {
		buffer.WriteString(`,expires="`)
		buffer.WriteString(strconv.FormatInt(signature.Expires, 10))
		buffer.WriteString(`"`)
	}

	return buffer.String()
}

// AlgorithmPrefix returns the first part of the algorithm name, such as "rsa", "hmac", or "ecdsa"
func (signature Signature) AlgorithmPrefix() string {
	return list.Head(signature.Algorithm, '-')
}

// SignatureBytes returns the signature as a slice of bytes
func (signature Signature) Base64() string {
	return base64.StdEncoding.EncodeToString(signature.Signature)
}
