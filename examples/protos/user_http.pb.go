// Code generated by protoc-gen-pallas-http. DO NOT EDIT.
//
// proto-gen-pallas-http version: 1.0.0
// protoc version: v4.25.0
// source file: protos/user.proto

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

type UserHTTPServer interface {
	Register(ctx context.Context, in *RegisterRequest) (*LoginReply, error)
	Login(ctx context.Context, in *LoginRequest) (*LoginReply, error)
}

func RegisterUserHTTPServer(s *httpx.Server, srv UserHTTPServer) {
	s.HandleMethod("POST", "/user/register", _User_Register_POST_HTTP_Handler(srv.(types.Service)).(httpx.Handler))
	s.HandleMethod("POST", "/user/login", _User_Login_POST_HTTP_Handler(srv.(types.Service)).(httpx.Handler))
}

func _User_Register_POST_HTTP_Handler(srv types.Service) any {
	return httpx.HandlerFunc(func(c *httpx.Context) error {
		req := new(RegisterRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		res, err := srv.(UserServer).Register(c.Context, req)
		if err != nil {
			return err
		}
		return c.Write(res)
	})
}

func _User_Login_POST_HTTP_Handler(srv types.Service) any {
	return httpx.HandlerFunc(func(c *httpx.Context) error {
		req := new(LoginRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		res, err := srv.(UserServer).Login(c.Context, req)
		if err != nil {
			return err
		}
		return c.Write(res)
	})
}
