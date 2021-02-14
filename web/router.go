package web

import (
	"net/http"
	"os"
	"time"

	"github.com/Songmu/flextime"
	"github.com/gin-contrib/cors"
	"github.com/masibw/blog-server/constant"
	"github.com/masibw/blog-server/log"

	"github.com/masibw/blog-server/domain/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/config"
	"github.com/masibw/blog-server/usecase"
	"github.com/masibw/blog-server/web/handler"
)

type login struct {
	MailAddress string `form:"mailAddress" json:"mailAddress" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
}

func NewServer(postUC *usecase.PostUseCase, tagUC *usecase.TagUseCase, authMW *AuthMiddleware, postsTagsService *service.PostsTagsService) (e *gin.Engine) {
	logger := log.GetLogger()
	e = gin.New()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	e.Use(cors.New(corsConfig))

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "authenticated zone",
		Key:             []byte(os.Getenv("AUTH_KEY")),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     constant.IdentityKey,
		PayloadFunc:     authMW.PayloadFunc,
		IdentityHandler: authMW.IdentityHandler,
		Authenticator:   authMW.Authenticate,
		Authorizator:    authMW.Authorize,
		Unauthorized:    authMW.UnAuthorize,
		TokenLookup:     "cookie: jwt",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// READMEには載ってなかったけど存在しているみたい
		CookieSameSite: http.SameSiteStrictMode,
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc:       flextime.Now,
		SendCookie:     true,
		SecureCookie:   !config.IsLocal(),
		CookieHTTPOnly: true, // JS can't modify
	})

	if err != nil {
		logger.Fatal("JWT Error:" + err.Error())
	}

	postHandler := handler.NewPostHandler(postUC, postsTagsService)
	tagHandler := handler.NewTagHandler(tagUC)

	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	v1 := e.Group("/api/v1")

	v1.POST("/login", authMiddleware.LoginHandler)
	v1.POST("/logout", authMiddleware.LogoutHandler)

	posts := v1.Group("/posts")
	posts.GET("", postHandler.GetPosts)
	posts.GET(":permalink", postHandler.GetPost)

	posts.Use(authMiddleware.MiddlewareFunc())
	{
		posts.POST("", postHandler.StorePost)
		posts.PUT(":id", postHandler.UpdatePost)
		posts.DELETE(":id", postHandler.DeletePost)
	}

	tags := v1.Group("/tags")
	tags.GET("", tagHandler.GetTags)
	tags.GET(":id", tagHandler.GetTag)
	tags.Use(authMiddleware.MiddlewareFunc())
	{
		tags.POST("", tagHandler.StoreTag)
		tags.DELETE(":id", tagHandler.DeleteTag)
	}
	return
}
