package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/charliego3/pallas/examples/protos"
)

type User struct {
	pb.UnimplementedUserDescServer
}

func (u *User) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.LoginReply, error) {
	fmt.Printf("%+v\n", in)
	return &pb.LoginReply{
		Message: "register api",
	}, nil
}

func (u *User) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	return nil, errors.New("balabala")
}
