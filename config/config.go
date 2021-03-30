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
	"github.com/blackdreamers/go-micro/v3/logger"
)

var (
	path    = []string{"daydream", "config"}
	Conf    *envConf
	Service = &service{
		EnableDB: true,
	}
)

type envConf struct {
	// env
	Env      string
	LogLevel string

	// etcd
	etcdConf

	// db
	dbConf
}

type etcdConf struct {
	EtcdTLS         bool
	EtcdAuth        bool
	EtcdAddress     []string
	EtcdUser        string
	EtcdPassword    string
	EtcdCaPath      string
	EtcdCertPath    string
	EtcdCertKeyPath string
}

type dbConf struct {
	DBUser       string `json:"user"`
	DBPassword   string `json:"password"`
	DBHost       string `json:"host"`
	DBPort       int    `json:"port"`
	LowThreshold int    `json:"low_threshold"`
}

type service struct {
	SrvName  string
	Name     string
	Type     string
	EnableDB bool
	Version  string
	DBName   string
}

func init() {
	Conf = &envConf{
		Env:      env.GetString(constant.Env, constant.Prod),
		LogLevel: env.GetString(constant.LogLevel, logger.InfoLevel.String()),
		etcdConf: etcdConf{
			EtcdUser:        env.GetString(constant.EtcdUser, ""),
			EtcdPassword:    env.GetString(constant.EtcdPassword, ""),
			EtcdAddress:     env.GetStrings(constant.EtcdAddress),
			EtcdCaPath:      env.GetString(constant.EtcdCaPath, ""),
			EtcdCertPath:    env.GetString(constant.EtcdCertPath, ""),
			EtcdCertKeyPath: env.GetString(constant.EtcdCertKeyPath, ""),
		},
	}
}

func Init() error {
	var err error
	Conf.EtcdAuth, err = env.GetBool(constant.EtcdAuth, false)
	if err != nil {
		return err
	}
	Conf.EtcdTLS, err = env.GetBool(constant.EtcdTLS, false)
	if err != nil {
		return err
	}

	etcdOpts := []source.Option{
		etcd.WithAddress(Conf.EtcdAddress...),
		etcd.WithDialTimeout(5 * time.Second),
		etcd.WithPrefix(fmt.Sprintf("/%s/%s", path[0], path[1])),
	}
	if Conf.EtcdAuth {
		etcdOpts = append(etcdOpts, etcd.Auth(Conf.EtcdUser, Conf.EtcdPassword))
	}

	if err = config.Load(etcd.NewSource(etcdOpts...)); err != nil {
		return err
	}

	var dbConfig dbConf
	if err = Get(constant.DBConfKey).Scan(&dbConfig); err != nil {
		return err
	}
	Conf.dbConf = dbConfig

	Service.SrvName = Service.Type + constant.Delimiter + Service.Name
	Service.Version = Get(Service.Name, "version").String("latest")
	Service.DBName = Service.Get("dbname").String(Service.Name)

	return nil
}

// 获取config路径下配置
func Get(key ...string) reader.Value {
	key = append(path, key...)
	return config.Get(key...)
}

// 获取config/[service]路径下配置
func (s *service) Get(key ...string) reader.Value {
	key = append([]string{s.Name}, key...)
	return Get(key...)
}
