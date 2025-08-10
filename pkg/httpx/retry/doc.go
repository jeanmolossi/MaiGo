// Package retry provides middleware for retrying HTTP client requests.
// It exposes a RoundTripper that can be composed with others to automatically
// retry failed requests with configurable backoff, allowed methods and body
// replay strategies.
package retry
