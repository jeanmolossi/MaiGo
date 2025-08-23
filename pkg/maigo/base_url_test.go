package maigo

import (
	"net/url"
	"sync"
	"testing"
)

func mustParse(t *testing.T, raw string) *url.URL {
	t.Helper()

	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse %q: %v", raw, err)
	}

	return u
}

func TestDefaultBaseURL_BaseURL(t *testing.T) {
	t.Parallel()

	u := mustParse(t, "https://example.com")

	d := newDefaultBaseURL(u)

	if got := d.BaseURL(); got.String() != u.String() {
		t.Errorf("BaseURL() = %q, want %q", got.String(), u.String())
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
		urls[i] = mustParse(t, r)
	}

	b := newBalancedBaseURL(urls)

	for i := 0; i < len(urls)*2; i++ {
		want := raw[i%len(raw)]
		got := b.BaseURL()

		if got == nil || got.String() != want {
			t.Errorf("call %d: BaseURL() = %q, want %q", i, got.String(), want)
		}
	}
}

func TestBalancedBaseURL_RoundRobin_Concurrent(t *testing.T) {
	t.Parallel()

	raw := []string{
		"https://server1.com",
		"https://server2.com",
		"https://server3.com",
	}

	urls := make([]*url.URL, len(raw))
	for i, r := range raw {
		urls[i] = mustParse(t, r)
	}

	b := newBalancedBaseURL(urls)

	const workers = 300

	counts := make([]int, len(urls))

	var wg sync.WaitGroup

	var mu sync.Mutex

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			u := b.BaseURL()
			if u == nil {
				t.Error("BaseURL returned nil")
				return
			}

			s := u.String()
			for idx, r := range raw {
				if s == r {
					mu.Lock()
					counts[idx]++
					mu.Unlock()

					return
				}
			}

			t.Errorf("unexpected URL %q", s)
		}()
	}

	wg.Wait()

	min, max := counts[0], counts[0]
	for _, c := range counts[1:] {
		if c < min {
			min = c
		}

		if c > max {
			max = c
		}
	}

	if max-min > 1 {
		t.Errorf("imbalanced distribution: %v", counts)
	}
}

func TestBalancedBaseURL_SingleURL(t *testing.T) {
	t.Parallel()

	u := mustParse(t, "https://only.com")

	b := newBalancedBaseURL([]*url.URL{u})

	for i := 0; i < 10; i++ {
		if got := b.BaseURL(); got != u {
			t.Fatalf("sequential call %d: got %v, want %v", i, got, u)
		}
	}

	const workers = 50

	var wg sync.WaitGroup

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			if got := b.BaseURL(); got != u {
				t.Errorf("concurrent call: got %v, want %v", got, u)
			}
		}()
	}

	wg.Wait()
}

func TestBalancedBaseURL_EmptyURLs(t *testing.T) {
	t.Parallel()

	b := newBalancedBaseURL(nil)

	for i := 0; i < 10; i++ {
		if got := b.BaseURL(); got != nil {
			t.Fatalf("sequential call %d: got %v, want nil", i, got)
		}
	}

	const workers = 50

	var wg sync.WaitGroup

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			if got := b.BaseURL(); got != nil {
				t.Errorf("concurrent call: got %v, want nil", got)
			}
		}()
	}

	wg.Wait()
}
