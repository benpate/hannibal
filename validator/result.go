package validator

// Result indicates the outcome of a validation attempt
type Result string

// ResultValid indicates that the current Validator has successfully validated the HTTP request
const ResultValid Result = "VALID"

// ResultInvalid indicates that the current Validator can say with certainty that the HTTP request is invalid
const ResultInvalid Result = "INVALID"

// ResultUnknown indicates that the current Validator cannot say that the HTTP request is valid or invalid
const ResultUnknown Result = "UNKNOWN"

// ResultError indicates that the current Validator encountered an error while attempting to validate the HTTP request
const ResultError Result = "ERROR"
