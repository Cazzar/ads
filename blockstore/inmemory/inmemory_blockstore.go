package inmemory

import (
	"github.com/c-mueller/ads"
	"github.com/c-mueller/ads/blockstore"
	"github.com/mholt/caddy"
)

func init() {
	ads.RegisterBlockListStorage(inmemoryBlockstoreBuilder{})
}

type inmemoryBlockstoreBuilder struct {
}

func (b inmemoryBlockstoreBuilder) Build(c *caddy.Controller) (blockstore.BlockListStore, error) {
	panic("implement me")
}

func (b inmemoryBlockstoreBuilder) Name() string {
	return "inmemory"
}
