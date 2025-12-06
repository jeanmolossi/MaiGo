package maigo

import (
	"testing"

	"github.com/jeanmolossi/maigo/pkg/maigo/header"
	"github.com/jeanmolossi/maigo/pkg/maigo/mime"
)

func TestRequestBuilderUnwrap(t *testing.T) {
	t.Parallel()

	builder := NewClient("https://example.com")
	builder.Header().Add(header.Accept, mime.JSON.String())

	client := builder.Build()

	reqBuilder := client.GET("/users")
	reqBuilder.Query().AddParam("active", "1").Header().AddUserAgent("maigo")

	req, err := reqBuilder.Unwrap()
	if err != nil {
		t.Fatalf("Unwrap() error = %v", err)
	}

	const expectedURL = "https://example.com/users?active=1"
	if got := req.URL.String(); got != expectedURL {
		t.Fatalf("req.URL = %q, want %q", got, expectedURL)
	}

	if got := req.Header.Get(header.Accept.String()); got != mime.JSON.String() {
		t.Fatalf("Accept header = %q, want %q", got, mime.JSON.String())
	}

	if got := req.Header.Get(header.UserAgent.String()); got != "maigo" {
		t.Fatalf("User-Agent header = %q, want %q", got, "maigo")
	}
}
