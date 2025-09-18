package maigo

import (
	"errors"
	"testing"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
	"github.com/jeanmolossi/maigo/pkg/maigo/method"
)

func TestClientConfigBase_HTTPMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(*ClientConfigBase, string) contracts.RequestBuilder
		want method.Type
	}{
		{"GET", (*ClientConfigBase).GET, method.GET},
		{"POST", (*ClientConfigBase).POST, method.POST},
		{"PUT", (*ClientConfigBase).PUT, method.PUT},
		{"DELETE", (*ClientConfigBase).DELETE, method.DELETE},
		{"PATCH", (*ClientConfigBase).PATCH, method.PATCH},
		{"HEAD", (*ClientConfigBase).HEAD, method.HEAD},
		{"CONNECT", (*ClientConfigBase).CONNECT, method.CONNECT},
		{"OPTIONS", (*ClientConfigBase).OPTIONS, method.OPTIONS},
		{"TRACE", (*ClientConfigBase).TRACE, method.TRACE},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newClientConfigBase("https://example.com")

			rbContract := tt.call(c, "/path")

			rb, ok := rbContract.(*RequestBuilder)
			if !ok {
				t.Fatalf("%s: expected *RequestBuilder, got %T", tt.name, rbContract)
			}

			if rb.request.config.Method() != tt.want {
				t.Errorf("Method() = %s, want %s", rb.request.config.Method(), tt.want)
			}

			if rb.request.config.Path() != "/path" {
				t.Errorf("Path() = %s, want /path", rb.request.config.Path())
			}
		})
	}
}

func TestNewClientConfigBase_Validations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		baseURL string
		wantErr error
	}{
		{"valid", "https://example.com", nil},
		{"empty", "", ErrEmptyBaseURL},
		{"invalid", "http://example.com:invalid", ErrParseURL},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newClientConfigBase(tt.baseURL)

			if tt.wantErr == nil {
				if !c.Validations().IsEmpty() {
					t.Fatalf("unexpected validation errors: %v", c.Validations().Unwrap())
				}

				return
			}

			if c.Validations().IsEmpty() {
				t.Fatalf("expected validation error, got none")
			}

			found := false

			for i := 0; i < c.Validations().Count(); i++ {
				if errors.Is(c.Validations().Get(i), tt.wantErr) {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected validation error %v not found", tt.wantErr)
			}
		})
	}
}

func TestNewBalancedClientConfigBase_Validations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		baseURLs []string
		wantErrs []error
	}{
		{
			name:     "empty slice",
			baseURLs: nil,
			wantErrs: []error{ErrEmptyBaseURL},
		},
		{
			name:     "with invalid and empty",
			baseURLs: []string{"", "http://example.com:invalid"},
			wantErrs: []error{ErrEmptyBaseURL, ErrParseURL},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newBalancedClientConfigBase(tt.baseURLs)

			if c.Validations().Count() < len(tt.wantErrs) {
				t.Fatalf("expected at least %d validations, got %d", len(tt.wantErrs), c.Validations().Count())
			}

			for _, we := range tt.wantErrs {
				found := false

				for i := 0; i < c.Validations().Count(); i++ {
					if errors.Is(c.Validations().Get(i), we) {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("expected error %v not found", we)
				}
			}
		})
	}
}

func TestNewBalancedClientConfigBase_RoundRobin(t *testing.T) {
	t.Parallel()

	baseURLs := []string{
		"https://server1.com",
		"https://server2.com",
		"https://server3.com",
	}

	c := newBalancedClientConfigBase(baseURLs)

	for i := 0; i < len(baseURLs)*2; i++ {
		want := baseURLs[i%len(baseURLs)]
		if got := c.BaseURL().String(); got != want {
			t.Errorf("call %d: BaseURL() = %q, want %q", i, got, want)
		}
	}
}
