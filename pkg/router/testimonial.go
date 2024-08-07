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
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Testimonial(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	controller := testimonial.Controller{Db: db, Logger: logger, Validator: validator, ExtReq: extReq}

	squeezeURL := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql, models.RoleIdentity.User))
	{
		squeezeURL.POST("/testimonials", controller.Create)
	}
	return r
}
