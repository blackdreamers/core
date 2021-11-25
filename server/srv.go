package server

import (
	"go-micro.dev/v4"
	microcli "go-micro.dev/v4/client"
	microsrv "go-micro.dev/v4/server"
)

var (
	srv = &server{}
)

type server struct {
	srv     micro.Service
	handles []interface{}
}

// AddHandles run before service starts
func AddHandles(srvHandles ...interface{}) {
	srv.handles = append(srv.handles, srvHandles...)
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

	return s.srv.Run()
}
