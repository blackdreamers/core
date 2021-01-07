package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	microsrv "github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/server/grpc"

	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/db"
	"github.com/blackdreamers/core/utils"
)

var (
	handles     []interface{}
	subscribers []interface{}
)

func Handles(srvHandles ...interface{}) {
	handles = append(handles, srvHandles...)
}

func Subscribers(srvSubscribers ...interface{}) {
	subscribers = append(subscribers, srvSubscribers...)
}

func Run(opts ...microsrv.Option) {
	if err := config.Init(); err != nil {
		panic(err)
	}

	if err := db.Init(); err != nil {
		panic(err)
	}

	regOpts := []registry.Option{
		registry.Addrs(config.Env.EtcdAddress...),
	}

	if config.Env.EtcdAuth {
		regOpts = append(regOpts, etcd.Auth(config.Env.EtcdUser, config.Env.EtcdPassword))
	}

	if config.Env.EtcdTls {
		tLSConf, err := utils.GetTLSConfig()
		if err != nil {
			panic(err)
		}
		regOpts = append(regOpts, registry.TLSConfig(tLSConf))
	}

	opts = append(opts, []microsrv.Option{
		microsrv.Name(config.Service.SrvName),
		microsrv.Version(config.Service.Version),
		microsrv.Registry(etcd.NewRegistry(regOpts...)),
	}...)

	srv := grpc.NewServer(opts...)

	// handles
	for _, handle := range handles {
		if err := srv.Handle(srv.NewHandler(handle)); err != nil {
			panic(err)
		}
	}

	// subscribers
	for _, subscribe := range subscribers {
		if err := srv.Subscribe(
			srv.NewSubscriber(
				config.Service.SrvName,
				subscribe,
			),
		); err != nil {
			panic(err)
		}
	}

	if err := srv.Start(); err != nil {
		panic(err)
	}

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait
	if err := srv.Stop(); err != nil {
		panic(err)
	}

}
