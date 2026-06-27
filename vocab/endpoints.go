package vocab

// `Endpoints` is a json object in the Actor object which maps additional (typically server/domain-wide)
// endpoints which may be useful either for this actor or someone referencing this actor.
// This mapping may be nested inside the actor document as the value or may be a link to a JSON-LD document
// with these properties.
// (from) https://www.w3.org/TR/activitypub/#endpoints

// EndpointFinishMigration is the "finishMigration" endpoint property.
// Proposed account migration extension (DO NOT USE)
// https://github.com/swicg/activitypub-data-portability/issues/56
const EndpointFinishMigration = "finishMigration"

// EndpointOAuthAuthorization is the "oauthAuthorizationEndpoint" endpoint property.
// https://www.w3.org/TR/activitypub/#oauthAuthorizationEndpoint
const EndpointOAuthAuthorization = "oauthAuthorizationEndpoint"

// EndpointOAuthToken is the "oauthTokenEndpoint" endpoint property.
// https://www.w3.org/TR/activitypub/#oauthTokenEndpoint
const EndpointOAuthToken = "oauthTokenEndpoint"

// EndpointOAuthMigration is the "oauthMigrationEndpoint" endpoint property.
// Proposed account migration extension
// https://swicg.github.io/activitypub-data-portability/lola
const EndpointOAuthMigration = "oauthMigrationEndpoint"

// EndpointOAuthMigrationToken is the "oauthMigrationTokenEndpoint" endpoint property.
// Proposed account migration extension
// https://swicg.github.io/activitypub-data-portability/lola
const EndpointOAuthMigrationToken = "oauthMigrationTokenEndpoint"

// EndpointProvideClientKey is the "provideClientKey" endpoint property.
// https://www.w3.org/TR/activitypub/#provideClientKey
const EndpointProvideClientKey = "provideClientKey"

// EndpointProxyURL is the "proxyUrl" endpoint property.
// https://www.w3.org/TR/activitypub/#proxyUrl
const EndpointProxyURL = "proxyUrl"

// EndpointSignClientKey is the "signClientKey" endpoint property.
// https://www.w3.org/TR/activitypub/#signClientKey
const EndpointSignClientKey = "signClientKey"

// EndpointSharedInbox is the "sharedInbox" endpoint property.
// https://www.w3.org/TR/activitypub/#sharedInbox
const EndpointSharedInbox = "sharedInbox"

// EndpointStartMigration is the "startMigration" endpoint property.
// Proposed account migration extension (DO NOT USE)
// https://github.com/swicg/activitypub-data-portability/issues/56
const EndpointStartMigration = "startMigration"
