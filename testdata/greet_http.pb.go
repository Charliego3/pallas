package testdata

import (
	"context"
	"net/http"

	"github.com/charliego3/pallas/httpx"
	"github.com/charliego3/pallas/types"
)

type GreeterHTTPClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...httpx.CallOption) (*HelloReply, error)
}

type greeterHTTPClient struct{}

func (g *greeterHTTPClient) SayHello(ctx context.Context, in *HelloRequest, opts ...httpx.CallOption) (*HelloReply, error) {
	return nil, nil
}

type GreeterHttpServer interface {
	SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error)
}

func _Greeter_SayHello_HTTP_Handler(srv types.Service) any {
	return httpx.HandlerFunc(func(c *httpx.Context) error {
		req := new(HelloRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		res, err := srv.(GreeterServer).SayHello(c.Context, req)
		if err != nil {
			return err
		}

		return c.JSON(res)
	})
}

func RegisterGreeterHTTPServer(s *httpx.Server, srv GreeterServer) {
	s.GET("/sayHello/{name}", _Greeter_SayHello_HTTP_Handler(srv.(types.Service)).(httpx.Handler))
}

var Greeter_HttpServiceDesc = types.HttpServiceDesc{
	ServiceName: "helloword.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []types.HttpMethodDesc{
		{
			Method:   http.MethodGet,
			Template: "/sayHello",
			Handler:  _Greeter_SayHello_HTTP_Handler,
		},
	},
}
