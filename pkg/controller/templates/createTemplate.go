package templates

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	templateService "github.com/hngprojects/hng_boilerplate_golang_web/services/templates"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) CreateTemplate(c *gin.Context) {
	var template models.TemplateRequest

	if err := c.ShouldBindJSON(&template); err != nil {
		base.Logger.Error("Failed to parse request body")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}


	err := base.Validator.Struct(&template)
	if err != nil {
		base.Logger.Error("Validation failed")
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	tempData, err := templateService.CreateTemplate(base.Db, template)
	if err != nil {
		base.Logger.Error("Failed to create template")
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to create template", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Template created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated,"Template created successfully", tempData)
	c.JSON(http.StatusCreated, rd)
}
