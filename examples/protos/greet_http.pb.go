// Code generated by protoc-gen-pallas-http. DO NOT EDIT.
//
// proto-gen-pallas-http version: 1.0.0
// protoc version: v4.25.0
// source file: protos/greet.proto

package pb

import (
	context "context"
	httpx "github.com/charliego3/pallas/httpx"
	types "github.com/charliego3/pallas/types"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the pallas package it is being compiled against.
var _ = new(httpx.CallOption)
var _ = new(types.Service)

type GreeterHTTPServer interface {
	SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterHTTPServer(s *httpx.Server, srv GreeterHTTPServer) {
	s.HandleMethod("GET", "/sayHello", _Greeter_SayHello_GET_HTTP_Handler(srv.(types.Service)).(httpx.Handler))
}

func _Greeter_SayHello_GET_HTTP_Handler(srv types.Service) any {
	return httpx.HandlerFunc(func(c *httpx.Context) error {
		req := new(HelloRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		res, err := srv.(GreeterServer).SayHello(c.Context, req)
		if err != nil {
			return err
		}
		return c.Write(res)
	})
}