package grpcx

import (
	"context"
	"testing"

	"log/slog"

	"github.com/charliego3/mspp/utils"
	"github.com/charliego3/shandler"
)

func TestListen(t *testing.T) {
	slog.SetDefault(slog.New(shandler.NewTextHandler(shandler.WithCaller(), shandler.WithTimeFormat(utils.TimeMillis))))
	s := NewServer()
	s.Start(context.Background())
}
