package vocab

// ContextTypeActivityStreams defines the standard ActivityStreams vocabulary.
// https://www.w3.org/TR/activitystreams-core/
const ContextTypeActivityStreams = "https://www.w3.org/ns/activitystreams"

// ContextTypeSecurity describes the standard security vocabulary for the Fediverse.
// https://w3c.github.io/vc-data-integrity/vocab/security/vocabulary.html
const ContextTypeSecurity = "https://w3id.org/security/v1"

// https://joinmastodon.org/ns#
var ContextTypeToot = map[string]any{
	"toot":         "https://joinmastodon.org/ns#",
	"discoverable": "toot:discoverable",
	"indexable":    "toot:indexable",
	"featured": map[string]any{
		"@id":   "http://joinmastodon.org/ns#featured",
		"@type": "@id",
	},
}
