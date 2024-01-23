package grpcx

import (
	"context"
	"github.com/charliego3/mspp/testdata"
	"github.com/stretchr/testify/require"
	"testing"

	"log/slog"

	"github.com/charliego3/mspp/utility"
	"github.com/charliego3/shandler"
)

func TestListen(t *testing.T) {
	slog.SetDefault(slog.New(shandler.NewTextHandler(shandler.WithCaller(), shandler.WithTimeFormat(utility.TimeMillis))))
	s := NewServer(WithAddr("tcp", "127.0.0.1:9999"))
	s.RegisterService(new(testdata.Greeter))
	err := s.Run(context.Background())
	require.NoError(t, err)
}
