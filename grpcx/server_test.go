package grpcx

import (
	"context"
	"testing"

	"github.com/charliego3/mspp/utils"
	"github.com/charliego3/shandler"
	"log/slog"
)

func TestListen(t *testing.T) {
	slog.SetDefault(slog.New(shandler.NewTextHandler(shandler.WithCaller(), shandler.WithTimeFormat(utils.TimeMillis))))
	s := NewServer()
	s.Run(context.Background())
	s.Wait()
}
