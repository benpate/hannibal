package queue

import "github.com/benpate/derp"

func IsServerError(err error) bool {
	code := derp.ErrorCode(err)
	return code >= 500 && code < 600
}
