package grpcx

import (
	"context"
	"testing"

	"log/slog"

	"github.com/charliego3/mspp/utility"
	"github.com/charliego3/shandler"
)

func TestListen(t *testing.T) {
	slog.SetDefault(slog.New(shandler.NewTextHandler(shandler.WithCaller(), shandler.WithTimeFormat(utility.TimeMillis))))
	s := NewServer()
	s.Run(context.Background())
}
