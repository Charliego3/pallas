package etcdx

import (
	"context"
	"crypto/tls"
	"time"

	ev3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type Option func(*ev3.Config)

func WithEndpoints(endpoints ...string) Option {
	return func(c *ev3.Config) {
		c.Endpoints = endpoints
	}
}

func WithAutoSyncInterval(interval time.Duration) Option {
	return func(c *ev3.Config) {
		c.AutoSyncInterval = interval
	}
}

func WithDialTimeout(timeout time.Duration) Option {
	return func(c *ev3.Config) {
		c.DialTimeout = timeout
	}
}

func WithDialKeepAliveTime(keepAliveTime time.Duration) Option {
	return func(c *ev3.Config) {
		c.DialKeepAliveTime = keepAliveTime
	}
}

func WithDialKeepAliveTimeout(timeout time.Duration) Option {
	return func(c *ev3.Config) {
		c.DialKeepAliveTimeout = timeout
	}
}

func WithMaxCallSendSize(maxSize int) Option {
	return func(c *ev3.Config) {
		c.MaxCallSendMsgSize = maxSize
	}
}

func WithMaxCallRecvSize(maxSize int) Option {
	return func(c *ev3.Config) {
		c.MaxCallRecvMsgSize = maxSize
	}
}

func WithRejectOldCluster(reject bool) Option {
	return func(c *ev3.Config) {
		c.RejectOldCluster = reject
	}
}

func WithTLS(tls *tls.Config) Option {
	return func(c *ev3.Config) {
		c.TLS = tls
	}
}

func WithUsername(username string) Option {
	return func(c *ev3.Config) {
		c.Username = username
	}
}

func WithPassword(password string) Option {
	return func(c *ev3.Config) {
		c.Password = password
	}
}

func WithDialOptions(opts []grpc.DialOption) Option {
	return func(c *ev3.Config) {
		c.DialOptions = opts
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *ev3.Config) {
		c.Context = ctx
	}
}

func WithPermWithoutStream(without bool) Option {
	return func(c *ev3.Config) {
		c.PermitWithoutStream = without
	}
}
