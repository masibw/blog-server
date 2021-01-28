package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	post, err := p.postUC.CreatePost()
	if err != nil {
		logger.Errorf("store post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"post": post,
	})
}

// UpdatePost は PUT /posts/:id に対応するハンドラーです。
func (p *PostHandler) UpdatePost(c *gin.Context) {

	type request struct {
		ID           string    `json:"id" binding:"required"`
		Title        string    `json:"title"`
		ThumbnailURL string    `json:"thumbnailUrl" binding:"required"`
		Content      string    `json:"content" `
		Permalink    string    `json:"permalink" `
		IsDraft      *bool     `json:"isDraft" binding:"required"`
		CreatedAt    time.Time `json:"createdAt" binding:"required"`
		UpdatedAt    time.Time `json:"updatedAt" binding:"required"`
		PublishedAt  time.Time `json:"publishedAt" `
	}

	req := &request{}

	logger := log.GetLogger()
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Errorf("failed to bind", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postDTO := &dto.PostDTO{
		ID:           req.ID,
		Title:        req.Title,
		ThumbnailURL: req.ThumbnailURL,
		Content:      req.Content,
		Permalink:    req.Permalink,
		IsDraft:      req.IsDraft,
		CreatedAt:    req.CreatedAt,
		UpdatedAt:    req.UpdatedAt,
		PublishedAt:  req.PublishedAt,
	}

	post, err := p.postUC.UpdatePost(postDTO)

	if err != nil {
		if errors.Is(err, entity.ErrPermalinkAlreadyExisted) {
			logger.Debugf("update post already existed", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, entity.ErrPostNotFound) {
			logger.Debugf("update post not found", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logger.Errorf("update post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
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
	if c.Query("page") != "" && c.Query("page-size") != "" {
		var page int
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			logger.Errorf("page invalid, %v : %v", c.Query("page"), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pageSize, err = strconv.Atoi(c.Query("page-size"))
		if err != nil {
			logger.Errorf("page-size invalid, %v : %v", c.Query("page-size"), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if page == 0 {
			page = 1
		}

		offset = (page - 1) * pageSize
	}

	if c.Query("is-draft") != "" {
		isDraft, err := strconv.ParseBool(c.Query("is-draft"))
		if err != nil {
			logger.Errorf("is-draft invalid, %v : %v", c.Query("is-draft"), err)
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
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrPostNotFound.Error()})
			return
		}
		logger.Errorf("get posts", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError.Error()})
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
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrPostNotFound.Error()})
			return
		}
		logger.Errorf("get post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError.Error()})
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
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrPostNotFound.Error()})
			return
		}
		logger.Errorf("delete post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted",
	})
}
