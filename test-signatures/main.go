package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/clients"
	"github.com/benpate/hannibal/sigs"
	"github.com/benpate/hannibal/streams"
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	// Logging Configuration
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    false,
		TimeFormat: "",
	})

	fmt.Println("HTTP Signature Tester")
	fmt.Println("This is a bare-bones tool for testing")
	fmt.Println("HTTP signatures.  Paste in an HTTP request")
	fmt.Println("and it will attempt to verify the signature")
	fmt.Println("using the publick key found in the document.")
	fmt.Println("")
	fmt.Println("Paste HTTP Request Now:")
	fmt.Println("-----------------------")

	request, err := http.ReadRequest(bufio.NewReader(os.Stdin))

	if err != nil {
		derp.Report(err)
		return
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Println("Processing HTTP Request")
	fmt.Println("-----------------------")

	verifier := sigs.NewVerifier()
	if err := verifier.Verify(request, keyFinder()); err != nil {
		spew.Dump(err)
	}
	fmt.Println("")
	fmt.Println("HTTP SIGNATURE VERIFIED SUCCESSFULLY.")
}

func keyFinder() sigs.PublicKeyFinder {

	return func(keyID string) (string, error) {

		hashClient := clients.NewHashLookup(streams.NewDefaultClient())

		document, err := hashClient.Load(keyID)

		if err != nil {
			return "", derp.Wrap(err, "hannibal.validator.HTTPSig.keyFinder", "Error retrieving Actor from ActivityPub document", keyID)
		}

		publicKeyPEM := document.PublicKeyPEM()

		return publicKeyPEM, nil
	}
}
