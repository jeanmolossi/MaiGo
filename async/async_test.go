package async

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/maigo"
	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
)

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}

func TestAllConcurrencyLimit(t *testing.T) {
	ts := makeDelayedServer(100 * time.Millisecond)
	defer ts.Close()

	client := maigo.NewClient(ts.URL).Build()

	makeBuilders := func() []contracts.RequestBuilder {
		return []contracts.RequestBuilder{
			client.GET("/"),
			client.GET("/"),
			client.GET("/"),
		}
	}

	cases := []struct {
		limit    int
		expected time.Duration
	}{
		{limit: 0, expected: 100 * time.Millisecond},
		{limit: 3, expected: 100 * time.Millisecond},
		{limit: 2, expected: 200 * time.Millisecond},
		{limit: 1, expected: 300 * time.Millisecond},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("limit=%d", tt.limit), func(t *testing.T) {
			start := time.Now()

			group, err := All(tt.limit, makeBuilders()...)
			if err != nil {
				t.Fatalf("All returned error: %v", err)
			}

			group.Wait()

			elapsed := time.Since(start)

			if elapsed < tt.expected || elapsed > tt.expected+100*time.Millisecond {
				t.Fatalf("limit=%d expected around %v got %v", tt.limit, tt.expected, elapsed)
			}
		})
	}
}

func TestGroupResultIndexOutOfRange(t *testing.T) {
	g := &Group{results: make([]result, 1)}

	_, err := g.Result(1)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	expected := fmt.Sprintf("index %d out of range [0-%d]", 1, 0)
	if err.Error() != expected {
		t.Fatalf("expected error %q got %q", expected, err.Error())
	}
}
