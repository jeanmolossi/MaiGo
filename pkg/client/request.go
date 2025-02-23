package client

import (
	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
	"github.com/jeanmolossi/MaiGo/pkg/client/method"
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
