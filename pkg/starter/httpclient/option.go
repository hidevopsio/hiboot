package httpclient

import (
	"time"
)

// Option represents the client options
type Option func(*client)

// WithHTTPTimeout sets hystrix timeout
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *client) {
		c.timeout = timeout
	}
}

// WithRetryCount sets the retrier count for the hystrixHTTPClient
func WithRetryCount(retryCount int) Option {
	return func(c *client) {
		c.retryCount = retryCount
	}
}

// WithRetrier sets the strategy for retrying
func WithRetrier(retrier Retriable) Option {
	return func(c *client) {
		c.retrier = retrier
	}
}

// WithHTTPClient sets a custom http client
func WithHTTPClient(doer Doer) Option {
	return func(c *client) {
		c.client = doer
	}
}
