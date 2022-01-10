package config

import (
	"github.com/blackdreamers/core/consts"
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
	return nil
}

func init() {
	auth, err := env.GetBool(consts.EtcdAuth, false)
	if err != nil {
		panic(err)
	}
	TLS, err := env.GetBool(consts.EtcdTLS, false)
	if err != nil {
		panic(err)
	}
	Etcd = &etcdConf{
		TLS:         TLS,
		Auth:        auth,
		User:        env.GetString(consts.EtcdUser, ""),
		Password:    env.GetString(consts.EtcdPassword, ""),
		Addrs:       env.GetStrings(consts.EtcdAddrs),
		CaPath:      env.GetString(consts.EtcdCaPath, ""),
		CertPath:    env.GetString(consts.EtcdCertPath, ""),
		CertKeyPath: env.GetString(consts.EtcdCertKeyPath, ""),
	}
	Configs(Etcd)
}
