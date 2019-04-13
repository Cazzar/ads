package blockstore

type BlockListStore interface {
	Init() error

	IsBlacklisted(qname string) bool
	IsWhitelisted(qname string) bool

	AddBlacklistRule(qname string) error
	AddWhitelistRule(qname string) error
	AddRegexBlacklistRule(qname string) error
	AddRegexWhitelistRule(qname string) error
}
