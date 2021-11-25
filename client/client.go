package client

import (
	microcli "go-micro.dev/v4/client"
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
