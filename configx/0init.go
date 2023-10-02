package configx

import (
	"encoding/json"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"

	"github.com/charliego3/argsx"
	"github.com/gookit/goutil/fsutil"
)

type StandardConfig struct {
	App      *App      `json:"app,omitempty" yaml:"app,omitempty" toml:"app,omitempty"`
	Etcd     *Etcd     `json:"etcd,omitempty" yaml:"etcd,omitempty" toml:"etcd,omitempty"`
	Database *Database `json:"database,omitempty" yaml:"database,omitempty" toml:"database,omitempty"`
	Redis    *Redis    `json:"redis,omitempty" yaml:"redis,omitempty" toml:"redis,omitempty"`
	Logger   *Logger   `json:"logger,omitempty" yaml:"logger,omitempty" toml:"logger,omitempty"`
}

var standard = StandardConfig{
	Logger: &Logger{
		Level:  slog.LevelInfo.String(),
		Format: HandlerTypeText,
	},
}

func init() {
	configPath := argsx.Fetch("config").MustString("./config.yaml")
	if !fsutil.FileExist(configPath) {
		return
	}

	bs, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("read config file", slog.String("path", configPath), slog.Any("err", err))
		os.Exit(1)
	}

	switch filepath.Ext(configPath) {
	case ".yaml":
		err = yaml.Unmarshal(bs, &standard)
	case ".json":
		err = json.Unmarshal(bs, &standard)
	}

	if err != nil {
		slog.Error("failed to load default config from file", slog.String("path", configPath), slog.Any("err", err))
		os.Exit(1)
	}

	register[Etcd](standard.Etcd, &standardEtcdFetcher{})
	register[App](standard.App, &standardAppFetcher{})
	register[Redis](standard.Redis, &standardRedisFetcher{})
	register[Database](standard.Database, &standardDatabaseFetcher{})
	register[Logger](standard.Logger, &standardLoggerConfig{})
}

// register register fetcher to fetchers if obj is not nil
func register[T any](obj *T, fetcher Fetcher[T]) {
	if obj == nil {
		return
	}

	RegisterFetcher(fetcher)
}
