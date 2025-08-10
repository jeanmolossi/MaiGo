package httpx

import "net/http"

// RoundTripperFn allows to create round trippers with a simple func.
type RoundTripperFn func(*http.Request) (*http.Response, error)

// Compile time check if [RoundTripperFn] implements [http.RoundTripper].
var _ http.RoundTripper = (*RoundTripperFn)(nil)

func (r RoundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

// ChainedRoundTripper is our handler chained style.
type ChainedRoundTripper func(next http.RoundTripper) http.RoundTripper

// Compose apply chained round trippers over a base round tripper (usually [http.DefaultTransport]).
// The received order: Compose(base, A, B, C) => request goes down by A -> B -> C -> base.
func Compose(base http.RoundTripper, chain ...ChainedRoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}

	for i := len(chain) - 1; i >= 0; i-- {
		base = chain[i](base)
	}

	return base
}
