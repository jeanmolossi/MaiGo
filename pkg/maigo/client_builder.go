package maigo

import "github.com/jeanmolossi/maigo/pkg/maigo/contracts"

var _ contracts.Builder[contracts.ClientHTTPMethods] = (*ClientBuilder)(nil)

type ClientBuilder struct {
	client contracts.Client
}

func NewClient(baseURL string) *ClientBuilder {
	return &ClientBuilder{
		client: newClientConfigBase(baseURL),
	}
}

func NewClientLoadBalancer(baseURLs []string) *ClientBuilder {
	return &ClientBuilder{
		client: newBalancedClientConfigBase(baseURLs),
	}
}

func DefaultClient(baseURL string) contracts.ClientHTTPMethods {
	return NewClient(baseURL).Build()
}

// Build implements contracts.Builder.
func (b *ClientBuilder) Build() contracts.ClientHTTPMethods {
	return b.client
}
