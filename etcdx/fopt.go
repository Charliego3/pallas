package etcdx

import (
	"time"

	ev3 "go.etcd.io/etcd/client/v3"
)

type fopt struct {
	timeout time.Duration
	opOpts  []ev3.OpOption
}

type OpOpt func(*fopt)

func getFOpts(opts []OpOpt) *fopt {
	f := &fopt{timeout: time.Second * 3}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithTimeout(timeout time.Duration) OpOpt {
	return func(f *fopt) {
		f.timeout = timeout
	}
}

func WithOpOpts(opts ...ev3.OpOption) OpOpt {
	return func(f *fopt) {
		f.opOpts = opts
	}
}
