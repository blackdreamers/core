package config

import (
	"github.com/blackdreamers/core/constant"
	"github.com/blackdreamers/core/env"

	"github.com/asim/go-micro/plugins/config/source/etcd/v3"
	"github.com/asim/go-micro/v3/config"
	"github.com/asim/go-micro/v3/config/reader"
	"github.com/asim/go-micro/v3/config/source"
)

var (
	path    = []string{"micro", "config"}
	Env     *envConf
	Service = &service{
		EnableDB: true,
	}
)

const (
	delimiter = "."
)

type envConf struct {
	// mode
	Mode string

	// etcd
	EtcdTLS         bool
	EtcdAuth        bool
	EtcdAddress     []string
	EtcdUser        string
	EtcdPassword    string
	EtcdCaPath      string
	EtcdCertPath    string
	EtcdCertKeyPath string
}

type service struct {
	SrvName  string
	Name     string
	Type     string
	EnableDB bool
	Version  string
}

func Init() error {
	etcdAuth, err := env.GetBool(constant.EtcdAuth, false)
	if err != nil {
		return err
	}
	etcdTLS, err := env.GetBool(constant.EtcdTLS, false)
	if err != nil {
		return err
	}

	Env = &envConf{
		EtcdTLS:         etcdTLS,
		EtcdAuth:        etcdAuth,
		Mode:            env.GetString(constant.Mode, "prod"),
		EtcdUser:        env.GetString(constant.EtcdUser, ""),
		EtcdPassword:    env.GetString(constant.EtcdPassword, ""),
		EtcdAddress:     env.GetStrings(constant.EtcdAddress),
		EtcdCaPath:      env.GetString(constant.EtcdCaPath, ""),
		EtcdCertPath:    env.GetString(constant.EtcdCertPath, ""),
		EtcdCertKeyPath: env.GetString(constant.EtcdCertKeyPath, ""),
	}

	etcdOpts := []source.Option{
		etcd.WithAddress(Env.EtcdAddress...),
	}
	if Env.EtcdAuth {
		etcdOpts = append(etcdOpts, etcd.Auth(Env.EtcdUser, Env.EtcdPassword))
	}

	if err := config.Load(etcd.NewSource(etcdOpts...)); err != nil {
		return err
	}

	Service.SrvName = Service.Type + delimiter + Service.Name
	Service.Version = Get(Service.Name, "version").String("latest")

	return nil
}

// 获取/micro/config路径下配置
func Get(key ...string) reader.Value {
	key = append(path, key...)
	return config.Get(key...)
}

// 获取/micro/config/[service]路径下配置
func (s *service) Get(key ...string) reader.Value {
	key = append([]string{s.Name}, key...)
	return Get(key...)
}
