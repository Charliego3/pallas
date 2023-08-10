package etcdx

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/charliego3/logger"
	"github.com/charliego3/mspp/configx"
	"github.com/charliego3/mspp/utils"
	"github.com/gookit/goutil/strutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	ev3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	*ev3.Client
}

// New returns etcd client from options
func New(opts ...Option) (*Client, error) {
	var config zap.Config
	if logger.GetLevel() == logger.LevelDebug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.ConsoleSeparator = " "
	} else {
		config = zap.NewProductionConfig()
	}

	config.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(time.DateTime))
	}
	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, pae zapcore.PrimitiveArrayEncoder) {
		sl := l.CapitalString()
		if len(sl) > 4 {
			sl = sl[:4]
		}
		pae.AppendString(sl)
	}
	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	cfg := &ev3.Config{Logger: log}
	for _, opt := range opts {
		opt(cfg)
	}

	if len(cfg.Endpoints) == 0 {
		cfg.Endpoints = []string{":2379"}
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = time.Minute
	}

	cli, err := ev3.New(*cfg)
	return &Client{Client: cli}, err
}

var (
	c *Client
	s sync.Once
)

// C returns etcd client from config
func C() *Client {
	s.Do(func() {
		cfg, err := configx.Fetch[configx.Etcd]()
		if err != nil {
			logger.Fatal("not resolved etcd config", "err", err)
		}

		var tlsConfig *tls.Config
		if strutil.IsNotBlank(cfg.RootCA) {
			if utils.AnyBlank(cfg.PemKey, cfg.PemCert) {
				logger.Fatalf("etcd certificate path is incorrect, Key: %q, Cert: %q", cfg.PemKey, cfg.PemCert)
			}

			etcdCA, err := os.ReadFile(cfg.RootCA)
			if err != nil {
				logger.Fatal("read etcd RootCA fail", "err", err)
			}

			pair, err := tls.LoadX509KeyPair(cfg.PemCert, cfg.PemKey)
			if err != nil {
				logger.Fatal("load etcd CA cert and key fail", "err", err)
			}

			rootCertPool := x509.NewCertPool()
			rootCertPool.AppendCertsFromPEM(etcdCA)

			tlsConfig = &tls.Config{
				RootCAs:      rootCertPool,
				Certificates: []tls.Certificate{pair},
			}
		}

		opts := []Option{
			WithEndpoints(cfg.Endpoints...),
			WithAutoSyncInterval(cfg.AutoSyncInterval),
			WithDialTimeout(cfg.DialTimeout),
			WithDialKeepAliveTime(cfg.DialKeepAliveTime),
			WithDialKeepAliveTimeout(cfg.DialKeepAliveTimeout),
			WithMaxCallSendSize(cfg.MaxCallSendSize),
			WithMaxCallRecvSize(cfg.MaxCallRecvSize),
			WithUsername(cfg.Username),
			WithPassword(cfg.Password),
			WithPermWithoutStream(cfg.PermWithoutStream),
			WithRejectOldCluster(cfg.RejectOldCluster),
		}

		if tlsConfig != nil {
			opts = append(opts, WithTLS(tlsConfig))
		}

		client, err := New(opts...)
		if err != nil {
			logger.Fatal("connect etcd server fail", "err", err)
		}

		c = client
		logger.Info("connected etcd", "endpoints", c.Endpoints())
	})
	return c
}

// Fetch returns response, default timeout is 3s
func Fetch(key string, opts ...OpOpt) (*ev3.GetResponse, error) {
	f := getFOpts(opts)
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()
	resp, err := C().Get(ctx, key, f.opOpts...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed fetch from etcd: %s", key)
	}

	if resp.Count == 0 {
		return nil, fmt.Errorf("there is no value in etcd for key: %s", key)
	}
	return resp, nil
}
