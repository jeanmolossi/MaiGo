// Package main demonstrates asynchronous dispatching with MaiGo.
package main

import (
	"log/slog"
	"time"

	"github.com/jeanmolossi/maigo/async"
	"github.com/jeanmolossi/maigo/examples/testserver"
	"github.com/jeanmolossi/maigo/pkg/maigo"
	"github.com/jeanmolossi/maigo/pkg/maigo/header"
)

func main() {
	ts := testserver.
		NewManager().
		NewServerBuilder().
		EnableBusy().
		EnableHeaderDebug().
		SleepFor(time.Second, 1.2).
		Build()
	defer ts.Close()

	client := maigo.NewClient(ts.URL).
		Build()

	start := time.Now()
	// Dispatch the request asynchronously without block current scope
	result, err := async.Dispatch(
		client.GET("/users").
			Retry().SetConstantBackoff(time.Millisecond*100, 5).
			Header().Add(header.Type("X-Request-Id"), "123"),
	)
	if err != nil {
		slog.Error("failed to dispatch request", "error", err)
		return
	}

	// Simulate some stuff while the request is at another thread
	slog.Info("waiting request to resolve", "since", time.Since(start))
	slog.Info("doing some stuff", "since", time.Since(start))
	time.Sleep(time.Millisecond * 200)
	slog.Info("some stuff while request waiting resolution", "since", time.Since(start))

	// Catch the response of the request
	resp, err := result.Response()
	if err != nil {
		slog.Error("request error", "error", err)
		return
	}

	slog.Info("request resolved", "duration", time.Since(start))

	var users []testserver.User

	if err := resp.Body().AsJSON(&users); err != nil {
		slog.Error("failed to parse users", "error", err)
		return
	}

	slog.Info("data received", "data", users)
}
