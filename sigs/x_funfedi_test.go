package sigs

import (
	"bufio"
	"crypto"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunFedi(t *testing.T) {

	// From transactions sent by https://verify.funfedi.dev/
	// FunFedi is a test suite that uses different capitalization for the Digest header (sha-256 instead of SHA-256)
	// This test ensures that the Digest header is case-insensitive.

	// zerolog.SetGlobalLevel(zerolog.TraceLevel)

	/*
		------------------------------------------
		HANNIBAL: Receiving Activity: /@64d68054a4bf39a519f27c67/pub/inbox
		Headers:
		Digest: sha-256=27p0TuEIcJbNLBjv/RQFROHFxe0K74PK2exvfyHkkDQ=
		Content-Type: application/activity+json
		Signature: keyId="https://verify.funfedi.dev/bob#main",algorithm="rsa-sha256",headers="(request-target) host date digest content-type",signature="b2k1vPoLJpuCk1MmAEk6pfWi5G8SFALBqOjywUNdOEiC9SeTEPULCDi5quLPqzlvsSoD+jHipzTlETYwnen9wkwYqKzBlp5sTMbdKEXI1L6dzE4mmqMqE5zCGgzJqAlK59Z7PQZGTegJ/qAXjywBPcJC7TB4yD9JpbNPBJ6DcqBk3wGMh0rTxMNg4m9Wj90lrmYF+fqNxUkUHPdXxG7TxlaiQ18Z5RWZoXGv0+lOpNrhRU44J9Dl98aiKnhm+xoRrE+QUBKLEmKpwJU+bBsd1R7s9IV6P2JjYL2paOWIOveaNt41GcPHUc5g5aUkQfmMbWVeWv6VM7lTzpfO3e93Ww=="
		Accept: application/activity+json
		Accept-Encoding: gzip, deflate
		Content-Length: 207
		User-Agent: bovine/0.5.3
		Date: Tue, 05 Dec 2023 21:22:25 GMT

		Body:
		{"@context": "https://www.w3.org/ns/activitystreams", "type": "Like", "actor": "https://verify.funfedi.dev/bob", "id": "https://verify.funfedi.dev/bobOEFNNp884mw", "object": "https://verify.funfedi.dev/bob"}

		2:22PM TRC ascache.Client.loadByURI: FOUND uri=https://verify.funfedi.dev/bob
		------------------------------------------
		hannibal.pub.validateRequest
		(map[string]interface {}) (len=9) {
		 (string) (len=4) "type": (string) (len=6) "Person",
		 (string) (len=8) "@context": (primitive.A) (len=2 cap=2) {
		  (string) (len=37) "https://www.w3.org/ns/activitystreams",
		  (string) (len=28) "https://w3id.org/security/v1"
		 },
		 (string) (len=17) "preferredUsername": (string) (len=3) "bob",
		 (string) (len=9) "publicKey": (mapof.Any) (len=3) {
		  (string) (len=2) "id": (string) (len=35) "https://verify.funfedi.dev/bob#main",
		  (string) (len=5) "owner": (string) (len=30) "https://verify.funfedi.dev/bob",
		  (string) (len=12) "publicKeyPem": (string) (len=451) "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo/r0o1lp0IUe6Y+IFm6Q\naHmMkyGXdHy9mE4l7+5AKQBGb8c3n6dDVIiECvrdmF1H8U1lsI/Q1nq8lQkuzxBV\nysmAPHFusW0ODy1NYGTEGYGnjfWuttltYGf8JgSzQMxUFnzg2PVXCmAq+QK3eENK\nm0xMc1EKagY5BBOtOljAP2iN0gdsb3RQ7mQHzBcZCataiMI52qVt/M/7Zony5W8e\nQWbLMPr3WMs+JPwz5TIVED4UMJxFswS5+yI1iQjgHgXdcw63ipJ/QWy/dtDU8llD\ne0TVR+KdKTxHpl2P3ky+OK6zYIO2MFfru8IDrax4i/zK1VTMzd9BipmoFdlK/5dw\n3wIDAQAB\n-----END PUBLIC KEY-----\n"
		 },
		 (string) (len=7) "summary": (string) (len=88) "user_part=bob\nrequires_signed_get_for_actor=False,\nrequires_signed_post_for_inbox=True,\n",
		 (string) (len=2) "id": (string) (len=30) "https://verify.funfedi.dev/bob",
		 (string) (len=5) "inbox": (string) (len=36) "https://verify.funfedi.dev/bob/inbox",
		 (string) (len=6) "outbox": (string) (len=30) "https://verify.funfedi.dev/bob",
		 (string) (len=4) "name": (string) (len=10) "Test Actor"
		}
		PEM: -----BEGIN PUBLIC KEY-----
		MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo/r0o1lp0IUe6Y+IFm6Q
		aHmMkyGXdHy9mE4l7+5AKQBGb8c3n6dDVIiECvrdmF1H8U1lsI/Q1nq8lQkuzxBV
		ysmAPHFusW0ODy1NYGTEGYGnjfWuttltYGf8JgSzQMxUFnzg2PVXCmAq+QK3eENK
		m0xMc1EKagY5BBOtOljAP2iN0gdsb3RQ7mQHzBcZCataiMI52qVt/M/7Zony5W8e
		QWbLMPr3WMs+JPwz5TIVED4UMJxFswS5+yI1iQjgHgXdcw63ipJ/QWy/dtDU8llD
		e0TVR+KdKTxHpl2P3ky+OK6zYIO2MFfru8IDrax4i/zK1VTMzd9BipmoFdlK/5dw
		3wIDAQAB
		-----END PUBLIC KEY-----
	*/

	rawHTTP := removeTabs(
		`POST /@64d68054a4bf39a519f27c67/pub/inbox HTTP/1.1
		Host: emdev.ddns.net
		Digest: sha-256=27p0TuEIcJbNLBjv/RQFROHFxe0K74PK2exvfyHkkDQ=
		Content-Type: application/activity+json
		Signature: keyId="https://verify.funfedi.dev/bob#main",algorithm="rsa-sha256",headers="(request-target) host date digest content-type",signature="b2k1vPoLJpuCk1MmAEk6pfWi5G8SFALBqOjywUNdOEiC9SeTEPULCDi5quLPqzlvsSoD+jHipzTlETYwnen9wkwYqKzBlp5sTMbdKEXI1L6dzE4mmqMqE5zCGgzJqAlK59Z7PQZGTegJ/qAXjywBPcJC7TB4yD9JpbNPBJ6DcqBk3wGMh0rTxMNg4m9Wj90lrmYF+fqNxUkUHPdXxG7TxlaiQ18Z5RWZoXGv0+lOpNrhRU44J9Dl98aiKnhm+xoRrE+QUBKLEmKpwJU+bBsd1R7s9IV6P2JjYL2paOWIOveaNt41GcPHUc5g5aUkQfmMbWVeWv6VM7lTzpfO3e93Ww=="
		Accept: */*
		Accept-Encoding: gzip, deflate
		Content-Length: 207
		User-Agent: bovine/0.5.3
		Date: Tue, 05 Dec 2023 21:22:25 GMT
		
		{"@context": "https://www.w3.org/ns/activitystreams", "type": "Like", "actor": "https://verify.funfedi.dev/bob", "id": "https://verify.funfedi.dev/bobOEFNNp884mw", "object": "https://verify.funfedi.dev/bob"}
		`)

	keyFinder := func(keyID string) (string, error) {
		return "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo/r0o1lp0IUe6Y+IFm6Q\naHmMkyGXdHy9mE4l7+5AKQBGb8c3n6dDVIiECvrdmF1H8U1lsI/Q1nq8lQkuzxBV\nysmAPHFusW0ODy1NYGTEGYGnjfWuttltYGf8JgSzQMxUFnzg2PVXCmAq+QK3eENK\nm0xMc1EKagY5BBOtOljAP2iN0gdsb3RQ7mQHzBcZCataiMI52qVt/M/7Zony5W8e\nQWbLMPr3WMs+JPwz5TIVED4UMJxFswS5+yI1iQjgHgXdcw63ipJ/QWy/dtDU8llD\ne0TVR+KdKTxHpl2P3ky+OK6zYIO2MFfru8IDrax4i/zK1VTMzd9BipmoFdlK/5dw\n3wIDAQAB\n-----END PUBLIC KEY-----", nil
	}

	// Make a new request
	request, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawHTTP)))
	require.Nil(t, err)

	err = VerifyDigest(request, crypto.SHA256)
	require.Nil(t, err)

	// Verify the request
	err = Verify(request, keyFinder, VerifierIgnoreTimeout())
	require.Nil(t, err)

}
