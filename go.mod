module github.com/blackdreamers/core

go 1.16

require (
	github.com/blackdreamers/go-micro/plugins/client/grpc/v3 v3.0.0-20210703052542-87049b6558d4
	github.com/blackdreamers/go-micro/plugins/config/source/etcd/v3 v3.0.0-20210703052542-87049b6558d4
	github.com/blackdreamers/go-micro/plugins/logger/logrus/v3 v3.0.0-20210703052542-87049b6558d4
	github.com/blackdreamers/go-micro/plugins/registry/etcd/v3 v3.0.0-20210703052542-87049b6558d4
	github.com/blackdreamers/go-micro/plugins/server/grpc/v3 v3.0.0-20210703052542-87049b6558d4
	github.com/blackdreamers/go-micro/v3 v3.5.2-0.20210703052542-87049b6558d4
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.7.1
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.5.0
	github.com/go-redis/redis/v8 v8.8.0
	github.com/jinzhu/copier v0.3.2
	github.com/sirupsen/logrus v1.8.1
	github.com/ulule/limiter/v3 v3.8.0
	golang.org/x/text v0.3.6
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.21.7
)
