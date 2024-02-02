package biz

import (
	"context"
)

type UserRepo interface {
	FetchByUsername(context.Context)
}
