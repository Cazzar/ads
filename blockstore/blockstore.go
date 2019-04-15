package blockstore

import "github.com/mholt/caddy"

type BlockListStoreBuilder interface {
	Build(c *caddy.Controller) (BlockListStore, error)
	Name() string
}

type BlockListStore interface {
	Init() error
	Name() string

	IsBlacklisted(qname string) bool
	IsWhitelisted(qname string) bool

	AddBlacklistRule(qname string) error
	AddWhitelistRule(qname string) error
	AddRegexBlacklistRule(qname string) error
	AddRegexWhitelistRule(qname string) error
}
