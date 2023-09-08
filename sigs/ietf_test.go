package sigs

import (
	"bufio"
	"bytes"
	"crypto"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

/******************************************
 * These tests come from the current IETF draft spec:
 * https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#section-2.4
******************************************/

// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#appendix-C.1
func Test_IETF_Default_Test_C1(t *testing.T) {
	signedFields := []string{"date"}
	request, body := test_IETF_Request()

	// Sign the document
	err := Sign(request, body, "Test", test_IETF_PrivateKey(), SignerFields(signedFields...))
	require.Nil(t, err)

	// Check the signature with "require"
	expectedSignature := `keyId="Test",algorithm="hs2019",headers="date",signature="SjWJWbWN7i0wzBvtPl8rbASWz5xQW6mcJmn+ibttBqtifLN7Sazz6m79cNfwwb8DMJ5cou1s7uEGKKCs+FLEEaDV5lp7q25WqS+lavg7T8hc0GppauB6hbgEKTwblDHYGEtbGmtdHgVCk9SuS13F0hZ8FD0k/5OxEPXe5WozsbM="`
	require.Equal(t, expectedSignature, request.Header.Get("Signature"))

	// Verify the signature
	err = Verify(request, body, test_IETF_PublicPEM(), VerifierFields(signedFields...), VerifierIgnoreTimeout())
	require.Nil(t, err)
}

// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#appendix-C.2
func Test_IETF_Basic_Test_C2(t *testing.T) {
	signedFields := []string{"(request-target)", "host", "date"}
	request, body := test_IETF_Request()

	// Sign the document
	err := Sign(request, body, "Test", test_IETF_PrivateKey(), SignerFields(signedFields...))
	require.Nil(t, err)

	// Check the signature with "require"
	expectedSignature := `keyId="Test",algorithm="hs2019",headers="(request-target) host date",signature="qdx+H7PHHDZgy4y/Ahn9Tny9V3GP6YgBPyUXMmoxWtLbHpUnXS2mg2+SbrQDMCJypxBLSPQR2aAjn7ndmw2iicw3HMbe8VfEdKFYRqzic+efkb3nndiv/x1xSHDJWeSWkx3ButlYSuBskLu6kd9Fswtemr3lgdDEmn04swr2Os0="`
	require.Equal(t, expectedSignature, request.Header.Get("Signature"))

	// Verify the signature
	err = Verify(request, body, test_IETF_PublicPEM(), VerifierFields(signedFields...), VerifierIgnoreTimeout())
	require.Nil(t, err)
}

func Test_IETF_All_Headers_Prep(t *testing.T) {
	signedFields := []string{"(request-target)", "(created)", "(expires)", "host", "date", "content-type", "digest", "content-length"}
	request, body := test_IETF_Request()

	signature := NewSignature()
	signature.Created = 1402170695
	signature.Expires = 1402170699

	plaintext := makePlaintext(request, signature, signedFields...)
	t.Log(plaintext)
	t.Log(string(body))
}

/*
	Removing this for now because we don't have a mechanism to pass created/expires vaues into the Signer.
	Plus, it may not be so important because the spec says that the created/expires values SHOULD NOT be used
	for RSA keys.

// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#appendix-C.3

	func Test_IETF_All_Headers_Test_C3(t *testing.T) {
		signedFields := []string{"(request-target)", "(created)", "(expires)", "host", "date", "content-type", "digest", "content-length"}
		request, body := test_IETF_Request()

		err := Sign(request, body, "Test", test_IETF_PrivateKey(), SignerFields(signedFields...))
		require.Nil(t, err)

		// handling this one differently because we can't *force* the created/expires dates to be these correct values
		// so let's just Verify this signature straight out.
		expectedSignature := `keyId="Test",algorithm="rsa-sha256",created=1402170695,expires=1402170699,headers="(request-target) (created) (expires) host date content-type digest content-length",signature="vSdrb+dS3EceC9bcwHSo4MlyKS59iFIrhgYkz8+oVLEEzmYZZvRs8rgOp+63LEM3v+MFHB32NfpB2bEKBIvB1q52LaEUHFv120V01IL+TAD48XaERZFukWgHoBTLMhYS2Gb51gWxpeIq8knRmPnYePbF5MOkR0Zkly4zKH7s1dE="`
		request.Header.Set("Signature", expectedSignature)

		err = Verify(request, body, test_IETF_PublicPEM(), VerifierIgnoreTimeout())
		require.Nil(t, err)
	}
*/
func test_IETF_Request() (*http.Request, []byte) {

	requestString := removeTabs(
		`POST /foo?param=value&pet=dog HTTP/1.1
		Host: example.com
		Date: Sun, 05 Jan 2014 21:31:40 GMT
		Content-Type: application/json
		Digest: SHA-256=X48E9qOokqqrvdts8nOJRJN3OWDUoyWxBf7kbu9DBPE=
		Content-Length: 18

		{"hello": "world"}`)

	bodyReader := bufio.NewReader(bytes.NewReader([]byte(requestString)))
	request := must(http.ReadRequest(bodyReader))
	body := must(io.ReadAll(request.Body))

	requestReader := bufio.NewReader(bytes.NewReader([]byte(requestString)))
	result := must(http.ReadRequest(requestReader))

	return result, body
}

func test_IETF_PublicPEM() string {
	return removeTabs(
		`-----BEGIN PUBLIC KEY-----
		MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
		6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
		Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
		oYi+1hqp1fIekaxsyQIDAQAB
		-----END PUBLIC KEY-----`)
}

func test_IETF_PrivateKey() crypto.PrivateKey {

	privatePEM := removeTabs(
		`-----BEGIN RSA PRIVATE KEY-----
		MIICXgIBAAKBgQDCFENGw33yGihy92pDjZQhl0C36rPJj+CvfSC8+q28hxA161QF
		NUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6Z4UMR7EOcpfdUE9Hf3m/hs+F
		UR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJwoYi+1hqp1fIekaxsyQIDAQAB
		AoGBAJR8ZkCUvx5kzv+utdl7T5MnordT1TvoXXJGXK7ZZ+UuvMNUCdN2QPc4sBiA
		QWvLw1cSKt5DsKZ8UETpYPy8pPYnnDEz2dDYiaew9+xEpubyeW2oH4Zx71wqBtOK
		kqwrXa/pzdpiucRRjk6vE6YY7EBBs/g7uanVpGibOVAEsqH1AkEA7DkjVH28WDUg
		f1nqvfn2Kj6CT7nIcE3jGJsZZ7zlZmBmHFDONMLUrXR/Zm3pR5m0tCmBqa5RK95u
		412jt1dPIwJBANJT3v8pnkth48bQo/fKel6uEYyboRtA5/uHuHkZ6FQF7OUkGogc
		mSJluOdc5t6hI1VsLn0QZEjQZMEOWr+wKSMCQQCC4kXJEsHAve77oP6HtG/IiEn7
		kpyUXRNvFsDE0czpJJBvL/aRFUJxuRK91jhjC68sA7NsKMGg5OXb5I5Jj36xAkEA
		gIT7aFOYBFwGgQAQkWNKLvySgKbAZRTeLBacpHMuQdl1DfdntvAyqpAZ0lY0RKmW
		G6aFKaqQfOXKCyWoUiVknQJAXrlgySFci/2ueKlIE1QqIiLSZ8V8OlpFLRnb1pzI
		7U1yQXnTAEFYM560yJlzUpOb1V4cScGd365tiSMvxLOvTA==
		-----END RSA PRIVATE KEY-----`)

	privateKey := must(DecodePrivatePEM(privatePEM))

	return privateKey
}
