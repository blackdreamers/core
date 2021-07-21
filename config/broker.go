package config

import (
	"github.com/blackdreamers/core/constant"
	"github.com/blackdreamers/core/env"
)

var (
	Broker = &brokerConf{}
)

type brokerConf struct {
	Addrs []string `json:"-"`
}

func (b *brokerConf) init() error {
	// env中配置的权重高于etcd中配置的，便于使用测试机nsq
	if len(b.Addrs) == 0 {
		b.Addrs = Get(constant.BrokerConfKey, "addrs").StringSlice([]string{})
	}
	return nil
}

func init() {
	Broker = &brokerConf{
		Addrs: env.GetStrings(constant.BrokerAddrs),
	}
	Configs(Broker)
}
