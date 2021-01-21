package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/config"
	"github.com/masibw/blog-server/database"
	"github.com/masibw/blog-server/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

type User struct {
	ID             string
	MailAddress    string
	Password       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastLoggedinAt time.Time
}

func main() {
	logger := log.GetLogger()
	logger.Infof("Initialized logger")

	m, err := migrate.New("file://"+os.Getenv("MIGRATION_FILE"), "mysql://"+config.PureDSN())
	if err != nil {
		logger.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			logger.Fatal(err)
		}
		logger.Infof("there were no changes to the schema")
	} else {
		logger.Info("updated db schema")
	}

	db, err := database.NewDB()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("db conn", db)
	db.Create(&User{ID: "aaaaaaaaaaaaaaaaaaaaaaaaa", MailAddress: "hoge@exmaple.com", Password: "hoge", LastLoggedinAt: time.Now()})
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	if err := engine.Run(":8080"); err != nil {
		if err != nil {
			logger.Fatal(err.Error())
		}
	}
}
