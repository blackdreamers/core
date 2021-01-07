package db

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/blackdreamers/core/config"
)

const dbConfKey = "database"

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
	var dbConf map[string]string
	dbConf = config.Get(dbConfKey).StringMap(dbConf)

	dialect := mysql.New(
		mysql.Config{
			DSN: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
				dbConf["user"],
				dbConf["password"],
				dbConf["host"],
				dbConf["port"],
				config.Service.Get("dbname").String(config.Service.Name),
			),
			DefaultStringSize: 255, // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		},
	)

	cfg := &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
	if os.Getenv("MODE") == "debug" {
		cfg.Logger = logger.Default.LogMode(logger.Info)
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
