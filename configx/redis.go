package configx

import "time"

type Redis struct {
	// Either a single address or a seed list of host:port addresses
	// of cluster/sentinel nodes.
	Address []string `yaml:"address,omitempty"`

	// Database to be selected after connecting to the server.
	// Only single-node and failover clients.
	DB int `yaml:"db,omitempty"`

	// Common options
	Username         string `yaml:"username,omitempty"`
	Password         string `yaml:"password,omitempty"`
	SentinelUsername string `yaml:"sentinelUsername,omitempty"`
	SentinelPassword string `yaml:"sentinelPassword,omitempty"`

	MaxRetries      int           `yaml:"maxRetries,omitempty"`
	MinRetryBackoff time.Duration `yaml:"minRetryBackoff,omitempty"`
	MaxRetryBackoff time.Duration `yaml:"maxRetryBackoff,omitempty"`

	DialTimeout  time.Duration `yaml:"dialTimeout,omitempty"`
	ReadTimeout  time.Duration `yaml:"readTimeout,omitempty"`
	WriteTimeout time.Duration `yaml:"writeTimeout,omitempty"`

	PoolFIFO    bool          `yaml:"poolFIFO,omitempty"`
	PoolSize    int           `yaml:"poolSize,omitempty"`
	PoolTimeout time.Duration `yaml:"poolTimeout,omitempty"`

	MinIdleConns int `yaml:"minIdleConns,omitempty"`
	MaxIdleConns int `yaml:"MaxIdleConns,omitempty"`

	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime,omitempty"`
	ConnMaxLifetime time.Duration `yaml:"ConnMaxLifetime,omitempty"`

	// Only cluster clients
	MaxRedirects   int  `yaml:"maxRedirects,omitempty"`
	ReadOnly       bool `yaml:"readOnly,omitempty"`
	RouteByLatency bool `yaml:"routeByLatency,omitempty"`
	RouteRandomly  bool `yaml:"routeRandomly,omitempty"`

	// The sentinel master name
	// Only failover clients
	MasterName string `yaml:"masterName,omitempty"`
}

type embeddedRedisFetcher struct{}

func (f *embeddedRedisFetcher) Fetch() (Redis, error) {
	return *embedded.Redis, nil
}
