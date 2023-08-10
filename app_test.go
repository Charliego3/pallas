package mapp

import (
	"context"
	"golang.org/x/exp/slog"
	"net"
	"net/http"
	"os"
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
	slog.Info("Host: %s, Port: %s", host, port)
}

func TestSlog(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})).WithGroup("Group")
	logger.Info("this is info message", slog.Int("", 11), slog.String("str", "string value"), slog.Group("request",
		"method", http.MethodOptions,
		"url", "/robot/details"))
}
