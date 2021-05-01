package database

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"gorm.io/gorm"

	"github.com/masibw/blog-server/config"
	"github.com/masibw/blog-server/log"
	"gorm.io/driver/mysql"
	"moul.io/zapgorm2"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	var mig *migrate.Migrate
	mig, err = migrate.New("file://"+os.Getenv("MIGRATION_FILE"), "mysql://"+config.PureDSN())
	if err != nil {
		panic(err)
	}
	if err = mig.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}

	db = NewTestDB()
	code := m.Run()
	os.Exit(code)
}

func NewTestDB() *gorm.DB {
	logger := zapgorm2.New(log.GetPureLogger())
	logger.SetAsDefault()
	var err error
	db, err = gorm.Open(mysql.Open(config.DSN()), &gorm.Config{Logger: logger})
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
