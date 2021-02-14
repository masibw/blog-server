package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/masibw/blog-server/domain/service"

	"github.com/masibw/blog-server/domain/entity"

	"github.com/masibw/blog-server/domain/dto"

	"github.com/masibw/blog-server/usecase"

	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/log"
)

type PostHandler struct {
	postUC           *usecase.PostUseCase
	postsTagsService *service.PostsTagsService
}

func NewPostHandler(postUC *usecase.PostUseCase, postsTagsservice *service.PostsTagsService) *PostHandler {
	return &PostHandler{
		postUC:           postUC,
		postsTagsService: postsTagsservice,
	}
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

	type postReq struct {
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

	type request struct {
		Post postReq  `json:"post" binding:"dive"`
		Tags []string `json:"tags"`
	}

	req := &request{}
	logger := log.GetLogger()
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Errorf("failed to bind", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postDTO := &dto.PostDTO{
		ID:           req.Post.ID,
		Title:        req.Post.Title,
		ThumbnailURL: req.Post.ThumbnailURL,
		Content:      req.Post.Content,
		Permalink:    req.Post.Permalink,
		IsDraft:      req.Post.IsDraft,
		CreatedAt:    req.Post.CreatedAt,
		UpdatedAt:    req.Post.UpdatedAt,
		PublishedAt:  req.Post.PublishedAt,
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

	if req.Tags == nil || len(req.Tags) == 0 {
		logger.Debug("tags nil")
		c.JSON(http.StatusOK, gin.H{
			"post": post,
			"tags": nil,
		})
		return
	}
	var tags []*entity.Tag
	tags, err = p.postsTagsService.LinkPostTags(req.Post.ID, req.Tags)
	if err != nil {
		if errors.Is(err, entity.ErrPostsTagsAlreadyExisted) {
			logger.Debugf("update post tags already exists", err)
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
		"tags": tags,
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

	if c.Query("tag") != "" {
		tagName := c.Query("tag")
		conditions = append(conditions, "tags.name = ?")
		params = append(params, tagName)
	}

	var sortCondition string
	if c.Query("sort") != "" {
		sort := c.Query("sort")
		order := " asc"
		if strings.HasPrefix(sort, "-") {
			order = " desc"
			sort = sort[1:]
		}

		var column string
		switch sort {
		case "id":
			column = "id"
		case "title":
			column = "title"
		case "thumbnailUrl", "thumbnail_url":
			column = "thumbnail_url"
		case "content":
			column = "content"
		case "permalink":
			column = "permalink"
		case "isDraft", "is_draft":
			column = "is_draft"
		case "createdAt", "created_at":
			column = "created_at"
		case "updatedAt", "updated_at":
			column = "updated_at"
		case "publishedAt", "published_at":
			column = "published_at"
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": entity.ErrPostColumnNotFound.Error()})
			return
		}
		sortCondition = column + order
	} else {
		sortCondition = "id asc"
	}

	condition := strings.Join(conditions, " AND ")
	posts, count, err := p.postUC.GetPosts(offset, pageSize, condition, params, sortCondition)
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
		"count": count,
	})
}

func (p *PostHandler) GetPost(c *gin.Context) {
	logger := log.GetLogger()
	permalink := c.Param("permalink")

	isMarkdown := false
	if c.Query("is-markdown") != "" {
		var err error
		isMarkdown, err = strconv.ParseBool(c.Query("is-markdown"))
		if err != nil {
			logger.Errorf("is-markdown invalid, %v : %v", c.Query("is-markdown"), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	post, err := p.postUC.GetPost(permalink, isMarkdown)
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
