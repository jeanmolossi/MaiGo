// Package main demonstrates how to enable custom TLS settings on the MaiGo client.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"github.com/jeanmolossi/maigo/pkg/maigo"
)

func main() {
	// Create a HTTPS test server with a self-signed certificate.
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "secure ok")
	}))
	defer server.Close()

	// Build a cert pool with the server's certificate so the client trusts it.
	certPool := x509.NewCertPool()
	certPool.AddCert(server.Certificate())

	tlsConfig := &tls.Config{RootCAs: certPool}

	client := maigo.NewClient(server.URL).
		Config().
		SetTLSConfig(tlsConfig).
		Build()

	resp, err := client.GET("/").Send()
	if err != nil {
		slog.Error("request failed", "error", err)
		return
	}

	body, err := resp.Body().AsString()
	if err != nil {
		slog.Error("failed to read body", "error", err)
		return
	}

	slog.Info("response received", "status", resp.Status().Code(), "body", body)
}
