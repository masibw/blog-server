package database

import (
	"database/sql"

	"github.com/masibw/blog-server/config"
	"github.com/masibw/blog-server/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func NewTestDB() *gorm.DB {
	logger := zapgorm2.New(log.GetPureLogger())
	logger.SetAsDefault()
	db, err := gorm.Open(mysql.Open(config.DSN()), &gorm.Config{Logger: logger})
	if err != nil {
		panic(err.Error())
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()

	if err != nil {
		panic(err.Error())
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(100)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	if err = sqlDB.Ping(); err != nil {
		panic(err.Error())
	}
	return db
}
