package maigo

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

func TestClientConfigBuilder_SetTLSConfig_NewTransport(t *testing.T) {
	t.Parallel()

	builder := NewClient("https://example.com")
	tlsConfig := &tls.Config{ServerName: "example.com"}

	builder.Config().SetTLSConfig(tlsConfig)

	client := builder.Build().(contracts.ClientCompat)

	transport, ok := client.HttpClient().Transport().(*http.Transport)
	if !ok {
		t.Fatalf("transport is not *http.Transport: %T", client.HttpClient().Transport())
	}

	if transport.TLSClientConfig != tlsConfig {
		t.Fatalf("TLSClientConfig = %v, want %v", transport.TLSClientConfig, tlsConfig)
	}
}

func TestClientConfigBuilder_SetTLSConfig_KeepsExistingTransport(t *testing.T) {
	t.Parallel()

	builder := NewClient("https://example.com")
	customTransport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	builder.Config().SetCustomTransport(customTransport)
	builder.Config().SetTLSConfig(tlsConfig)

	client := builder.Build().(contracts.ClientCompat)

	transport, ok := client.HttpClient().Transport().(*http.Transport)
	if !ok {
		t.Fatalf("transport is not *http.Transport: %T", client.HttpClient().Transport())
	}

	if transport != customTransport {
		t.Fatalf("transport was replaced; got %p, want %p", transport, customTransport)
	}

	if transport.TLSClientConfig != tlsConfig {
		t.Fatalf("TLSClientConfig = %v, want %v", transport.TLSClientConfig, tlsConfig)
	}
}
