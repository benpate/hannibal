package sigs

import (
	"strconv"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
)

// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#section-2.1
type Signature struct {
	KeyID     string   // ID (URL) of the key used to create this signature
	Algorithm string   // Algorithm used to create this signature
	Headers   []string // List of headers that were signed
	Signature string   // Base64 encoded signature
	Created   int64    // Unix epoch (in seconds) when this signature was created
	Expires   int64    // Unix epoch (in seconds) when this signature expires
}

func NewSignature() Signature {
	return Signature{
		Headers: make([]string, 0),
	}
}

func ParseSignature(signature string) (Signature, error) {
	return Signature{}, derp.NewInternalError("", "not implemented")
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
	buffer.WriteString(signature.Signature)
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
