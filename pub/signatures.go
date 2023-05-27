package pub

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/mapof"
	"github.com/go-fed/httpsig"
)

/******************************************
 * HTTP Signatures
 *
 * https://docs.joinmastodon.org/spec/security/
 *
 ******************************************/

// ValidateHTTPSignature verifies that the HTTP request is signed with a valid key.
// This function loads the public key from the ActivityPub actor, then verifies their signature.
func ValidateHTTPSignature(request *http.Request, document streams.Document) error {

	headerValues := ParseSignatureHeader(request.Header.Get("Signature"))

	const location = "activitypub.validateRequest"

	// Get the Actor from the document
	actor, err := document.Actor().Load()

	if err != nil {
		return derp.Wrap(err, location, "Error retrieving Actor from ActivityPub document")
	}

	// Get the Actor's Public Key
	actorPublicKey, err := actor.PublicKey().Load()

	if err != nil {
		return derp.Wrap(err, location, "Error retrieving Public Key from Actor")
	}

	// Get the PEM from the public key
	actorPublicKeyPEM := actorPublicKey.PublicKeyPEM()

	// Finally, Verify request signatures
	verifier, err := httpsig.NewVerifier(request)

	if err != nil {
		return derp.Wrap(err, location, "Error creating HTTP Signature verifier")
	}

	// Parse the correct key type from the request
	algorithm := httpsig.Algorithm(headerValues["algorithm"])
	key, err := ParsePublicKeyFromPEM(actorPublicKeyPEM)

	if err != nil {
		return derp.Wrap(err, location, "Error parsing Public Key")
	}

	// Use httpsig to verify the request
	if err := verifier.Verify(key, algorithm); err != nil {
		return derp.Wrap(err, location, "Error verifying HTTP Signature")
	}

	return nil
}

func ParsePublicKeyFromPEM(pemString string) (crypto.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemString))

	if block == nil {
		return nil, derp.New(derp.CodeInternalError, "pub.ParseKeyFromPEM", "Block is nil", pemString)
	}

	switch block.Type {

	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)

	case "PUBLIC KEY":
		return x509.ParsePKIXPublicKey(block.Bytes)

	default:
		return nil, derp.New(derp.CodeInternalError, "pub.ParseKeyFromPEM", "Invalid block type", block.Type)
	}
}

func ParseSignatureHeader(value string) mapof.String {

	result := mapof.NewString()

	item := ""
	itemList := list.ByComma(value)

	for !itemList.IsEmpty() {
		item, itemList = itemList.Split()
		name, value := list.Split(item, '=')
		value = strings.Trim(value, "\" ")

		result[name] = value
	}

	return result
}
