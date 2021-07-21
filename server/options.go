package server

import (
	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/go-micro/v3"
)

const (
	SRV = "srv"
	API = "api"
)

func Name(n string) micro.Option {
	return func(o *micro.Options) {
		config.Service.Name = n
	}
}

func Type(t string) micro.Option {
	return func(o *micro.Options) {
		config.Service.Type = t
	}
}

// EnableDB Whether the server use db, default true
func EnableDB(b bool) micro.Option {
	return func(o *micro.Options) {
		config.Service.EnableDB = b
	}
}

// EnableBroker Whether the server use broker, default true
func EnableBroker(b bool) micro.Option {
	return func(o *micro.Options) {
		config.Service.EnableBroker = b
	}
}
