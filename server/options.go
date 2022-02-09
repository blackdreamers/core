package server

import (
	"go-micro.dev/v4"

	"github.com/blackdreamers/core/config"
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

// DisableDB Server not using db
func DisableDB() micro.Option {
	return func(o *micro.Options) {
		config.Service.EnableDB = false
	}
}

// DisableBroker Server not using broker
func DisableBroker() micro.Option {
	return func(o *micro.Options) {
		config.Service.EnableBroker = false
	}
}

// Private Server not using public components, must enable db
func Private() micro.Option {
	return func(o *micro.Options) {
		config.Service.Private = true
	}
}
