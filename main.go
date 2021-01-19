package main

import (
	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/database"
	"log"
	"net/http"
)

func main() {
	var err error
	engine := gin.Default()

	_, err = database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

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
