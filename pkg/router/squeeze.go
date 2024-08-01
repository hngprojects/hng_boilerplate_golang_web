package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/squeeze"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Squeeze(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	controller := squeeze.Controller{DB: db, Logger: logger, Validator: validator}

	squeezeURL := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		squeezeURL.POST("/squeeze", controller.Create)
	}
	return r
}
