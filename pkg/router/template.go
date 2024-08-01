package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/templates"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Template(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	template := templates.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	templateUrl := r.Group(fmt.Sprintf("%v/", ApiVersion), middleware.Authorize(db.Postgresql))
	{
		templateUrl.POST("/template", template.CreateTemplate)
		templateUrl.GET("/template", template.GetTemplates)
		templateUrl.GET("/template/:id", template.GetTemplate)
		templateUrl.DELETE("/template/:id", template.DeleteTemplate)
		templateUrl.PUT("/template/:id", template.UpdateTemplate)
	}
	return r
}
