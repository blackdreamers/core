package server

import (
	"context"
	"time"

	"github.com/blackdreamers/core/broker"
	"github.com/blackdreamers/core/cache/redis"
	"github.com/blackdreamers/core/client"
	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/constant"
	"github.com/blackdreamers/core/cron"
	"github.com/blackdreamers/core/db"
	"github.com/blackdreamers/core/logger"
	"github.com/blackdreamers/core/utils"
	"github.com/blackdreamers/go-micro/plugins/broker/nsq/v3"
	cgrpc "github.com/blackdreamers/go-micro/plugins/client/grpc/v3"
	"github.com/blackdreamers/go-micro/plugins/registry/etcd/v3"
	sgrpc "github.com/blackdreamers/go-micro/plugins/server/grpc/v3"
	"github.com/blackdreamers/go-micro/v3"
	"github.com/blackdreamers/go-micro/v3/registry"
)

type Server interface {
	init(opts ...micro.Option) error
	run() error
}

func Init(opts ...micro.Option) {
	for _, o := range opts {
		o(&micro.Options{})
	}

	if err := config.Init(); err != nil {
		panic(err)
	}

	if err := logger.Init(); err != nil {
		panic(err)
	}

	if config.Service.EnableDB {
		if err := db.Init(); err != nil {
			panic(err)
		}
	}

	if err := redis.Init(); err != nil {
		panic(err)
	}

	regOpts := []registry.Option{
		registry.Addrs(config.Etcd.Addrs...),
	}

	opts = append(
		opts,
		micro.Server(sgrpc.NewServer()),
		micro.Client(cgrpc.NewClient()),
		micro.Name(config.Service.SrvName),
		micro.Version(config.Service.Version),
		micro.Broker(nsq.NewBroker()),
		micro.AfterStart(func() error {
			client.Init(Client())
			return nil
		}),
		micro.AfterStart(func() error {
			return cron.Init()
		}),
		micro.BeforeStop(func() error {
			return broker.Broker().Disconnect()
		}),
		micro.BeforeStop(func() error {
			cron.Stop()
			return nil
		}),
	)

	if config.Service.EnableBroker {
		opts = append(opts, micro.AfterStart(func() error {
			return broker.Init(srv.srv.Options().Broker)
		}))
	}

	if config.Registry == "" || config.Registry == constant.Etcd {
		if config.Etcd.Auth {
			regOpts = append(regOpts, etcd.Auth(config.Etcd.User, config.Etcd.Password))
		}

		if config.Etcd.TLS {
			tLSConf, err := utils.GetTLSConfig(config.Etcd.CaPath, config.Etcd.CertPath, config.Etcd.CertKeyPath)
			if err != nil {
				panic(err)
			}
			regOpts = append(regOpts, registry.TLSConfig(tLSConf))
		}

		opts = append(opts, micro.Registry(etcd.NewRegistry(regOpts...)))
	}

	if config.Service.Type == API {
		opts = append(
			opts,
			micro.AfterStart(func() error {
				return api.run()
			}),
			micro.BeforeStop(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				return api.s.Shutdown(ctx)
			}),
		)
	}

	if err := srv.init(opts...); err != nil {
		panic(err)
	}

	if config.Service.Type == API {
		if err := api.init(opts...); err != nil {
			panic(err)
		}
	}

}

func Run() {
	if err := srv.run(); err != nil {
		panic(err)
	}
}
