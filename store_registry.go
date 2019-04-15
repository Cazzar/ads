package ads

import "github.com/c-mueller/ads/blockstore"

var blockstores map[string]blockstore.BlockListStoreBuilder

func RegisterBlockListStorage(builder blockstore.BlockListStoreBuilder) {
	log.Infof("Registering Blocklist Store %q", builder.Name())
	blockstores[builder.Name()] = builder
}
