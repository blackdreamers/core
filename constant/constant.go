package constant

const (
	Dev     = "dev"
	Test    = "test"
	Prod    = "prod"
	Release = "release"
)

const (
	Timestamp = "2006-01-02 15:04:05"
)

const (
	Delimiter   = "."
	SourceField = "source"
	ErrKey      = "error"
)

const (
	DBConfKey      = "database"
	RedisConfKey   = "redis"
	LimiterConfKey = "limiter"
	SessionConfKey = "session"
)

const (
	Env      = "ENV"
	LogLevel = "LOG_LEVEL"
	Registry = "MICRO_REGISTRY"
)

const (
	Etcd = "etcd"
)

const (
	EtcdUser        = "ETCD_USER"
	EtcdPassword    = "ETCD_PASSWORD"
	EtcdAddrs       = "ETCD_ADDRS"
	EtcdAuth        = "ETCD_AUTH"
	EtcdTLS         = "ETCD_TLS"
	EtcdCaPath      = "ETCD_CA_PATH"
	EtcdCertPath    = "ETCD_CERT_PATH"
	EtcdCertKeyPath = "ETCD_CERT_KEY_PATH"
)

const (
	MemoryStore = "memory"
	RedisStore  = "redis"
	CookieStore = "cookie"
)
