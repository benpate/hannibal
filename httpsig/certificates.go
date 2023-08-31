package httpsig

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/benpate/derp"
)

// EncodePrivatePEM converts a private key into a PEM string
func EncodePrivatePEM(privateKey *rsa.PrivateKey) string {

	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return string(privatePEM)
}

// EncodePublicPEM converts a public key into a PEM string
func EncodePublicPEM(privateKey *rsa.PrivateKey) string {

	// Get ASN.1 DER format
	publicDER := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	// pem.Block
	publicBlock := pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   publicDER,
	}

	// Private key in PEM format
	publicPEM := pem.EncodeToMemory(&publicBlock)

	return string(publicPEM)
}

// DecodePrivatePEM converts a PEM string into a private key
func DecodePrivatePEM(pemString string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemString))

	if block == nil {
		return nil, derp.New(derp.CodeInternalError, "hannibal.httpsig.DecodePrivatePEM", "Block is nil", pemString)
	}

	switch block.Type {

	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)

	case "PRIVATE KEY":
		return x509.ParsePKCS8PrivateKey(block.Bytes)

	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)

	default:
		return nil, derp.New(derp.CodeInternalError, "hannibal.httpsig.DecodePrivatePEM", "Invalid block type", block.Type)
	}
}

// DecodePublicPEM converts a PEM string into a public key
func DecodePublicPEM(pemString string) (crypto.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemString))

	if block == nil {
		return nil, derp.New(derp.CodeInternalError, "hannibal.httpsig.DecodePublicPEM", "Block is nil", pemString)
	}

	switch block.Type {

	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)

	case "PUBLIC KEY":
		return x509.ParsePKIXPublicKey(block.Bytes)

	default:
		return nil, derp.New(derp.CodeInternalError, "hannibal.httpsig.DecodePublicPEM", "Invalid block type", block.Type)
	}
}
