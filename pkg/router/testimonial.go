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
	testimonialService "github.com/hngprojects/hng_boilerplate_golang_web/services/testimonial"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Testimonial(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	testimonialServiceImp := testimonialService.NewTestimonialService(db.Postgresql.DB())
	controller := testimonial.Controller{Db: db, Logger: logger, Validator: validator, ExtReq: extReq, TestimonialService: testimonialServiceImp}

	squeezeURL := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql.DB(), models.RoleIdentity.User))
	{
		squeezeURL.POST("/testimonials", controller.Create)
	}

	publicGroup := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		publicGroup.GET("/testimonials/user/:user_id", controller.GetUserTestimonials)
	}

	return r
}
