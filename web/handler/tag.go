package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/masibw/blog-server/domain/entity"

	"github.com/masibw/blog-server/domain/dto"

	"github.com/masibw/blog-server/usecase"

	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/log"
)

type TagHandler struct {
	tagUC *usecase.TagUseCase
}

func NewTagHandler(tagUC *usecase.TagUseCase) *TagHandler {
	return &TagHandler{tagUC: tagUC}
}

// StoreTag は POST /tags に対応するハンドラーです。
func (p *TagHandler) StoreTag(c *gin.Context) {
	logger := log.GetLogger()
	tagDTO := &dto.TagDTO{}
	if err := c.ShouldBindJSON(tagDTO); err != nil {
		logger.Errorf("failed to bind", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag, err := p.tagUC.StoreTag(tagDTO)
	if err != nil {
		if errors.Is(err, entity.ErrTagNameAlreadyExisted) {
			logger.Debugf("store tag already tag name existed :%w", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": entity.ErrTagAlreadyExisted})
			return
		}
		logger.Errorf("store tag", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"tag": tag,
	})
}

// GetTags は POST /tags に対応するハンドラーです。
func (p *TagHandler) GetTags(c *gin.Context) {
	logger := log.GetLogger()
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
	tags, err := p.tagUC.GetTags(offset, pageSize)
	if err != nil {
		if errors.Is(err, entity.ErrTagNotFound) {
			logger.Debug("get tags not found", err)
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrTagNotFound})
			return
		}
		logger.Errorf("get tags", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
	})
}

func (p *TagHandler) GetTag(c *gin.Context) {
	logger := log.GetLogger()
	id := c.Param("id")
	tag, err := p.tagUC.GetTag(id)
	if err != nil {
		if errors.Is(err, entity.ErrTagNotFound) {
			logger.Debug("get tag not found", err)
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrTagNotFound})
			return
		}
		logger.Errorf("get tag", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tag": tag,
	})
}

func (p *TagHandler) DeleteTag(c *gin.Context) {
	logger := log.GetLogger()
	id := c.Param("id")
	err := p.tagUC.DeleteTag(id)
	if err != nil {
		if errors.Is(err, entity.ErrTagNotFound) {
			logger.Debug("delete tag not found", err)
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrTagNotFound})
			return
		}
		logger.Errorf("delete tag", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted",
	})
}
