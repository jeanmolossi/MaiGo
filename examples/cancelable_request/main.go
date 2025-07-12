package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/jeanmolossi/MaiGo/async"
	"github.com/jeanmolossi/MaiGo/examples/testserver"
	"github.com/jeanmolossi/MaiGo/pkg/maigo"
)

func main() {
	ts := testserver.NewManager().
		NewServerBuilder().
		EnableBusy().
		SleepFor(time.Second*5, 0.1).
		Build()

	defer ts.Close()

	ctx, abort := context.WithCancel(context.Background())

	client := maigo.NewClient(ts.URL).Build()

	result, err := async.Dispatch(client.GET("/users").Context().Set(ctx))
	if err != nil {
		slog.Error("error to dispatch request", "error", err)
		return
	}

	time.Sleep(time.Millisecond * 100)

	abort()

	_, err = result.Response()
	if err == nil {
		slog.Error("was expecting an error but <nil> received")
		return
	}

	slog.Info("request aborted", "error", err)
}
