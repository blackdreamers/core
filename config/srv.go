package config

import (
	"github.com/blackdreamers/core/constant"
)

var (
	Service = &service{
		EnableDB:     true,
		EnableBroker: true,
	}
)

type service struct {
	SrvName      string
	Name         string
	Type         string
	Version      string
	EnableDB     bool
	EnableBroker bool
	DBName       string
	Addr         string
	AllowOrigins []string
}

func (s *service) init() error {
	s.SrvName = s.Type + constant.Delimiter + s.Name
	s.Version = s.Get("version").String("latest")
	s.DBName = s.Get("dbname").String(s.Name)
	s.Addr = s.Get("addr").String(":8080")
	s.AllowOrigins = s.Get("allow_origins").StringSlice([]string{"*"})

	return nil
}

func init() {
	Configs(Service)
}
