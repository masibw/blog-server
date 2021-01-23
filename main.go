package main

import (
	"errors"
	"os"
	"time"

	"github.com/masibw/blog-server/usecase"

	"github.com/masibw/blog-server/web"

	"github.com/masibw/blog-server/config"
	"github.com/masibw/blog-server/database"
	"github.com/masibw/blog-server/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func main() {
	logger := log.GetLogger()
	logger.Infof("Initialized logger")
	time.Local = time.FixedZone("JST", 9*60*60)
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

	postRepository := database.NewPostRepository(db)
	postUC := usecase.NewPostUseCase(postRepository)

	e := web.NewServer(postUC)

	if err := e.Run(":8080"); err != nil {
		if err != nil {
			logger.Fatal(err.Error())
		}
	}
}
