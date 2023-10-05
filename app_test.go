package mapp

import (
	"context"
	"github.com/charliego3/shandler"
	"log/slog"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewApp(t *testing.T) {
	NewApp().Run(context.Background())
}

func TestDefaultFunc(t *testing.T) {
	var arr []string
	t.Log(append([]string{"first"}, arr...))
}

func TestCheckAddress(t *testing.T) {
	host, port, err := net.SplitHostPort(":8080")
	require.NoError(t, err)
	slog.Info("address", slog.String("host", host), slog.String("port", port))
}

func TestSlog(t *testing.T) {
	slog.SetDefault(slog.New(shandler.NewTextHandler(
		//shandler.WithWriter(os.Stdout),
		shandler.WithCaller(),
	)).WithGroup("g"))
	slog.Info("this is info message",
		slog.Int("sss", 11),
		slog.String("str", "string value"),
		slog.Group("request",
			"method", http.MethodOptions,
			"url", "/robot/details"))
}
