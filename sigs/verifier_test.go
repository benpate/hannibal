package sigs

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"net/http"
	"testing"

	"github.com/benpate/derp"
	"github.com/stretchr/testify/require"
)

func TestVerify_Manual(t *testing.T) {

	// Create a new Request
	var body bytes.Buffer
	body.WriteString("This is the body of the request")

	request, err := http.NewRequest("GET", "http://example.com/something?test=true", &body)
	require.Nil(t, err)
	request.Header.Set("Content-Type", "text/plain")

	// Create a Private Key to sign the request
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	// Sign the Request
	err = Sign(request, "test-key", privateKey)
	require.Nil(t, err)

	require.Equal(t, "SHA-256=65F8+S1Bg7oPQS/fIxVg4x7PoLWnOxWlGMFB/hafojg=", request.Header.Get("Digest"))
	require.NotEmpty(t, request.Header.Get("Signature"))

	// Verify the Request
	keyFinder := func(keyID string) (string, error) {
		return EncodePublicPEM(privateKey), nil
	}

	_, err = Verify(request, keyFinder)
	require.Nil(t, err)
}

func TestSignAndVerify_RSA_SHA256(t *testing.T) {

	// Create an RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	// Create a test digest
	digest := sha256.Sum256([]byte("this is the message"))

	// Sign the digest
	signature, err := makeSignedDigest(digest[:], crypto.SHA256, privateKey)
	require.Nil(t, err)

	err = verifySignature(&privateKey.PublicKey, crypto.SHA256, digest[:], signature)
	require.Nil(t, err)
}

func TestSignAndVerify_RSA_SHA512(t *testing.T) {

	// Create an RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	// Create a test digest
	digest := sha512.Sum512([]byte("this is the message"))

	// Sign the digest
	signature, err := makeSignedDigest(digest[:], crypto.SHA512, privateKey)
	require.Nil(t, err)

	err = verifySignature(&privateKey.PublicKey, crypto.SHA512, digest[:], signature)
	require.Nil(t, err)
}

func TestSignAndVerify_Emissary(t *testing.T) {

	privateKey, publicKey := getTestKeys()
	require.NotNil(t, privateKey)
	require.NotNil(t, publicKey)

	plaintextSignature := removeTabs(
		`(request-target): post /@63f3c38ab0d6029bf1c915c9/pub/inbox
		host: 127.0.0.1
		date: Sat, 02 Sep 2023 15:29:05 GMT
		digest: SHA-256=VgRwIlTTedMmiNk71Zlod3Bq6gXOzJNIRlUXenyNdLw=`)

	hashedSignature, err := makeSignatureHash(plaintextSignature, crypto.SHA256)
	require.Nil(t, err)

	signedSignature, err := makeSignedDigest(hashedSignature, crypto.SHA256, privateKey)
	require.Nil(t, err)

	derp.Report(verifyHashAndSignature(plaintextSignature, crypto.SHA256, publicKey, signedSignature))
}

