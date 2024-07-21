package router

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
<<<<<<< HEAD

=======
	"github.com/hngprojects/hng_boilerplate_golang_web/auth"
>>>>>>> feature/google_OAuth2
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Setup(logger *utility.Logger, validator *validator.Validate, db *storage.Database, appConfiguration *config.App) *gin.Engine {
	if appConfiguration.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// Middlewares
	// r.Use(gin.Logger())
	r.ForwardedByClientIP = true
	r.SetTrustedProxies(config.GetConfig().Server.TrustedProxies)
	r.Use(middleware.Security())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.MaxMultipartMemory = 1 << 20 // 1MB

	// routers
	ApiVersion := "api/v1"

	Health(r, ApiVersion, validator, db, logger)
	Seed(r, ApiVersion, validator, db, logger)
	User(r, ApiVersion, validator, db, logger)
	Organisation(r, ApiVersion, validator, db, logger)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "HNGi Golang Boilerplate",
			"status":  http.StatusOK,
		})
	})

	//OAuth implementation by BlacAc3
	r.GET("/api/v1/auth/login/google", auth.Handle_Google_Login)
	r.GET("/api/v1/auth/callback/google", auth.Handle_Google_Callback)
	r.POST("/api/v1/auth/token/refresh", auth.Handle_Token_Refresh)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"name":    "Not Found",
			"message": "Page not found.",
			"code":    404,
			"status":  http.StatusNotFound,
		})
	})

	return r
}
