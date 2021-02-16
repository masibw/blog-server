package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/domain/entity"
	"github.com/masibw/blog-server/log"
	"github.com/masibw/blog-server/usecase"
)

type ImageHandler struct {
	imageUC usecase.Image
}

func NewImageHandler(imageUC usecase.Image) *ImageHandler {
	return &ImageHandler{imageUC: imageUC}
}

func (i *ImageHandler) GetPresignedURL(c *gin.Context) {
	logger := log.GetLogger()
	var fileName string
	if c.Query("objectName") != "" {
		fileName = c.Query("objectName")
	}
	var contentType string
	if c.Query("contentType") != "" {
		contentType = c.Query("contentType")
	}

	url, err := i.imageUC.CreatePresignedURL(&fileName, &contentType)
	if err != nil {
		logger.Errorf("create presigned url", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": entity.ErrInternalServerError.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"signedUrl": url,
	})
}
