package config

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/reader"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/config/source/etcd"

	"github.com/blackdreamers/core/constant"
	"github.com/blackdreamers/core/env"
)

var (
	path    = []string{"micro", "config"}
	Env     *envConf
	Service *service
)

type envConf struct {
	// mode
	Mode string
	// service
	ServiceName string
	// etcd
	EtcdTls         bool
	EtcdAuth        bool
	EtcdAddress     []string
	EtcdUser        string
	EtcdPassword    string
	EtcdCaPath      string
	EtcdCertPath    string
	EtcdCertKeyPath string
}

type service struct {
	SrvName string
	Name    string
	Version string
}

func Init() error {
	etcdAuth, err := env.GetBool(constant.EtcdAuth, false)
	if err != nil {
		return err
	}
	etcdTls, err := env.GetBool(constant.EtcdTls, false)
	if err != nil {
		return err
	}

	Env = &envConf{
		EtcdTls:         etcdTls,
		EtcdAuth:        etcdAuth,
		Mode:            env.GetString(constant.Mode, "prod"),
		ServiceName:     env.GetString(constant.ServiceName, ""),
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

	Service = &service{
		Name:    Env.ServiceName,
		SrvName: "srv." + Env.ServiceName,
		Version: Get(Env.ServiceName, "version").String("latest"),
	}

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
