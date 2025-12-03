//nolint:revive
package main

import (
	"log/slog"
	"time"

	"github.com/jeanmolossi/maigo/examples/testserver"
	"github.com/jeanmolossi/maigo/pkg/maigo"
)

func main() {
	ts := testserver.NewManager().NewServer()
	defer ts.Close()

	client := maigo.DefaultClient(ts.URL)

	newUser := &testserver.User{
		Name:      "John Doe",
		Birthdate: time.Date(1995, 3, 8, 0, 0, 0, 0, time.UTC),
	}

	resp, err := client.POST("/users").
		Body().AsJSON(newUser).
		Send()
	if err != nil {
		slog.Error("error sending the request.", "error", err)
		return
	}

	if resp.Status().Is5xxServerError() {
		defer resp.Body().Close()

		slog.Error("failed to create user, server error.", "status", resp.Status().Text())

		return
	}

	if resp.Status().Is4xxClientError() {
		defer resp.Body().Close()

		slog.Error("failed to create user, client error.", "status", resp.Status().Text())

		return
	}

	slog.Info("user created!", "status", resp.Status().Text())

	var user *testserver.User
	if err := resp.Body().AsJSON(&user); err != nil {
		slog.Error("error parsing response.", "error", err)
		return
	}

	slog.Info("user data", "data", user)
}
