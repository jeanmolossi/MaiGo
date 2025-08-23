package maigo

import (
	"net/url"
	"testing"
)

func TestDefaultBaseURL_BaseURL(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}

	d := newDefaultBaseURL(u)

	if got := d.BaseURL(); got != u {
		t.Errorf("BaseURL() = %v, want %v", got, u)
	}
}

func TestBalancedBaseURL_RoundRobin(t *testing.T) {
	t.Parallel()

	raw := []string{
		"https://server1.com",
		"https://server2.com",
		"https://server3.com",
	}

	urls := make([]*url.URL, len(raw))

	for i, r := range raw {
		u, err := url.Parse(r)
		if err != nil {
			t.Fatalf("parse url %d: %v", i, err)
		}

		urls[i] = u
	}

	b := newBalancedBaseURL(urls)

	for i := 0; i < len(urls)*2; i++ {
		want := urls[i%len(urls)]
		if got := b.BaseURL(); got != want {
			t.Errorf("call %d: BaseURL() = %v, want %v", i, got, want)
		}
	}
}
