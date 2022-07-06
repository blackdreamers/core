package config

import "github.com/blackdreamers/core/consts"

var (
	ES = &esConf{}
)

type esConf struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (e *esConf) init() error {
	if err := Get(consts.EsKey).Scan(e); err != nil {
		return err
	}
	return nil
}

func init() {
	Configs(ES)
}
