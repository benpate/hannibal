package queue

import (
	"strconv"
	"time"

	"github.com/benpate/derp"
)

type TaskFunc func() error
type BackoffFunc func(int) time.Duration

func Retry(taskFunc TaskFunc, options ...RetryOption) {

	retryConfig := NewRetryConfig(options...)

	var lastError error

	for attempt := 0; attempt < retryConfig.MaxAttempts; attempt++ {

		time.Sleep(retryConfig.Backoff(attempt))

		lastError = taskFunc()

		if lastError == nil {
			return
		}

		if retryConfig.ReportErrors {
			derp.Report(derp.Wrap(lastError, "outbox.Retry", "Attempt: "+strconv.Itoa(attempt)))
		}
	}

	if retryConfig.ReportErrors {
		derp.Report(derp.Wrap(lastError, "outbox.Retry", "Max Attempts Exceeded"))
	}
}
