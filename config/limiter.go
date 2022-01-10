package config

import (
	"github.com/blackdreamers/core/consts"
)

var (
	Limiter = &limiterConf{}
)

type limiterConf struct {
	Store string `json:"store"`
	// format:<limit>-<period>
	// 5 reqs/second: "5-S"
	// 10 reqs/minute: "10-M"
	// 1000 reqs/hour: "1000-H"
	// 2000 reqs/day: "2000-D"
	Limit string `json:"limit"`
}

func (l *limiterConf) init() error {
	if err := Get(consts.LimiterConfKey).Scan(l); err != nil {
		return err
	}
	if l.Store == "" {
		l.Store = consts.MemoryStore
	}
	if l.Limit == "" {
		l.Limit = "10-S"
	}
	return nil
}

func init() {
	Configs(Limiter)
}
