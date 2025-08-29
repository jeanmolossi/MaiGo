package maigo

import (
	"errors"
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

	builder := NewClientLoadBalancer(baseURLs)

	for i := 0; i < len(baseURLs)*2; i++ {
		want := baseURLs[i%len(baseURLs)]
		if got := builder.client.BaseURL().String(); got != want {
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

	if builder.client.Validations().Count() != 2 {
		t.Fatalf("expected 2 validations, got %d", builder.client.Validations().Count())
	}

	if err := builder.client.Validations().Get(0); !errors.Is(err, ErrEmptyBaseURL) {
		t.Errorf("first validation = %v, want ErrEmptyBaseURL", err)
	}

	if err := builder.client.Validations().Get(1); !errors.Is(err, ErrParseURL) {
		t.Errorf("second validation = %v, want ErrParseURL", err)
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

			c := DefaultClient(tt.baseURL).(*ClientConfigBase)
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
