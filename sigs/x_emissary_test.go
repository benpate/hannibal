package sigs

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Emissary(t *testing.T) {

	rawHTTP := removeTabs(
		`POST /users/angusius_dukernen/inbox HTTP/1.1
		Host: activitypub.academy
		Accept: application/activity+json
		Date: Sun, 10 Sep 2023 10:37:22 GMT
		Digest: SHA-256=gH5r8XfomFHwvjZUZAqtbE1MpJkmbgCZ3cQjYGt7ppA=
		Content-Type: application/activity+json
		Signature: keyId="https://emdev.ddns.net/@64d68054a4bf39a519f27c67#main-key",algorithm="hs2019",headers="(request-target) host date digest",signature="NjivuB3qFMopxLw/JMdMw2pi7ao31lEZDLM+GuyRXrUHBOwGc5jZ0tO+wOgCNbZ2wVhJq7P6jk0WzEjRSAXvmntv5xlNVNFepwWhhmR9ti+HA836wr9VNhijfK/J4UaXTAbV14tCKmrmPLptDAAXz4kpWMXrCQD9OORnNUXZJShpBj+ahmMdQYOfcPFNBxy3qpfRlGhKHl10kl1ehlDXGcwG+yldJaxBQM5Rln4edyliI6Ij5l1RMo8Xq0NUlIo6ZSg2FdVWAQfGKeW/m5EP4N4wPRdItmJ2C2BR/+zzFW1fy5//FboYPzr/887exxtw4aVvF73m6QwOx+zDs8oJ+w=="
		Content-Length: 338

		{"@context":"application/activity+json","actor":"https://emdev.ddns.net/@64d68054a4bf39a519f27c67","id":"https://emdev.ddns.net/@64d68054a4bf39a519f27c67/pub/following/64fdf0c1a835dc65840d6ffb","object":"https://activitypub.academy/users/angusius_dukernen","to":"https://activitypub.academy/users/angusius_dukernen/inbox","type":"Follow"}`)

	publicPEM := "-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEAusAu2SkpV2zpD4Yhzwc4fV2RDMNLoTQp+gS2w5xUaLHwyPJckCZi\n5l4Aj9OvebwGXKicAUund8vMiRFeTeOlV6+UUK8/kV0WWbO2cWlvvqcOyBm4VKxz\nigHKGVD0HM2QmSAf4XyXRLuCbJjawUm36xZgGRgkG/YDFD+0YqQWVmjpwVSjiXc8\nx0GqV9LottkNeQf8YMgNEp8sW0fY7h2RmkdILNqJ7rX3UjE6MBXg+rlJqAGy/Zge\nuGfQIAe+bA53onjJNRonSLKpeYMSGKAKWONfHMnWScqBy7b79OOKLGIzh9v/3hk7\nEvKGydvXckxu9KW2u/0wFV7NVLbS5XEo0wIDAQAB\n-----END RSA PUBLIC KEY-----\n"

	// Make a new request
	request, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawHTTP)))
	require.Nil(t, err)

	// Verify the request
	err = Verify(request, publicPEM, VerifierIgnoreTimeout())
	require.Nil(t, err)
}

func Test_Emissary_Digest(t *testing.T) {
	expected := "SHA-256=gH5r8XfomFHwvjZUZAqtbE1MpJkmbgCZ3cQjYGt7ppA="
	body := `{"@context":"application/activity+json","actor":"https://emdev.ddns.net/@64d68054a4bf39a519f27c67","id":"https://emdev.ddns.net/@64d68054a4bf39a519f27c67/pub/following/64fdf0c1a835dc65840d6ffb","object":"https://activitypub.academy/users/angusius_dukernen","to":"https://activitypub.academy/users/angusius_dukernen/inbox","type":"Follow"}`

	digestBytes := sha256.Sum256([]byte(body))
	digest := base64.StdEncoding.EncodeToString(digestBytes[:])

	require.Equal(t, expected, "SHA-256="+digest)
	require.Equal(t, expected, DigestSHA256([]byte(body)))
}
