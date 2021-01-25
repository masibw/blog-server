package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/masibw/blog-server/domain/entity"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"post": post,
	})
}

// GetPosts は POST /posts に対応するハンドラーです。
func (p *PostHandler) GetPosts(c *gin.Context) {
	logger := log.GetLogger()

	conditions := make([]string, 0)
	params := make([]interface{}, 0)
	var offset int
	var pageSize int
	var err error

	// ページネーションの設定
	if c.Query("page") != "" && c.Query("page_size") != "" {
		var page int
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			logger.Errorf("page invalid, %v : %v", c.Query("page"), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pageSize, err = strconv.Atoi(c.Query("page_size"))
		if err != nil {
			logger.Errorf("page_size invalid, %v : %v", c.Query("page_size"), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if page == 0 {
			page = 1
		}

		offset = (page - 1) * pageSize
	}

	if c.Query("is_draft") != "" {
		isDraft, err := strconv.ParseBool(c.Query("is_draft"))
		if err != nil {
			logger.Errorf("is_draft invalid, %v : %v", c.Query("is_draft"), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		conditions = append(conditions, " is_draft = ? ")
		params = append(params, isDraft)
	}
	condition := strings.Join(conditions, "")
	posts, err := p.postUC.GetPosts(offset, pageSize, condition, params)
	if err != nil {
		if errors.Is(err, entity.ErrPostNotFound) {
			logger.Debug("get posts not found", err)
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrPostNotFound})
			return
		}
		logger.Errorf("get posts", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func (p *PostHandler) GetPost(c *gin.Context) {
	logger := log.GetLogger()
	id := c.Param("id")
	post, err := p.postUC.GetPost(id)
	if err != nil {
		if errors.Is(err, entity.ErrPostNotFound) {
			logger.Debug("get post not found", err)
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrPostNotFound})
			return
		}
		logger.Errorf("get post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func (p *PostHandler) DeletePost(c *gin.Context) {
	logger := log.GetLogger()
	id := c.Param("id")
	err := p.postUC.DeletePost(id)
	if err != nil {
		if errors.Is(err, entity.ErrPostNotFound) {
			logger.Debug("delete post not found", err)
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrPostNotFound})
			return
		}
		logger.Errorf("delete post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted",
	})
}
