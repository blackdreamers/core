package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/blackdreamers/core/config"
)

var (
	DB           *gorm.DB
	repositories []Repository
)

type Repository interface {
	FetchById() error
}

func Repositories(dbRepositories ...Repository) {
	repositories = append(repositories, dbRepositories...)
}

func Init() error {
	var err error

	dialect := mysql.New(
		mysql.Config{
			DSN: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
				config.Conf.DBUser,
				config.Conf.DBPassword,
				config.Conf.DBHost,
				config.Conf.DBPort,
				config.Service.DBName,
			),
			DefaultStringSize: 255, // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		},
	)

	cfg := &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLog(),
	}

	DB, err = gorm.Open(dialect, cfg)
	if err != nil {
		return err
	}

	for _, repository := range repositories {
		if err = DB.AutoMigrate(repository); err != nil {
			return err
		}
	}

	return nil
}
