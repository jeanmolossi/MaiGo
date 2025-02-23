package maigo

import "github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"

var _ contracts.Builder[contracts.ClientHTTPMethods] = (*ClientBuilder)(nil)

type ClientBuilder struct {
	client contracts.Client
}

func NewClient(baseURL string) *ClientBuilder {
	return &ClientBuilder{
		client: newClientConfigBase(baseURL),
	}
}

func DefaultClient(baseURL string) contracts.ClientHTTPMethods {
	return NewClient(baseURL).Build()
}

// Build implements contracts.Builder.
func (h *ClientBuilder) Build() contracts.ClientHTTPMethods {
	return h.client
}
