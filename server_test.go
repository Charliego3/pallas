package pallas

import (
	"context"
	"log/slog"
	"testing"

	"github.com/charliego3/pallas/httpx"
	"github.com/charliego3/pallas/testdata"
	"github.com/charliego3/shandler"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer(t *testing.T) {
	slog.SetDefault(slog.New(shandler.NewTextHandler(shandler.WithCaller())))
	server := httpx.NewServer(
		httpx.WithAddr("tcp", "127.0.0.1:8888"),
	)
	server.RegisterService(new(testdata.Greeter))
	require.NoError(t, server.Run(context.Background()))
}
