package server

import (
	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/go-micro/v3"
	microcli "github.com/blackdreamers/go-micro/v3/client"
	microsrv "github.com/blackdreamers/go-micro/v3/server"
)

var (
	srv = &server{}
)

type server struct {
	srv         micro.Service
	handles     []interface{}
	subscribers []interface{}
}

func Handles(srvHandles ...interface{}) {
	srv.handles = append(srv.handles, srvHandles...)
}

func Subscribers(srvSubscribers ...interface{}) {
	srv.subscribers = append(srv.subscribers, srvSubscribers...)
}

// Service server.Service().Init(microsrv.Wait(nil))
func Service() microsrv.Server {
	return srv.srv.Server()
}

func Client() microcli.Client {
	return srv.srv.Client()
}

func (s *server) init(opts ...micro.Option) error {
	s.srv = micro.NewService(opts...)
	s.srv.Init()

	return nil
}

func (s *server) run() error {
	// handles
	for _, handle := range s.handles {
		if err := micro.RegisterHandler(s.srv.Server(), handle); err != nil {
			panic(err)
		}
	}

	// subscribers
	for _, subscribe := range s.subscribers {
		if err := micro.RegisterSubscriber(config.Service.SrvName, s.srv.Server(), subscribe); err != nil {
			panic(err)
		}
	}

	return s.srv.Run()
}
