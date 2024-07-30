package test_faq

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/faq"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func SetupFAQTestRouter() (*gin.Engine, *faq.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tst.Setup()
	db := storage.Connection()
	validator := validator.New()

	faqController := &faq.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupFAQRoutes(r, faqController)
	return r, faqController
}

func SetupFAQRoutes(r *gin.Engine, faqController *faq.Controller) {
	r.POST("/api/v1/faq", middleware.Authorize(faqController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		faqController.AddToFaq)
	r.GET("/api/v1/faq", faqController.GetFaq)
	r.DELETE("/api/v1/faq/:id", middleware.Authorize(faqController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		faqController.DeleteFaq)
	r.PUT("/api/v1/faq/:id", middleware.Authorize(faqController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		faqController.UpdateFaq)
}
