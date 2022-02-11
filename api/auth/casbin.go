package auth

import (
	_ "embed"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"

	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/db"
)

const (
	ruleDBName    = "auth"
	ruleTableName = "rule"
)

var (
	//go:embed model.conf
	conf string
	e    *casbin.Enforcer
)

func Init() error {
	m, err := model.NewModelFromString(conf)
	if err != nil {
		return err
	}

	var a *gormadapter.Adapter
	if config.Service.EnableDB && config.Service.Private {
		a, err = gormadapter.NewAdapterByDBUseTableName(db.DB, ruleDBName, ruleTableName)
	} else {
		a, err = gormadapter.NewAdapter(
			"mysql",
			fmt.Sprintf("%v:%v@tcp(%v:%v)/",
				config.DB.User,
				config.DB.Password,
				config.DB.Host,
				config.DB.Port,
			),
			ruleDBName,
			ruleTableName,
		)
	}
	if err != nil {
		return err
	}

	e, err = casbin.NewEnforcer(m, a)
	if err != nil {
		return err
	}

	if err = e.LoadPolicy(); err != nil {
		return err
	}

	return nil
}

func Enforce(vs ...interface{}) (bool, error) {
	return e.Enforce(vs...)
}

func LoadPolicy() error {
	return e.LoadPolicy()
}

func SavePolicy() error {
	return e.SavePolicy()
}
