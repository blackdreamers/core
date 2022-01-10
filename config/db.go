package config

import "github.com/blackdreamers/core/consts"

var (
	DB = &dbConf{}
)

type dbConf struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	LowThreshold int    `json:"low_threshold"`
}

func (db *dbConf) init() error {
	if err := Get(consts.DBConfKey).Scan(db); err != nil {
		return err
	}
	return nil
}

func init() {
	Configs(DB)
}
