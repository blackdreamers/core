package config

import (
	"github.com/blackdreamers/core/consts"
)

var (
	Session = &sessionConf{}
)

type sessionConf struct {
	Store    string `json:"store"`
	Secret   string `json:"secret"`
	MaxAge   int    `json:"max_age"`
	HttpOnly bool   `json:"http_only"`
	DB       int    `json:"db"`
}

func (s *sessionConf) init() error {
	if err := Get(consts.SessionConfKey).Scan(s); err != nil {
		return err
	}
	if s.Store == "" {
		s.Store = consts.MemoryStore
	}
	if s.Secret == "" {
		s.Secret = "daydream"
	}
	if s.MaxAge == 0 {
		s.MaxAge = 604800
	}
	return nil
}

func init() {
	Configs(Session)
}