func getTestKeys() (crypto.PrivateKey, crypto.PublicKey) {

	privateKeyPEM := removeTabs(
		`-----BEGIN RSA PRIVATE KEY-----
		MIIEpQIBAAKCAQEAy1xxGJw8d+FouEHikqkmNo/X8/tPAMtZtzXXj03Uzr3Pxfpy
		4a0MhZwd3duWqhINPKWbDFgn9W2z6I+nIziBLD+YxHWqvahpsRqGkmu86CoOLKom
		mbULjAIzyAMtPqBpOQJ8xJtq6Evz09avUku08iPrjP64wKNESyu5mDFvfpW31F6B
		7C0y+QC6vbhDanOnvV9QIxMDEbU87iY3nyyt8ZkSj5I2bHb80LQ0BEWN4WkOZB+w
		c0+fhQ9+pJobSSsyGJ21graTbkEKcr1LGo+Xe+rqPYT1IcDwpMTD7es1AiqbZwlI
		xNoh9wvJygZsqB4Iok8iatc+I1fGl6XiJcnxAQIDAQABAoIBAGANotG4AguxqU/W
		ttkFEiqVWLBCFHfQlOingtCKN6kLGJdvi1Gy9gYpziWbcZeU/TGXGxwCi6UuEtsW
		9x/4sXKf+11YIrSAVqOzXrrMLqcOLjHEkITrca/I3oJrlbRN+kVWOm525lEghuOZ
		NKhPYAE7HCg1rDg5JanH1lrfhsUnxF7XXhR48+8T1iifd31DpY6M3FAP4fGSdsOA
		23bW7VG/nOwKrZ/V6ij0F3KvJ5dCOc1tw5SwhlTH9kuEoMq8Uy0i+3ngNkv7e0qJ
		k2/0B7nc6B6HP+nsswN3KoY4qKvPfktZDZeowmcNrJOBiw0N3PBrOIOI5UjkwBUN
		TOfnRUECgYEA1J1WMeag0Mr3d/rVRVqTXlwEkF7T9s7+MAPfJ5vF6u5ndrZdwy6g
		6WdbiJMkUbAH1QFQ1q7MTNoo9L72dRkS5cnSwv7WLBC2Oy3ZpsBT9Lc9SZ7GAj/X
		9x3XstBBIGkPMAfva/OCACPKd3Z84gTeDNTU+cYghv+EYhLd7UKhQ1kCgYEA9Nu3
		bpS2FguZ/CwkmRQTgTGAGysETjgLFHDxODEZPKBkTpnAf49In+Vx6BixZ6xqRWe9
		lSnlF20x8+pUOqhMpfnwxx+OEMayr4dGMDhDj6pQ+D8ETEOQ9PXqyuyyHq8FfIAH
		kdXddjCS8SSreBMgzjUZjn/md+ArA6UMN9T1LekCgYEAop/13hVZzFpzDwJ9Pp8Z
		OYOIuiTOXGnXY0KS3ej4acoQuWykKzbvPZghG0Xw8cqDMxnei1cITYBQ82NdgBO9
		sKW+4AesKeheesWHRVS24ueFqVoYen/64Lmi0tMX/YJea46mQxvuw8ycgOPQgdDX
		R1lDzgkNuDSZParQtTnRv4ECgYEAq/3RgPkwVZfcl8ciBeyWLr9obqzun0q6badP
		qNrEEVPQYW2aS3+H0djHA/KkWmA/XXUbM7Vz19q5pc1JUNJ61HMV76h4j8wiIy1v
		3dsHidhme5k4GaG0Jny+ab+M9gSWY/dCWevRXX2NGZlaYEN/XZjq1K9+YWGylSLP
		zD/n4FECgYEAssvQs++XNyNi1YQvUOX0oZjKH1l9sx/JdE+O+UYsvOsCMRXPfLu4
		ipiwM6zWyC/tumuDlM9qReAKcqBtSNUOHcd9Tbs5zTUvm3RbhqUdgMPZfVG/O9qp
		7FXferRokfJyjJpo7SW8K57HKPJ4aNMpa1yhG+uZ1rXzqvjD6zwRvto=
		-----END RSA PRIVATE KEY-----`)

	privateKey, err := DecodePrivatePEM(privateKeyPEM)

	if err != nil {
		panic(err)
	}

	testPrivateKey := privateKey.(*rsa.PrivateKey)
	testPublicKeyPEM := EncodePublicPEM(testPrivateKey)

	// in database
	publicKeyPEM := removeTabs(
		`-----BEGIN RSA PUBLIC KEY-----
		MIIBCgKCAQEAy1xxGJw8d+FouEHikqkmNo/X8/tPAMtZtzXXj03Uzr3Pxfpy4a0M
		hZwd3duWqhINPKWbDFgn9W2z6I+nIziBLD+YxHWqvahpsRqGkmu86CoOLKommbUL
		jAIzyAMtPqBpOQJ8xJtq6Evz09avUku08iPrjP64wKNESyu5mDFvfpW31F6B7C0y
		+QC6vbhDanOnvV9QIxMDEbU87iY3nyyt8ZkSj5I2bHb80LQ0BEWN4WkOZB+wc0+f
		hQ9+pJobSSsyGJ21graTbkEKcr1LGo+Xe+rqPYT1IcDwpMTD7es1AiqbZwlIxNoh
		9wvJygZsqB4Iok8iatc+I1fGl6XiJcnxAQIDAQAB
		-----END RSA PUBLIC KEY-----
		`)

	if publicKeyPEM != testPublicKeyPEM {
		panic("Public PEMs do not match")
	}

	publicKey, err := DecodePublicPEM(publicKeyPEM)

	if err != nil {
		panic(err)
	}

	return privateKey, publicKey
}
