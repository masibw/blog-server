package main

import (
	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/config"
	"github.com/masibw/blog-server/database"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func main() {

	m, err := migrate.New("file://"+os.Getenv("MIGRATION_FILE"), "mysql://"+config.PureDSN())
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})
	err = engine.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}
