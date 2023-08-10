package configx

import "time"

type Database struct {
	DSN     string `json:"dsn" yaml:"dsn"`
	Options struct {
		MaxOpenConns int           `json:"maxOpenConns,omitempty" yaml:"maxOpenConns,omitempty"`
		MaxIdleConns int           `json:"maxIdleConns,omitempty" yaml:"maxIdleConns,omitempty"`
		MaxIdleTime  time.Duration `json:"maxIdleTime,omitempty" yaml:"maxIdleTime,omitempty"`
		MaxLifetime  time.Duration `json:"maxLifetime,omitempty" yaml:"maxLifetime,omitempty"`
	} `json:"options,omitempty" yaml:"options,omitempty"`
}

type embeddedDatabaseFetcher struct{}

func (f *embeddedDatabaseFetcher) Fetch() (Database, error) {
	return *embedded.Database, nil
}
