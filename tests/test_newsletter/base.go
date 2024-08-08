package test_newsletter

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/newsletter"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func SetupNewsLetterTestRouter() (*gin.Engine, *newsletter.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tst.Setup()
	db := storage.Connection()
	validator := validator.New()

	newsController := &newsletter.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupNewsLetterRoutes(r, newsController)
	return r, newsController
}

func SetupNewsLetterRoutes(r *gin.Engine, newsController *newsletter.Controller) {
	r.POST("/api/v1/newsletter-subscription", newsController.SubscribeNewsLetter)
	r.GET("/api/v1/newsletter-subscription", middleware.Authorize(newsController.Db.Postgresql,
		models.RoleIdentity.SuperAdmin),
		newsController.GetNewsLetters)
	r.DELETE("/api/v1/newsletter-subscription/:id", middleware.Authorize(newsController.Db.Postgresql,
		models.RoleIdentity.SuperAdmin),
		newsController.DeleteNewsLetter)
	r.GET("/api/v1/newsletter-subscription/deleted",
		middleware.Authorize(newsController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		newsController.GetDeletedNewsLetters)
	r.PATCH("/api/v1/newsletter-subscription/restore/:id",
		middleware.Authorize(newsController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		newsController.RestoreDeletedNewsLetter)
}
