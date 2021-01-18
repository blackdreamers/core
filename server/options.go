package server

import (
	"github.com/blackdreamers/core/config"

	microsrv "github.com/micro/go-micro/v2/server"
)

const (
	SRV = "srv"
	API = "api"
)

// Server name
func Name(n string) microsrv.Option {
	return func(o *microsrv.Options) {
		config.Service.Name = n
	}
}

// Server type
func Type(t string) microsrv.Option {
	return func(o *microsrv.Options) {
		config.Service.Type = t
	}
}

// Whether the server use db, default true
func EnableDB(b bool) microsrv.Option {
	return func(o *microsrv.Options) {
		config.Service.EnableDB = b
	}
}
