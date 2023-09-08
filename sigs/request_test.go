package sigs

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPRequest(t *testing.T) {

	request, err := http.NewRequest("GET", "http://example.com/foo", nil)
	require.Nil(t, err)

	request.Header.Add("Cache-Control", "max-age=60")
	request.Header.Add("Cache-Control", "must-revalidate")

	var buffer bytes.Buffer
	if err := request.Header.Write(&buffer); err != nil {
		panic(err)
	}

	fmt.Println("---------")
	fmt.Println(buffer.String())
	fmt.Println("---------")
	fmt.Println(request.Header.Get("Cache-Control"))

	fmt.Println("---------")
	value := request.Header.Get("cAcHe-CoNtRoL")
	fmt.Println(value)

	fmt.Println("---------")
	valueSlice := request.Header[http.CanonicalHeaderKey("CAcHe-CoNtroL")]
	fmt.Println(strings.Join(valueSlice, ", "))
}
