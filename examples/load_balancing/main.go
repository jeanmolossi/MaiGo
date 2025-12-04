// Package main showcases load balancing between multiple endpoints.
package main

import (
	"log/slog"

	"github.com/jeanmolossi/maigo/examples/testserver"
	"github.com/jeanmolossi/maigo/pkg/maigo"
	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

func main() {
	serverManager := testserver.NewManager()

	ts1 := serverManager.NewServer()
	defer ts1.Close()

	ts2 := serverManager.NewServerBuilder().EnableBusy().Build()
	defer ts2.Close()

	ts3 := serverManager.NewServer()
	defer ts3.Close()

	client := maigo.NewClientLoadBalancer([]string{
		ts1.URL,
		ts2.URL,
		ts3.URL,
	}).Build()

	for range 100 {
		healthcheck(client)
	}
}

func healthcheck(client contracts.ClientHTTPMethods) {
	resp, err := client.GET("/").Send()
	if err != nil {
		slog.Error("error sending the request", "error", err)
		return
	}

	if resp.Status().Is5xxServerError() {
		slog.Error("healthcheck failed.", "status", resp.Status().Text())
		return
	}

	var healthResponse map[string]any

	if err := resp.Body().AsJSON(&healthResponse); err != nil {
		slog.Error("error parsing response.", "error", err)
		return
	}

	slog.Info("healthcheck received.", "data", healthResponse)
}
