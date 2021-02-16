package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/masibw/blog-server/mock_usecase"

	"github.com/golang/mock/gomock"

	"github.com/gin-gonic/gin"
)

func TestImageHandler_StoreImage(t *testing.T) {
	tests := []struct {
		name                 string
		prepareMockImageUCFn func(mock *mock_usecase.MockImage)
		queryParam           string
		wantCode             int
	}{
		{
			name: "正常にタグを保存できる",
			prepareMockImageUCFn: func(mock *mock_usecase.MockImage) {
				mock.EXPECT().CreatePresignedURL(gomock.Any(), gomock.Any()).Return("url", nil)
			},
			queryParam: `objectName=image&contentType=image%2Fpng`,
			wantCode:   http.StatusCreated,
		},
		{
			name: "urlの作成に失敗した時はStatusInternalServerErrorエラーが返る",
			prepareMockImageUCFn: func(mock *mock_usecase.MockImage) {
				mock.EXPECT().CreatePresignedURL(gomock.Any(), gomock.Any()).Return("", errors.New("dummy error"))
			},
			queryParam: `objectName=image&contentType=image%2Fpng`,
			wantCode:   http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mu := mock_usecase.NewMockImage(ctrl)
			tt.prepareMockImageUCFn(mu)

			// HTTPRequestをテストするために必要な部分
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/images"+"?"+tt.queryParam, nil)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &ImageHandler{
				imageUC: mu,
			}
			p.GetPresignedURL(c)
			if w.Code != tt.wantCode {
				t.Errorf("StoreImage() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}
