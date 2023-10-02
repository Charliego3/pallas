package configx

import (
	"testing"
)

func TestFetcherTypes(t *testing.T) {
	RegisterFetcher[Redis](&standardRedisFetcher{})
	RegisterFetcher[Etcd](&standardEtcdFetcher{})

	config, ok := Fetch[Etcd]()
	t.Log(config, ok)
	config, ok = Fetch[Etcd]()
	t.Log(config, ok)
}
