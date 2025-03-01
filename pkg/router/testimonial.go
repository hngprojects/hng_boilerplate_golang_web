package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/testimonial"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/testimonial"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Testimonial(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}

	testimonialService := service.NewTestimonialService(db.Postgresql.DB())

	controller := testimonial.Controller{
		Db:            db,
		Logger:        logger,
		Validator:     validator,
		ExtReq:        extReq,
		TestimonialSvc: testimonialService,
	}

	publicGroup := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql.DB()))
	{
		publicGroup.GET("/testimonials/user/:user_id", controller.GetUserTestimonials)
	}

	protectedGroup := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql.DB(), models.RoleIdentity.User))
	{
		protectedGroup.POST("/testimonials", controller.Create)
	}

	return r
}


