package maigo

import "github.com/jeanmolossi/maigo/pkg/maigo/contracts"

var _ contracts.Builder[contracts.ClientHTTPMethods] = (*ClientBuilder)(nil)

type ClientBuilder struct {
	client contracts.ClientCompat
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

func DefaultClientCompat(baseURL string) contracts.ClientCompat {
	return NewClient(baseURL).client
}

// Build implements contracts.Builder.
func (b *ClientBuilder) Build() contracts.ClientHTTPMethods {
	return b.client
}
