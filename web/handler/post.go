package handler

import (
	"net/http"

	"github.com/masibw/blog-server/domain/dto"

	"github.com/masibw/blog-server/usecase"

	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/log"
)

type PostHandler struct {
	postUC *usecase.PostUseCase
}

func NewPostHandler(postUC *usecase.PostUseCase) *PostHandler {
	return &PostHandler{postUC: postUC}
}

// StorePost は POST /posts に対応するハンドラーです。
func (p *PostHandler) StorePost(c *gin.Context) {
	logger := log.GetLogger()
	postDTO := &dto.PostDTO{}
	if err := c.ShouldBindJSON(postDTO); err != nil {
		logger.Errorf("failed to bind", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, err := p.postUC.StorePost(postDTO)
	if err != nil {
		logger.Errorf("store post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"post": post,
	})
}

// GetPosts は POST /posts に対応するハンドラーです。
func (p *PostHandler) GetPosts(c *gin.Context) {
	logger := log.GetLogger()

	posts, err := p.postUC.GetPosts()
	if err != nil {
		logger.Errorf("get posts", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}
