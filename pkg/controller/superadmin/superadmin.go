package superadmin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/superadmin"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) AddToRegion(c *gin.Context) {
	var (
		req = models.Region{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed",
			utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	err = service.AddToRegion(&req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, rd)

		return
	}

	base.Logger.Info("region added successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Region added successfully", nil)
	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) AddToTimeZone(c *gin.Context) {
	var (
		req = models.Timezone{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed",
			utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	err = service.AddToTimeZone(&req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, rd)

		return
	}

	base.Logger.Info("timezone added successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Timezone added successfully", nil)
	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) AddToLanguage(c *gin.Context) {
	var (
		req = models.Language{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed",
			utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	err = service.AddToLanguage(&req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, rd)

		return
	}

	base.Logger.Info("language added successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Language added successfully", nil)
	c.JSON(http.StatusCreated, rd)

}
