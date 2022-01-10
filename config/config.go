package config

import (
	"fmt"
	"time"

	"github.com/asim/go-micro/plugins/config/source/etcd/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/source"

	"github.com/blackdreamers/core/consts"
	"github.com/blackdreamers/core/env"
)

var (
	path     = []string{"daydream", "config"}
	Env      string
	Registry string
	configs  []conf
)

type conf interface {
	init() error
}

func init() {
	Env = env.GetString(consts.Env, consts.Prod)
	Registry = env.GetString(consts.Registry, "")
}

func Init() error {
	etcdOpts := []source.Option{
		etcd.WithAddress(Etcd.Addrs...),
		etcd.WithDialTimeout(5 * time.Second),
		etcd.WithPrefix(fmt.Sprintf("/%s/%s", path[0], path[1])),
	}
	if Etcd.Auth {
		etcdOpts = append(etcdOpts, etcd.Auth(Etcd.User, Etcd.Password))
	}

	if err := config.Load(etcd.NewSource(etcdOpts...)); err != nil {
		return err
	}

	for _, c := range configs {
		if err := c.init(); err != nil {
			return err
		}
	}

	return nil
}

func IsDevEnv() bool {
	return Env == consts.Dev
}

func IsTestEnv() bool {
	return Env == consts.Test
}

func IsProdEnv() bool {
	return Env == consts.Prod || Env == consts.Release
}

func Configs(cs ...conf) {
	configs = append(configs, cs...)
}

// Get path: daydream/config
func Get(key ...string) reader.Value {
	key = append(path, key...)
	return config.Get(key...)
}

// Get path: daydream/config/[service]
func (s *service) Get(key ...string) reader.Value {
	key = append([]string{s.Name}, key...)
	return Get(key...)
}
