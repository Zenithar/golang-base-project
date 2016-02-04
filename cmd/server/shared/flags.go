package shared

// Flags represents system flags
type Flags struct {
	BindAddress      string
	RavenDSN         string
	LogFormatterType string
	ForceColors      bool

	DatabaseDriver    string
	DatabaseHost      string
	DatabaseNamespace string
	DatabaseUser      string
	DatabasePassword  string

	MemcachedHosts string
	RedisHost      string
}
