package maigo

import (
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

func BenchmarkNewClient(t *testing.B) {
	url := "https://example.com"
	for range t.N {
		NewClient(url)
	}
}
