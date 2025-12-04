// Package main shows a basic GET request with MaiGo.
package main

import (
	"log/slog"

	"github.com/jeanmolossi/maigo/examples/testserver"
	"github.com/jeanmolossi/maigo/pkg/maigo"
	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

func main() {
	// Start test server.
	ts := testserver.NewManager().NewServer()
	defer ts.Close()

	// Creates a default client.
	client := maigo.DefaultClient(ts.URL)

	// Get all users.
	getUsers(client)

	// Get a user by id.
	getUser(client, "1")

	// Get a user that does not exists.
	getUser(client, "99")
}

func getUsers(c contracts.ClientHTTPMethods) {
	slog.Info("Get all users.")

	resp, err := c.GET("/users").Send()
	if err != nil {
		slog.Error("error getting response", "error", err)
		return
	}

	handleResponse(resp, &[]testserver.User{})
}

func getUser(c contracts.ClientHTTPMethods, id string) {}

func handleResponse(resp contracts.Response, data any) {
	slog.Info("Response:", "status", resp.Status().Text())

	if resp.Status().IsError() {
		slog.Error("Failed to get data.")
		return
	}

	if err := resp.Body().AsJSON(data); err != nil {
		slog.Error("Error parsing response.", "error", err)
		return
	}

	slog.Info("Data received!", "data", data)
}
