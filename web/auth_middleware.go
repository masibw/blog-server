package web

import (
	"unsafe"

	"github.com/masibw/blog-server/constant"

	"github.com/masibw/blog-server/log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/domain/dto"
	"github.com/masibw/blog-server/usecase"
	"golang.org/x/crypto/bcrypt"
)

type AuthMiddleware struct {
	identityKey string
	userUC      *usecase.UserUseCase
}

func NewAuthMiddleware(userUC *usecase.UserUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		identityKey: constant.IdentityKey,
		userUC:      userUC,
	}
}

func (m *AuthMiddleware) Authenticate(c *gin.Context) (interface{}, error) {
	logger := log.GetLogger()

	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	mailAddress := loginVals.MailAddress
	password := *(*[]byte)(unsafe.Pointer(&loginVals.Password))

	user, err := m.userUC.GetUserByMailAddress(mailAddress)
	if err != nil {
		logger.Errorf("admin user not found mailAddress= %v", loginVals.MailAddress, err)
		return nil, jwt.ErrFailedAuthentication
	}

	diff := bcrypt.CompareHashAndPassword([]byte(user.Password), password)
	if user.MailAddress == mailAddress && diff == nil {
		return user, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

func (m *AuthMiddleware) Authorize(data interface{}, c *gin.Context) bool {
	// User = 管理者用ユーザーなので認可
	if _, ok := data.(*dto.UserDTO); ok {
		return true
	}
	return false
}

func (m *AuthMiddleware) UnAuthorize(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func (m *AuthMiddleware) PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*dto.UserDTO); ok {
		return jwt.MapClaims{
			m.identityKey: v.MailAddress,
		}
	}
	return jwt.MapClaims{}
}

func (m *AuthMiddleware) IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	if v, ok := claims[m.identityKey]; ok {
		return &dto.UserDTO{
			MailAddress: v.(string),
		}
	}
	return nil
}
