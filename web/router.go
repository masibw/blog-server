package web

import (
	"net/http"

	"github.com/masibw/blog-server/usecase"

	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/web/handler"
)

func NewServer(postUC *usecase.PostUseCase) (e *gin.Engine) {
	e = gin.Default()

	postHandler := handler.NewPostHandler(postUC)

	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	v1 := e.Group("/api/v1")

	posts := v1.Group("/posts")
	posts.POST("", postHandler.StorePost)

	return
}
