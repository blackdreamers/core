module github.com/blackdreamers/core

go 1.15

require (
	github.com/blackdreamers/go-micro/plugins/config/source/etcd/v3 v3.0.0-20210327053124-3d3c2b7a6fa2
	github.com/blackdreamers/go-micro/plugins/logger/logrus/v3 v3.0.0-20210327062615-32425875b001
	github.com/blackdreamers/go-micro/plugins/registry/etcd/v3 v3.0.0-20210327061138-79eae41a7e43
	github.com/blackdreamers/go-micro/plugins/server/grpc/v3 v3.0.0-20210327062615-32425875b001
	github.com/blackdreamers/go-micro/v3 v3.5.1-0.20210328134617-1dfe2330793d
	github.com/sirupsen/logrus v1.8.1
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.21.4
)
