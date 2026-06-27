package sigs

// FieldRequestTarget is the "(request-target)" pseudo-header used in the signature base string.
const FieldRequestTarget = "(request-target)"

// FieldCreated is the "(created)" pseudo-header. It is not supported at this time, and will generate an error.
// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#section-2.3
const FieldCreated = "(created)"

// FieldExpires is the "(expires)" pseudo-header. It is not supported at this time, and will generate an error.
// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#section-2.3
const FieldExpires = "(expires)"

// FieldDate is the "date" header field.
const FieldDate = "date"

// FieldDigest is the "digest" header field that validates the request body.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Digest
// https://datatracker.ietf.org/doc/draft-ietf-httpbis-digest-headers/
const FieldDigest = "digest"

// FieldHost is the "host" header field.
const FieldHost = "host"

// Algorithm_HS2019 is the "hs2019" signature algorithm: the only non-deprecated algorithm.
// Unlike the others, its hash and digest functions are not implied by the algorithm name but
// chosen based on the key used (RSA, HMAC, and ECDSA keys are all supported).
const Algorithm_HS2019 = "hs2019"

// Algorithm_RSA_SHA256 is the deprecated "rsa-sha256" signature algorithm.
// Deprecated by the standard because it reveals which hash and digest algorithm is used.
const Algorithm_RSA_SHA256 = "rsa-sha256"

// Algorithm_RSA_SHA512 is the "rsa-sha512" signature algorithm.
const Algorithm_RSA_SHA512 = "rsa-sha512"

// Algorithm_HMAC_SHA256 is the deprecated "hmac-sha256" signature algorithm.
// Deprecated by the standard because it reveals which hash and digest algorithm is used.
const Algorithm_HMAC_SHA256 = "hmac-sha256"

// Algorithm_ECDSA_SHA256 is the deprecated "ecdsa-sha256" signature algorithm.
// Deprecated by the standard because it reveals which hash and digest algorithm is used.
const Algorithm_ECDSA_SHA256 = "ecdsa-sha256"

// Algorithm_HMAC_SHA512 is the "hmac-sha512" signature algorithm.
const Algorithm_HMAC_SHA512 = "hmac-sha512"

// Algorithm_ECDSA_SHA512 is the "ecdsa-sha512" signature algorithm.
const Algorithm_ECDSA_SHA512 = "ecdsa-sha512"
