package configx

import (
	"time"
)

type Etcd struct {
	Endpoints            []string      `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	AutoSyncInterval     time.Duration `json:"autoSyncInterval,omitempty" yaml:"autoSyncInterval,omitempty"`
	DialTimeout          time.Duration `json:"dialTimeout,omitempty" yaml:"dialTimeout,omitempty"`
	DialKeepAliveTime    time.Duration `json:"dialKeepAliveTime,omitempty" yaml:"dialKeepAliveTime,omitempty"`
	DialKeepAliveTimeout time.Duration `json:"dialKeepAliveTimeout,omitempty" yaml:"dialKeepAliveTimeout,omitempty"`
	MaxCallSendSize      int           `json:"maxCallSendSize,omitempty" yaml:"maxCallSendSize,omitempty"`
	MaxCallRecvSize      int           `json:"maxCallRecvSize,omitempty" yaml:"maxCallRecvSize,omitempty"`
	Username             string        `json:"username,omitempty" yaml:"username,omitempty"`
	Password             string        `json:"password,omitempty" yaml:"password,omitempty"`
	PermWithoutStream    bool          `json:"permWithoutStream,omitempty" yaml:"permWithoutStream,omitempty"`
	RejectOldCluster     bool          `json:"rejectOldCluster,omitempty" yaml:"rejectOldCluster,omitempty"`

	// TLS
	RootCA  string `json:"rootCA,omitempty" yaml:"rootCA,omitempty"`
	PemKey  string `json:"pemKey,omitempty" yaml:"pemKey,omitempty"`
	PemCert string `json:"pemCert,omitempty" yaml:"pemCert,omitempty"`
}

type embeddedEtcdFetcher struct{}

func (f *embeddedEtcdFetcher) Fetch() (Etcd, error) {
	return *embedded.Etcd, nil
}
