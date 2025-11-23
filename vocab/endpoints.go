package vocab

// `Endpoints` is a json object in the Actor object which maps additional (typically server/domain-wide)
// endpoints which may be useful either for this actor or someone referencing this actor.
// This mapping may be nested inside the actor document as the value or may be a link to a JSON-LD document
// with these properties.
// (from) https://www.w3.org/TR/activitypub/#endpoints

// https://www.w3.org/TR/activitypub/#oauthAuthorizationEndpoint
const EndpointOAuthAuthorization = "oauthAuthorizationEndpoint"

// Proposed account migration extension
// https://swicg.github.io/activitypub-data-portability/lola
const EndpointOAuthMigration = "oauthMigrationEndpoint"

// https://www.w3.org/TR/activitypub/#oauthTokenEndpoint
const EndpointOAuthToken = "oauthTokenEndpoint"

// https://www.w3.org/TR/activitypub/#provideClientKey
const EndpointProvideClientKey = "provideClientKey"

// https://www.w3.org/TR/activitypub/#proxyUrl
const EndpointProxyURL = "proxyUrl"

// https://www.w3.org/TR/activitypub/#signClientKey
const EndpointSignClientKey = "signClientKey"

// https://www.w3.org/TR/activitypub/#sharedInbox
const EndpointSharedInbox = "sharedInbox"
