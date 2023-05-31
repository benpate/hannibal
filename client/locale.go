package client

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

type Locale struct {
	locales     []string
	innerClient streams.Client
}

func NewLocale(locales []string, innerClient streams.Client) Locale {
	return Locale{
		locales:     locales,
		innerClient: innerClient,
	}
}

func (locale Locale) Load(uri string) (streams.Document, error) {
	result, err := locale.innerClient.Load(uri)

	if err != nil {
		return result, derp.Wrap(err, "hannibal.client.Locale.Load", "Error loading document", uri)
	}

	result.WithOptions(
		streams.WithClient(locale),
		streams.WithLocales(locale.locales...),
	)

	return result, nil
}
