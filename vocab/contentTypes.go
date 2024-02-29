package vocab

// Accept is the string used in the HTTP header to request a response be encoded as a MIME type
const Accept = "Accept"

// ContentType is the string used in the HTTP header to designate a MIME type
const ContentType = "Content-Type"

// ContentTypeActivityPub is the standard MIME type for ActivityPub content
const ContentTypeActivityPub = "application/activity+json"

// ContentTypeJSONLD is the standard MIME Type for JSON-LD content
// https://en.wikipedia.org/wiki/JSON-LD
const ContentTypeJSONLD = "application/ld+json"

// ContentTypeJSONLDWithProfile is the standard MIME Type for JSON-LD content, with profile
// to designate ActivityPub content.
// https://www.w3.org/TR/activitystreams-core/#media-type
const ContentTypeJSONLDWithProfile = `application/ld+json; profile="https://www.w3.org/ns/activitystreams"`

// ContentTypeHTML is the standard MIME type for HTML content
const ContentTypeHTML = "text/html"

// ContentTypeJSON is the standard MIME Type for JSON content
const ContentTypeJSON = "application/json"

// ContentTypeJSONResourceDescriptor is the standard MIME Type for JSON Resource Descriptor content
// which is used by WebFinger: https://datatracker.ietf.org/doc/html/rfc7033#section-10.2
const ContentTypeJSONResourceDescriptor = "application/jrd+json"
