package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
	r.Use(middleware.Metrics(config.GetConfig()))
	r.Use(middleware.GzipWithExclusion("/metrics"))
	r.MaxMultipartMemory = 3 << 20

	// routers
	ApiVersion := "api/v1"

	Health(r, ApiVersion, validator, db, logger)
	Seed(r, ApiVersion, validator, db, logger)
	Invite(r, ApiVersion, validator, db, logger)
	Testimonial(r, ApiVersion, validator, db, logger)
	Squeeze(r, ApiVersion, validator, db, logger)
	Blog(r, ApiVersion, validator, db, logger)
	Waitlist(r, ApiVersion, validator, db, logger)
	User(r, ApiVersion, validator, db, logger)
	Organisation(r, ApiVersion, validator, db, logger)
	Newsletter(r, ApiVersion, validator, db, logger)
	Product(r, ApiVersion, validator, db, logger)
	Auth(r, ApiVersion, validator, db, logger)
	JobPost(r, ApiVersion, validator, db, logger)
	FAQ(r, ApiVersion, validator, db, logger)
	SuperAdmin(r, ApiVersion, validator, db, logger)
	Category(r, ApiVersion, validator, db, logger)
	Notification(r, ApiVersion, validator, db, logger)
	Template(r, ApiVersion, validator, db, logger)
	Python(r, ApiVersion, validator, db, logger)
	HelpCenter(r, ApiVersion, validator, db, logger)
	Profile(r, ApiVersion, validator, db, logger)
	Contact(r, ApiVersion, validator, db, logger)
	NotificationSettings(r, ApiVersion, validator, db, logger)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 200,
			"message":     "HNGi Golang Boilerplate",
			"status":      http.StatusOK,
		})
	})

	r.StaticFile("/swagger.yaml", "static/swagger.yaml")
	url := ginSwagger.URL("/swagger.yaml")
	r.GET("/api/docs/*any", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'sha256-2TOI2ugkuROHHfKZr6kdGv+XxhrVUI8uHycXqXUIR4g='; img-src 'self' data:;")
		ginSwagger.WrapHandler(swaggerFiles.Handler, url)(c)
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"name":        "Not Found",
			"message":     "Page not found.",
			"status_code": 404,
			"status":      http.StatusNotFound,
		})
	})

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
