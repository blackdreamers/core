package config

import (
	"github.com/blackdreamers/core/constant"
)

var (
	Redis = &redisConf{}
)

type redisConf struct {
	DB        int      `json:"db"`
	Addrs     []string `json:"-"`
	Password  string   `json:"password"`
	KeyPrefix string   `json:"-"`
}

func (r *redisConf) init() error {
	if err := Get(constant.RedisConfKey).Scan(r); err != nil {
		return err
	}
	r.Addrs = Get(constant.RedisConfKey, "addrs").StringSlice([]string{"localhost:6379"})
	r.KeyPrefix = Service.Name
	return nil
}

func init() {
	Configs(Redis)
}
