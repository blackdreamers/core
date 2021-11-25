package auth

import (
	_ "embed"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"

	"github.com/blackdreamers/core/db"
)

var (
	//go:embed rbac.conf
	config string
	e      *casbin.Enforcer
)

func Init() error {
	m, err := model.NewModelFromString(config)
	if err != nil {
		return err
	}

	a, err := gormadapter.NewAdapterByDBUseTableName(db.DB.Table("db_name.platform"), "", "rule")
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
