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
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxLifeTime  int    `json:"max_life_time"`
	MaxIdleTime  int    `json:"max_idle_time"`
}

func (db *dbConf) init() error {
	if err := Get(consts.DBConfKey).Scan(db); err != nil {
		return err
	}
	if db.MaxOpenConns == 0 {
		db.MaxOpenConns = 100
	}
	if db.MaxIdleConns == 0 {
		db.MaxIdleConns = 25
	}
	return nil
}

func init() {
	Configs(DB)
}
