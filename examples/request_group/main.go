// Package main groups multiple requests executed together.
package main

import (
	"log/slog"
	"time"

	"github.com/jeanmolossi/maigo/async"
	"github.com/jeanmolossi/maigo/examples/testserver"
	"github.com/jeanmolossi/maigo/pkg/maigo"
)

func main() {
	ts := testserver.NewManager().
		NewServerBuilder().
		SleepFor(time.Millisecond*300, 1.5).
		Build()

	defer ts.Close()

	client := maigo.NewClient(ts.URL).Build()

	start := time.Now()

	group, err := async.All(2,
		client.GET("/users"),
		client.GET("/resources"),
	)
	if err != nil {
		slog.Error("error building request group", "error", err)
		return
	}

	slog.Info("waiting group resolve", "until now", time.Since(start))

	group.Wait()

	slog.Info("group resolved", "resolved in", time.Since(start))

	resp, err := group.Result(0)
	if err != nil {
		slog.Error("failed to get users", "error", err)
		return
	}

	var users []any
	if err := resp.Body().AsJSON(&users); err != nil {
		slog.Error("failed to read json result", "error", err)
		return
	}

	resp, err = group.Result(1)
	if err != nil {
		slog.Error("failed to get resources", "error", err)
		return
	}

	var resources []any
	if err := resp.Body().AsJSON(&resources); err != nil {
		slog.Error("failed to read json result", "error", err)
		return
	}

	slog.Info("requests was made", "users", users, "resources", resources)
}
