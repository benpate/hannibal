package sigs

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestVerify_Manual(t *testing.T) {

	// Configure logging
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// zerolog.SetGlobalLevel(zerolog.TraceLevel)

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
	err = Sign(request, body.Bytes(), "test-key", privateKey)
	require.Nil(t, err)

	require.Equal(t, "SHA-256=65F8+S1Bg7oPQS/fIxVg4x7PoLWnOxWlGMFB/hafojg=", request.Header.Get("Digest"))
	require.NotEmpty(t, request.Header.Get("Signature"))

	// Verify the Request
	publicKeyPEM := EncodePublicPEM(privateKey)

	err = Verify(request, body.Bytes(), publicKeyPEM)
	require.Nil(t, err)
}

func TestVerify_Manual2(t *testing.T) {

	// sample values from https://blog.cubieserver.de/2016/go-verify-cryptographic-signatures/
	var rawPubKey = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvtjdLkS+FP+0fPC09j25\ny/PiuYDDivIT86COVedvlElk99BBYTrqNaJybxjXbIZ1Q6xFNhOY+iTcBr4E1zJu\ntizF3Xi0V9tOuP/M8Wn4Y/1lCWbQKlWrNQuqNBmhovF4K3mDCYswVbpgTmp+JQYu\nBm9QMdieZMNry5s6aiMA9aSjDlNyedvSENYo18F+NYg1J0C0JiPYTxheCb4optr1\n5xNzFKhAkuGs4XTOA5C7Q06GCKtDNf44s/CVE30KODUxBi0MCKaxiXw/yy55zxX2\n/YdGphIyQiA5iO1986ZmZCLLW8udz9uhW5jUr3Jlp9LbmphAC61bVSf4ou2YsJaN\n0QIDAQAB\n-----END PUBLIC KEY-----"
	var rawSignature = "c2pkYWpuY2sgZmphbm9panF3b2lqYWRvbmFzbWQgc2EsbWMgc2FuZHBvZHA5cTN1cjA5M3Vyajg4OUoocHEqaDlIUkZKU0ZLQkZPSDk4"
	var message = []byte("authenticmessage")

	publicKey, err := DecodePublicPEM(rawPubKey)
	require.Nil(t, err)

	signature, err := base64.StdEncoding.DecodeString(rawSignature)
	require.Nil(t, err)

	hashedMessage := sha1.Sum(message)

	// err = verifySignature(publicKey, crypto.SHA1, hashedMessage[:], signature)
	err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hashedMessage[:], signature)
	require.Nil(t, err)

	require.Nil(t, err)
}

func TestVerify_PixelFed(t *testing.T) {

	actor := streams.NewDocument("https://pixelfed.social/users/benpate")
	publicKeyPEM := actor.PublicKey().PublicKeyPEM()
	log.Trace().Interface("actor", actor.Value()).Send()
	log.Trace().Str("publicKeyPEM", publicKeyPEM).Send()

	request, body := getTestPixelFedRequest()

	signature := NewSignature()
	plaintext := makePlaintext(request, signature, "(request-target)", "host", "date", "digest", "content-type", "user-agent")
	fmt.Println("---------")
	fmt.Println(string(body))
	fmt.Println("---------")
	fmt.Println(plaintext)

	hashed, err := makeSignatureHash(plaintext, crypto.SHA256)
	require.Nil(t, err)
	fmt.Println("---------")
	fmt.Println("hashed plaintext")
	fmt.Println(base64.StdEncoding.EncodeToString(hashed))
	fmt.Println("---------")

	err = Verify(request, body, publicKeyPEM, VerifierIgnoreTimeout())
	require.Nil(t, err)
}

func TestSignAndVerify_RSASHA256(t *testing.T) {

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

func TestSignAndVerify_RSASHA512(t *testing.T) {

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

	fmt.Println("plaintextSignature -------------------")
	fmt.Println(plaintextSignature)
	fmt.Println("")

	hashedSignature, err := makeSignatureHash(plaintextSignature, crypto.SHA256)
	require.Nil(t, err)

	fmt.Println("hashedSignature -------------------")
	fmt.Println(hashedSignature)
	fmt.Println("")
	fmt.Println(base64.StdEncoding.EncodeToString(hashedSignature))
	fmt.Println("")

	signedSignature, err := makeSignedDigest(hashedSignature, crypto.SHA256, privateKey)
	require.Nil(t, err)
	fmt.Println("signedSignature -------------------")
	fmt.Println(signedSignature)
	fmt.Println("")
	fmt.Println(base64.StdEncoding.EncodeToString(signedSignature))
	fmt.Println("")

	fmt.Println("original Signature --------------------")
	signature := `xFVaZiSGXc/6KBHEeg8qssbqVwowTA2pQKD46QcbojZLVP90kJpDutZ0DuSYlMjPlRc95meFb+O0B3ikqA8MEyoDQ+1xyn5o7o+zKteTb6FQf3feBDvGWxJh3DWIog3h8Hqxzha++wctUYm1AXSabd1enQwhzlseSJFfie3P0Wr5MTd96H+NcUizZIrapsE8Q8Wbdixaywp1lehrHt/ah72ocbvHX1AVlPAjdfxJ+tftI2MJ03qRshssNtk+r2DMufGiNb8wIEA4E7928+LqDpjZiUmQoSbxto5MTOtNV58wurz5YaBQjpD8zSHmzV17RGWVQYGII+v/aHqXPKD7Gw==`
	fmt.Println(signature)
	fmt.Println("")

	fmt.Println("verifyHashAndSignature -------------------")
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
	fmt.Println(privateKeyPEM)

	privateKey, err := DecodePrivatePEM(privateKeyPEM)

	if err != nil {
		panic(err)
	}

	testPrivateKey := privateKey.(*rsa.PrivateKey)
	testPublicKeyPEM := EncodePublicPEM(testPrivateKey)
	fmt.Println(testPublicKeyPEM)

	/* sent via browser??
	publicKeyPEM := removeTabs(
		`-----BEGIN RSA PUBLIC KEY-----
		MEgCQQDbLVt+d4EGWdMOgG6lS2xvhP6kbb0OgdkG26jmqWfUCqzYhyuhoL3JgijV
		N+Y0Jbb4iEU2aQXMNHM+Rq1bfkLTAgMBAAE=
		-----END RSA PUBLIC KEY-----`)
	*/

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

	fmt.Println(publicKeyPEM)

	if publicKeyPEM != testPublicKeyPEM {
		panic("Public PEMs do not match")
	}

	publicKey, err := DecodePublicPEM(publicKeyPEM)

	if err != nil {
		panic(err)
	}

	return privateKey, publicKey
}
