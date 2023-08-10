package configx

import (
	"testing"
)

func TestFetcherTypes(t *testing.T) {
	RegisterFetcher[Redis](&embeddedRedisFetcher{})
	RegisterFetcher[Etcd](&embeddedEtcdFetcher{})

	config, ok := Fetch[Etcd]()
	t.Log(config, ok)
	config, ok = Fetch[Etcd]()
	t.Log(config, ok)
}
