package main

import (
	"log/slog"
	"time"

	"github.com/jeanmolossi/maigo/examples/testserver"
	"github.com/jeanmolossi/maigo/pkg/maigo"
)

func main() {
	ts := testserver.NewManager().
		NewServerBuilder().
		EnableBusy().
		Build()

	defer ts.Close()

	client := maigo.DefaultClient(ts.URL)

	slog.Info("get resource:", "id", 2)

	res, err := client.GET("/resources/2").
		Retry().SetExponentialBackoff(time.Millisecond*50, 10, 2).
		Retry().WithMaxDelay(time.Millisecond * 500).
		Send()
	if err != nil {
		slog.Error("error sending the request.", "error", err)
		return
	}

	if res.Status().IsError() {
		defer res.Body().Close()

		slog.Error("failed to fetch a resource.", "status", res.Status().Text())

		return
	}

	var resource *testserver.Resource

	if err := res.Body().AsJSON(&resource); err != nil {
		slog.Error("error parsing response.", "error", err)
		return
	}

	slog.Info("resource got!", "data", resource)
}
