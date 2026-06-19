package sigs

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

/******************************************
 * ParseSignature Rejections
 ******************************************/

// TestParseSignature_Reject confirms ParseSignature rejects inputs that are
// missing required fields or carry an undecodable signature.
func TestParseSignature_Reject(t *testing.T) {

	valid := `signature="Y2FiYWIxNGRiZDk4ZA=="`

	rejects := func(name string, input string) {
		t.Run(name, func(t *testing.T) {
			_, err := ParseSignature(input)
			require.NotNil(t, err, "input must be rejected")
		})
	}

	// Missing keyId.
	rejects("missing keyId", `headers="host",`+valid)

	// Missing headers.
	rejects("missing headers", `keyId="k",`+valid)

	// Missing signature.
	rejects("missing signature", `keyId="k",headers="host"`)

	// Signature that is not valid base64.
	rejects("bad base64 signature", `keyId="k",headers="host",signature="!!!not-base64!!!"`)

	// Completely empty input.
	rejects("empty", "")
}

/******************************************
 * IsExpired Truth Table
 ******************************************/

// TestSignature_IsExpired walks the IsExpired decision table across the duration,
// Expires, and Created inputs.
func TestSignature_IsExpired(t *testing.T) {

	now := time.Now().Unix()
	hourAgo := now - 3600
	hourAhead := now + 3600

	check := func(name string, sig Signature, duration int, expected bool) {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, expected, sig.IsExpired(duration))
		})
	}

	// Zero duration means "no timeout" -> never expired, regardless of fields.
	check("zero duration", Signature{Expires: hourAgo, Created: hourAgo}, 0, false)

	// Explicit Expires in the past -> expired.
	check("expires in past", Signature{Expires: hourAgo}, 3600, true)

	// Explicit Expires in the future -> not expired.
	check("expires in future", Signature{Expires: hourAhead}, 3600, false)

	// Created + duration in the past -> expired.
	check("created+duration past", Signature{Created: now - 7200}, 3600, true)

	// Created + duration still in the future -> not expired.
	check("created+duration future", Signature{Created: now}, 3600, false)

	// No Created and no Expires -> not expired (nothing to test against).
	check("no timestamps", Signature{}, 3600, false)
}

/******************************************
 * Accessors
 ******************************************/

// TestSignature_ActorID confirms ActorID strips the URL fragment from the keyID.
func TestSignature_ActorID(t *testing.T) {

	withFragment := Signature{KeyID: "https://example.com/users/alice#main-key"}
	require.Equal(t, "https://example.com/users/alice", withFragment.ActorID())

	// No fragment -> the keyID is returned unchanged.
	noFragment := Signature{KeyID: "https://example.com/users/bob"}
	require.Equal(t, "https://example.com/users/bob", noFragment.ActorID())
}

// TestSignature_AlgorithmPrefix confirms AlgorithmPrefix returns the family name.
func TestSignature_AlgorithmPrefix(t *testing.T) {
	require.Equal(t, "rsa", Signature{Algorithm: "rsa-sha256"}.AlgorithmPrefix())
	require.Equal(t, "ecdsa", Signature{Algorithm: "ecdsa-sha512"}.AlgorithmPrefix())
	require.Equal(t, "hs2019", Signature{Algorithm: "hs2019"}.AlgorithmPrefix())
}

// TestSignature_CreatedExpiresString confirms the Created/Expires string getters
// render the timestamp, or empty when zero.
func TestSignature_CreatedExpiresString(t *testing.T) {

	require.Equal(t, "1700000000", Signature{Created: 1700000000}.CreatedString())
	require.Equal(t, "", Signature{Created: 0}.CreatedString())

	require.Equal(t, "1700000000", Signature{Expires: 1700000000}.ExpiresString())
	require.Equal(t, "", Signature{Expires: 0}.ExpiresString())
}

// TestSignature_Bytes confirms Bytes returns the same content as String.
func TestSignature_Bytes(t *testing.T) {

	signature := Signature{
		KeyID:     "k",
		Algorithm: "rsa-sha256",
		Headers:   []string{"host", "date"},
		Signature: []byte("0123456789"),
	}

	require.Equal(t, signature.String(), string(signature.Bytes()))
}

/******************************************
 * HasSignature / GetSignature
 ******************************************/

// TestHasSignature confirms HasSignature reflects the presence of the header.
func TestHasSignature(t *testing.T) {

	withSig, err := http.NewRequest("GET", "https://example.com", nil)
	require.Nil(t, err)
	withSig.Header.Set("Signature", `keyId="k"`)
	require.True(t, HasSignature(withSig))

	without, err := http.NewRequest("GET", "https://example.com", nil)
	require.Nil(t, err)
	require.False(t, HasSignature(without))
}

/******************************************
 * Round Trip
 ******************************************/

// TestSignature_StringRoundTrip confirms a Signature rendered with String can be
// parsed back into an equivalent Signature.
func TestSignature_StringRoundTrip(t *testing.T) {

	original := Signature{
		KeyID:     "https://example.com/users/alice#main-key",
		Algorithm: "rsa-sha256",
		Headers:   []string{"(request-target)", "host", "date", "digest"},
		Signature: []byte("this-is-a-signature-payload"),
		Created:   1700000000,
		Expires:   1700003600,
	}

	parsed, err := ParseSignature(original.String())
	require.Nil(t, err)

	require.Equal(t, original.KeyID, parsed.KeyID)
	require.Equal(t, original.Algorithm, parsed.Algorithm)
	require.Equal(t, original.Headers, parsed.Headers)
	require.Equal(t, original.Signature, parsed.Signature)
	require.Equal(t, original.Created, parsed.Created)
	require.Equal(t, original.Expires, parsed.Expires)
}
