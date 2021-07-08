package config

import (
	"fmt"
	"time"

	"github.com/blackdreamers/core/constant"
	"github.com/blackdreamers/core/env"
	"github.com/blackdreamers/go-micro/plugins/config/source/etcd/v3"
	"github.com/blackdreamers/go-micro/v3/config"
	"github.com/blackdreamers/go-micro/v3/config/reader"
	"github.com/blackdreamers/go-micro/v3/config/source"
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
	Env = env.GetString(constant.Env, constant.Prod)
	Registry = env.GetString(constant.Registry, "")
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
	return Env == constant.Dev
}

func IsTestEnv() bool {
	return Env == constant.Test
}

func IsProdEnv() bool {
	return Env == constant.Prod || Env == constant.Release
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
