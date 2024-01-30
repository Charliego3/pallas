package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/charliego3/pallas/testdata"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestGrpcClient(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	greeter := testdata.NewGreeterClient(conn)
	resp, err := greeter.SayHello(context.Background(), &testdata.HelloRequest{
		Name: "client",
	}, grpc.Header(&metadata.MD{
		"Content-Type": []string{"application/grpc"},
	}))
	require.NoError(t, err)
	fmt.Println(resp.String())
}
