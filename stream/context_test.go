package stream

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {

	// Test default context data
	{
		c := DefaultContext()

		require.True(t, c.IsEmptyTail())
		require.True(t, c.Tail().IsEmpty())
		head := c.Head()
		require.Equal(t, "https://www.w3.org/ns/activitystreams", head.Vocabulary)
		require.Equal(t, "und", head.Language)
		require.Zero(t, len(head.Extensions))

		result, err := c.MarshalJSON()
		require.Nil(t, err)
		require.Equal(t, `"https://www.w3.org/ns/activitystreams"`, string(result))
	}

	// Test custom context, and chaining multiple contexts
	{
		c := NewContext()
		entry := c.Add("https://test.com").WithLanguage("en-us")

		require.Equal(t, "https://test.com", c.Head().Vocabulary)
		require.Equal(t, "en-us", c.Head().Language)
		require.Zero(t, len(c.Head().Extensions))

		{
			result, err := json.Marshal(c)
			require.Nil(t, err)
			require.Equal(t, `{"@language":"en-us","@vocab":"https://test.com"}`, string(result))
		}

		entry.WithExtension("ext", "https://extension.com/ns/activitystreams")

		json1, err1 := c.MarshalJSON()
		require.Nil(t, err1)
		require.Equal(t, `{"@language":"en-us","@vocab":"https://test.com","ext":"https://extension.com/ns/activitystreams"}`, string(json1))

		c.Add("https://www.w3.org/ns/activitystreams")
		json2, err2 := c.MarshalJSON()

		require.Equal(t, `[{"@language":"en-us","@vocab":"https://test.com","ext":"https://extension.com/ns/activitystreams"},"https://www.w3.org/ns/activitystreams"]`, string(json2))
		require.Nil(t, err2)
	}

	// Test safely adding an extension to an improperly initialized context
	{
		c := NewContext()
		c.Add("https://test.com").
			WithExtension("dog", "https://dog.com/ns/activitystreams")

		require.Equal(t, "https://test.com", c.Head().Vocabulary)
		require.Equal(t, "und", c.Head().Language)
		require.Equal(t, c.Head().Extensions["dog"], "https://dog.com/ns/activitystreams")
	}
}
