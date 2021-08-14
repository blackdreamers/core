package client

import (
	"github.com/blackdreamers/core/client"
	microcli "github.com/blackdreamers/go-micro/v3/client"
	"github.com/blackdreamers/platform/proto/auth"
)

var (
	Platform = &platform{}
)

type platform struct {
	Auth auth.AuthService
}

func (p *platform) Name() string {
	return "srv.platform"
}

func (p *platform) Init(client microcli.Client) {
	p.Auth = auth.NewAuthService(p.Name(), client)
}

func init() {
	client.AddClients(Platform)
}
