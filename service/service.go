package service

import (
	"google.golang.org/grpc"
)

type Service interface {
	ServiceDesc() *grpc.ServiceDesc
}
