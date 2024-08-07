package test_testimonial

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/testimonial"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func SetupTestimonialTestRouter() (*gin.Engine, *testimonial.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tst.Setup()
	db := storage.Connection()
	validator := validator.New()

	testimonialController := &testimonial.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupTestimonialRoutes(r, testimonialController)
	return r, testimonialController
}

func SetupTestimonialRoutes(r *gin.Engine, testimonialController *testimonial.Controller) {
	r.POST(
		"/api/v1/testimonials", 
		middleware.Authorize(testimonialController.Db.Postgresql, models.RoleIdentity.User), 
		testimonialController.Create,
	)
}
