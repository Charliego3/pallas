package configx

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"

	"github.com/charliego3/argsx"
	"github.com/charliego3/logger"
	"github.com/gookit/goutil/fsutil"
)

type EmbeddedConfig struct {
	App      *App      `json:"app,omitempty" yaml:"app,omitempty"`
	Etcd     *Etcd     `json:"etcd,omitempty" yaml:"etcd,omitempty"`
	Database *Database `json:"database,omitempty" yaml:"database,omitempty"`
	Redis    *Redis    `json:"redis,omitempty" yaml:"redis,omitempty"`
}

var embedded EmbeddedConfig

func init() {
	configPath := argsx.Fetch("config").MustString("./config.yaml")
	if !fsutil.FileExist(configPath) {
		return
	}

	bs, err := os.ReadFile(configPath)
	if err != nil {
		logger.Fatal("read config file", "path", configPath, "err", err)
	}

	switch filepath.Ext(configPath) {
	case ".yaml":
		err = yaml.Unmarshal(bs, &embedded)
	case ".json":
		err = json.Unmarshal(bs, &embedded)
	}

	if err != nil {
		logger.Fatal("failed to load default config from file", "file", configPath, "err", err)
	}

	register[Etcd](embedded.Etcd, &embeddedEtcdFetcher{})
	register[App](embedded.App, &embeddedAppFetcher{})
	register[Redis](embedded.Redis, &embeddedRedisFetcher{})
	register[Database](embedded.Database, &embeddedDatabaseFetcher{})
}

// register register fetcher to fetchers if obj is not nil
func register[T any](obj *T, fetcher Fetcher[T]) {
	if obj == nil {
		return
	}

	RegisterFetcher(fetcher)
}
