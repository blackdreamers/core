package server

import (
	"context"
	"time"

	"github.com/go-micro/plugins/v4/broker/nsq"
	cgrpc "github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/etcd"
	sgrpc "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	"go-micro.dev/v4"
	microcli "go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/registry"

	"github.com/blackdreamers/core/api/auth"
	"github.com/blackdreamers/core/api/websocket"
	"github.com/blackdreamers/core/broker"
	"github.com/blackdreamers/core/cache/redis"
	corecli "github.com/blackdreamers/core/client"
	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/consts"
	"github.com/blackdreamers/core/cron"
	"github.com/blackdreamers/core/cron/jobs"
	"github.com/blackdreamers/core/db"
	"github.com/blackdreamers/core/logger"
	"github.com/blackdreamers/core/utils"
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
		micro.WrapCall(func(cf microcli.CallFunc) microcli.CallFunc {
			return func(ctx context.Context, node *registry.Node, req microcli.Request, rsp interface{},
				opts microcli.CallOptions) error {
				err := cf(metadata.Set(ctx, consts.SrvNameKey, config.Service.SrvName), node, req, rsp, opts)
				return err
			}
		}),
		micro.AfterStart(func() error {
			corecli.Init(Client())
			return nil
		}),
		micro.AfterStart(func() error {
			return cron.Init()
		}),
		micro.BeforeStop(func() error {
			cron.Stop()
			return nil
		}),
	)

	if !config.IsDevEnv() {
		opts = append(opts, micro.WrapHandler(
			prometheus.NewHandlerWrapper(
				prometheus.ServiceName(config.Service.Name),
			),
		))
	}

	if config.Service.EnableBroker {
		opts = append(opts, micro.Broker(nsq.NewBroker()))
		opts = append(opts, micro.AfterStart(func() error {
			return broker.Init(srv.srv.Options().Broker)
		}))
		opts = append(opts, micro.BeforeStop(func() error {
			return broker.Broker().Disconnect()
		}))
	}

	if config.Registry == "" || config.Registry == consts.Etcd {
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
				_ = websocket.WS.Close()
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
		cron.AddJobs(&jobs.CasbinPolicy{})
		if err := auth.Init(); err != nil {
			panic(err)
		}
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
