package api

import "net/http"

type RetryOptions struct {
	count    int
	err      error
	response *http.Response
}

// retryStrategy passes in details about the last response that failed
// the strategy should return true if another attempt should be made
// the strategy should return false if no additional attempts should be made
type retryStrategy func(retry RetryOptions) bool
