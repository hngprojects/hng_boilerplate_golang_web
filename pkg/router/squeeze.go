package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/squeeze"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	SqueezeService "github.com/hngprojects/hng_boilerplate_golang_web/services/squeeze"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Squeeze(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	SqueezeServiceImp := SqueezeService.NewSqueezeUserService(db.Postgresql.DB())
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	controller := squeeze.Controller{Db: db, Logger: logger, Validator: validator, ExtReq: extReq, SqueezeService: SqueezeServiceImp}

	squeezeURL := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		squeezeURL.POST("/squeeze", controller.Create)
	}
	return r
}
