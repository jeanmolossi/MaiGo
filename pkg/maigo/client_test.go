package maigo

import (
	"errors"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	testtable := []struct {
		name        string
		baseURL     string
		expectError bool
	}{
		{
			name:        "Successful client creation",
			baseURL:     "https://example.com",
			expectError: false,
		},
	}

	for _, tt := range testtable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clientBuilder := NewClient(tt.baseURL)
			if !clientBuilder.client.Validations().IsEmpty() != tt.expectError {
				t.Errorf("NewClient() error = %v, expectedError %v", clientBuilder.client.Validations().Count() > 0, tt.expectError)
			}
		})
	}
}

func TestNewClientLoadBalancer(t *testing.T) {
	t.Parallel()

	baseURLs := []string{
		"https://server1.com",
		"https://server2.com",
		"https://server3.com",
	}

	for _, raw := range baseURLs {
		if raw == "" {
			t.Fatalf("base URL is empty")
		}

		if _, err := url.Parse(raw); err != nil {
			t.Fatalf("invalid base URL %q: %v", raw, err)
		}
	}

	builder := NewClientLoadBalancer(baseURLs)
	if builder == nil {
		t.Fatalf("NewClientLoadBalancer returned nil")
	}

	for i := 0; i < len(baseURLs)*2; i++ {
		want := baseURLs[i%len(baseURLs)]

		base := builder.client.BaseURL()
		if base == nil {
			t.Fatalf("call %d: BaseURL() returned nil", i)
		}

		if got := base.String(); got != want {
			t.Errorf("call %d: BaseURL() = %q, want %q", i, got, want)
		}
	}
}

func TestNewClientLoadBalancer_InvalidBaseURLs(t *testing.T) {
	t.Parallel()

	builder := NewClientLoadBalancer([]string{
		"",
		"http://example.com:invalid",
	})

	if builder.client.Validations().Count() < 2 {
		t.Fatalf("expected at least 2 validations, got %d", builder.client.Validations().Count())
	}

	foundEmpty, foundParse := false, false

	for i := 0; i < builder.client.Validations().Count(); i++ {
		err := builder.client.Validations().Get(i)
		if errors.Is(err, ErrEmptyBaseURL) {
			foundEmpty = true
			continue
		}

		if errors.Is(err, ErrParseURL) {
			foundParse = true
		}
	}

	if !foundEmpty {
		t.Errorf("expected validation error %v not found", ErrEmptyBaseURL)
	}

	if !foundParse {
		t.Errorf("expected validation error %v not found", ErrParseURL)
	}
}

func TestDefaultClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		baseURL string
		wantErr bool
	}{
		{name: "valid", baseURL: "https://example.com", wantErr: false},
		{name: "empty", baseURL: "", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := DefaultClient(tt.baseURL)

			c, ok := client.(*ClientConfigBase)
			if !ok {
				t.Fatalf("DefaultClient(%q) returned %T, want *ClientConfigBase", tt.baseURL, client)
			}

			gotErr := !c.Validations().IsEmpty()

			if gotErr != tt.wantErr {
				t.Errorf("Validations().IsEmpty() = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func BenchmarkNewClient(t *testing.B) {
	url := "https://example.com"
	for range t.N {
		NewClient(url)
	}
}
