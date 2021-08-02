package client

import (
	microcli "github.com/blackdreamers/go-micro/v3/client"
)

var (
	clients []Client
)

type Client interface {
	Name() string
	Init(client microcli.Client)
}

func AddClients(clis ...Client) {
	clients = append(clients, clis...)
}

func Init(client microcli.Client) {
	for _, c := range clients {
		c.Init(client)
	}
}
