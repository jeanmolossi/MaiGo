package maigo

import (
	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
	"github.com/jeanmolossi/MaiGo/pkg/maigo/method"
)

type Request struct {
	client contracts.ClientConfig
	config *RequestConfigBase
}

func newRequest(client contracts.ClientConfig, method method.Type, path string) *RequestBuilder {
	return &RequestBuilder{
		request: &Request{
			client: client,
			config: newRequestConfigBase(method, path),
		},
	}
}
