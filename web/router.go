package web

import (
	"net/http"

	"github.com/masibw/blog-server/usecase"

	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/web/handler"
)

func NewServer(postUC *usecase.PostUseCase, tagUC *usecase.TagUseCase) (e *gin.Engine) {
	e = gin.Default()

	postHandler := handler.NewPostHandler(postUC)
	tagHandler := handler.NewTagHandler(tagUC)

	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	v1 := e.Group("/api/v1")

	posts := v1.Group("/posts")
	posts.GET("", postHandler.GetPosts)
	posts.POST("", postHandler.StorePost)
	posts.GET(":id", postHandler.GetPost)
	posts.DELETE(":id", postHandler.DeletePost)

	tags := v1.Group("/tags")
	tags.GET("", tagHandler.GetTags)
	tags.POST("", tagHandler.StoreTag)
	tags.GET(":id", tagHandler.GetTag)
	tags.DELETE(":id", tagHandler.DeleteTag)
	return
}
