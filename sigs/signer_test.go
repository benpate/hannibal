package sigs

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPathAndQuery(t *testing.T) {
	url, _ := url.Parse("http://example.com")
	require.Equal(t, "/", getPathAndQuery(url))

	url, _ = url.Parse("http://example.com/")
	require.Equal(t, "/", getPathAndQuery(url))

	url, _ = url.Parse("http://example.com/something")
	require.Equal(t, "/something", getPathAndQuery(url))

	url, _ = url.Parse("http://example.com/something?test=true")
	require.Equal(t, "/something?test=true", getPathAndQuery(url))
}

func TestGetFields(t *testing.T) {

	bodyReader := strings.NewReader("This is the body of the request")

	request, err := http.NewRequest("GET", "http://example.com/something?test=true", bodyReader)
	require.Nil(t, err)
	request.Header.Set("Content-Type", "text/plain")

	result := makePlaintext(request, FieldRequestTarget, FieldHost, "Content-Type")
	fmt.Println(result)
}
