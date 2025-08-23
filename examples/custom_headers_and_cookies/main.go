package main

import (
	"log/slog"
	"net/http"

	"github.com/jeanmolossi/maigo/examples/testserver"
	"github.com/jeanmolossi/maigo/pkg/maigo"
	"github.com/jeanmolossi/maigo/pkg/maigo/mime"
)

func main() {
	testsrv := testserver.NewManager().
		NewServerBuilder().
		EnableHeaderDebug().
		Build()

	defer testsrv.Close()

	client := maigo.NewClient(testsrv.URL).
		Header().AddUserAgent("MaiGo/1.0").
		Header().Add("X-My-Header", "header-value").
		Cookie().Add(&http.Cookie{
		Name:  "__Secure-session-id",
		Value: "xyz",
	}).Build()

	resp, err := client.GET("/resources").
		Header().AddAccept(mime.JSON).
		Send()
	if err != nil {
		slog.Error("Error sending the request.", "error", err)
	}

	var data []testserver.Resource

	if err := resp.Body().AsJSON(&data); err != nil {
		slog.Error("Error parsing response", "error", err)
		return
	}

	slog.Info("Data received!", "data", data)
}
