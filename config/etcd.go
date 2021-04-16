package config

import (
	"github.com/blackdreamers/core/constant"
	"github.com/blackdreamers/core/env"
)

var (
	Etcd *etcdConf
)

type etcdConf struct {
	TLS         bool
	Auth        bool
	Addrs       []string
	User        string
	Password    string
	CaPath      string
	CertPath    string
	CertKeyPath string
}

func (e *etcdConf) init() error {
	var err error
	e.Auth, err = env.GetBool(constant.EtcdAuth, false)
	if err != nil {
		return err
	}
	e.TLS, err = env.GetBool(constant.EtcdTLS, false)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	Etcd = &etcdConf{
		User:        env.GetString(constant.EtcdUser, ""),
		Password:    env.GetString(constant.EtcdPassword, ""),
		Addrs:       env.GetStrings(constant.EtcdAddrs),
		CaPath:      env.GetString(constant.EtcdCaPath, ""),
		CertPath:    env.GetString(constant.EtcdCertPath, ""),
		CertKeyPath: env.GetString(constant.EtcdCertKeyPath, ""),
	}
	Configs(Etcd)
}
