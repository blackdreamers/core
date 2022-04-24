package test

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/blackdreamers/core/db"
	"github.com/blackdreamers/core/env"
	"github.com/blackdreamers/core/logger"
)

func Init() {
	initDB()
}

func initDB() {
	var err error

	if err = logger.Init(); err != nil {
		panic(err)
	}

	dialect := mysql.New(
		mysql.Config{
			DSN: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
				env.GetString("DB_USER", "test"),
				env.GetString("DB_PASSWORD", "test"),
				env.GetString("DB_HOST", "localhost"),
				env.GetString("DB_PORT", "3306"),
				env.GetString("DB_NAME", "test"),
			),
			DefaultStringSize:         255,
			SkipInitializeWithVersion: true,
		},
	)

	cfg := &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	db.DB, err = gorm.Open(dialect, cfg)
	if err != nil {
		panic(err)
	}

	sqldb, err := db.DB.DB()
	if err != nil {
		panic(err)
	}
	maxOpenConns, err := env.GetInt("DB_MAX_OPEN_CONNS", 100)
	if err != nil {
		panic(err)
	}
	maxIdleConns, err := env.GetInt("DB_MAX_IDLE_CONNS", 25)
	if err != nil {
		panic(err)
	}
	maxLifeTime, err := env.GetInt("DB_MAX_LIFE_TIME", 0)
	if err != nil {
		panic(err)
	}
	maxIdleTime, err := env.GetInt("DB_MAX_IDLE_TIME", 0)
	if err != nil {
		panic(err)
	}
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxIdleConns)
	if maxLifeTime != 0 {
		sqldb.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Second)
	}
	if maxIdleTime != 0 {
		sqldb.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Second)
	}

	for _, repository := range db.Repositories() {
		if err = db.DB.Migrator().DropTable(repository); err != nil {
			panic(err)
		}
		if err = db.DB.AutoMigrate(repository); err != nil {
			panic(err)
		}
	}
}
