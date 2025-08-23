package httpx

import (
	"net/http"
	"testing"
)

func BenchmarkCompose(b *testing.B) {
	noop := func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFn(func(r *http.Request) (*http.Response, error) {
			return next.RoundTrip(r)
		})
	}

	chain := make([]ChainedRoundTripper, 10)
	for i := range chain {
		chain[i] = noop
	}

	base := RoundTripperFn(func(r *http.Request) (*http.Response, error) {
		return NewResp(200, ""), nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Compose(base, chain...)
	}
}

func BenchmarkCompose100(b *testing.B) {
        noop := func(next http.RoundTripper) http.RoundTripper {
                return RoundTripperFn(func(r *http.Request) (*http.Response, error) {
                        return next.RoundTrip(r)
                })
        }

        chain := make([]ChainedRoundTripper, 100)
        for i := range chain {
                chain[i] = noop
        }

        base := RoundTripperFn(func(r *http.Request) (*http.Response, error) {
                return NewResp(200, ""), nil
        })

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
                _ = Compose(base, chain...)
        }
}
