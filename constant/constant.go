package constant

const (
	Dev     = "dev"
	Test    = "test"
	Prod    = "prod"
	Release = "release"
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
	BrokerConfKey  = "broker"
)

const (
	Env         = "ENV"
	LogLevel    = "LOG_LEVEL"
	Registry    = "MICRO_REGISTRY"
	BrokerAddrs = "BROKER_ADDRS"
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
