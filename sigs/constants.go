package sigs

const FieldRequestTarget = "(request-target)"

// FieldExpires is not supported at this time, and will generate an error.
// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#section-2.3
// const FieldExpires = "(expires)"

// FieldCreated is not supported at this time, and will generate an error.
// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#section-2.3
// const FieldCreated = "(created)"

const FieldHost = "host"

const FieldDate = "date"

// FieldDigest represents the Digest header field that validates the request body.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Digest
// https://datatracker.ietf.org/doc/draft-ietf-httpbis-digest-headers/
const FieldDigest = "digest"

// The “hs2019” signature algorithm. This is the only non-deprecated algorithm. Unlike the other algorithms, the hash and digest functions are not implied by the choice of this signature algorithm. Instead, the hash and digest functions are chosen based on the key used. RSA, HMAC, and ECDSA keys are all supported.
// TODO: How to implement hs2019?
const Algorithm_HS2019 = "hs2019"

// Deprecated. The “rsa-sha256” signature algorithm. Deprecated by the standard because it reveals which hash and digest algorithm is used.
const Algorithm_RSA_SHA256 = "rsa-sha256"

const Algorithm_RSA_SHA512 = "rsa-sha512"

// Deprecated. The “hmac-sha256” signature algorithm. Deprecated by the standard because it reveals which hash and digest algorithm is used.
const Algorithm_HMAC_SHA256 = "hmac-sha256"

// Deprecated. The “ecdsa-sha256” signature algorithm. Deprecated by the standard because it reveals which hash and digest algorithm is used.
const Algorithm_ECDSA_SHA256 = "ecdsa-sha256"

// TODO: Are these supported by the actual specs?
const Algorithm_HMAC_SHA512 = "hmac-sha512"
const Algorithm_ECDSA_SHA512 = "ecdsa-sha512"

// DigestSha256 represents SHA-256 digest algorithm
const Digest_SHA256 = "sha-256"

// DigestSha512 represents SHA-512 digest algorithm
const Digest_SHA512 = "sha-512"

// Additional algorithms specified by https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Digest
// unixsum, unixcksum, crc32c, sha-256 and sha-512, id-sha-256, id-sha-512
