package main

import (
	"context"
	"log/slog"

	"github.com/charliego3/pallas"
	"github.com/charliego3/pallas/testdata"
	"github.com/charliego3/shandler"
)

func init() {
	slog.SetDefault(slog.New(shandler.NewTextHandler(shandler.WithCaller())))
}

func main() {
	app := pallas.NewApp(
		pallas.WithTCPAddr("127.0.0.1:50051"),
		pallas.WithLogger(shandler.CopyWithPrefix("Application")),
		pallas.WithName("Example"),
	)

	app.RegisterService(new(testdata.Greeter))
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
