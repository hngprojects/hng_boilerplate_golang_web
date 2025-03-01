package test_testimonial

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/testimonial"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/testimonial"
	testimonialService "github.com/hngprojects/hng_boilerplate_golang_web/services/testimonial"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func SetupTestimonialTestRouter() (*gin.Engine, *testimonial.Controller) {
    gin.SetMode(gin.TestMode)

    logger := tst.Setup()
    db := storage.Connection()
    validator := validator.New()
    testimonialService := service.NewTestimonialService(db.Postgresql.DB())

    testimonialController := &testimonial.Controller{
        Db:        db,
        Validator: validator,
        Logger:    logger,
        TestimonialSvc:   testimonialService,  
    }
	testimonialService := testimonialService.NewTestimonialService(db.Postgresql.DB())
	testimonialController := &testimonial.Controller{
		Db:                 db,
		Validator:          validator,
		Logger:             logger,
		TestimonialService: testimonialService,
	}

    r := gin.Default()
    SetupTestimonialRoutes(r, testimonialController)
    return r, testimonialController
}


func SetupTestimonialRoutes(r *gin.Engine, testimonialController *testimonial.Controller) {
	r.POST(
		"/api/v1/testimonials",
		middleware.Authorize(testimonialController.Db.Postgresql.DB(), models.RoleIdentity.User),
		testimonialController.Create,
	)
	r.GET("/api/v1/testimonials/user/:user_id", testimonialController.GetUserTestimonials)
}
