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

// GetRegion godoc
// @Summary Get all regions
// @Description Retrieve all regions
// @Tags Superadmin
// @Accept json
// @Produce json
// @Success 200 {object} utility.Response
// @Failure 400 {object} utility.Response
// @Router /superadmin/regions [get]
func (base *Controller) GetRegion(c *gin.Context) {

	regionData, err := service.GetRegions(base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Regions retrieved successfully", regionData)
	c.JSON(http.StatusOK, rd)

}

// GetTimeZone godoc
// @Summary Get all timezones
// @Description Retrieve all timezones
// @Tags Superadmin
// @Accept json
// @Produce json
// @Success 200 {object} utility.Response
// @Failure 400 {object} utility.Response
// @Router /superadmin/timezones [get]
func (base *Controller) GetTimeZone(c *gin.Context) {

	timezoneData, err := service.GetTimeZones(base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Timezones retrieved successfully", timezoneData)
	c.JSON(http.StatusOK, rd)

}

// GetLanguage godoc
// @Summary Get all languages
// @Description Retrieve all languages
// @Tags Superadmin
// @Accept json
// @Produce json
// @Success 200 {object} utility.Response
// @Failure 400 {object} utility.Response
// @Router /superadmin/languages [get]
func (base *Controller) GetLanguage(c *gin.Context) {

	languageData, err := service.GetLanguages(base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Languages retrieved successfully", languageData)
	c.JSON(http.StatusOK, rd)

}

// AddToRegion godoc
// @Summary Add a new region
// @Description Add a new region to the system
// @Tags Superadmin
// @Accept json
// @Produce json
// @Param region body models.Region true "Region details"
// @Success 201 {object} utility.Response
// @Failure 400,422 {object} utility.Response
// @Router /superadmin/regions [post]
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

// AddToTimeZone godoc
// @Summary Add a new timezone
// @Description Add a new timezone to the system
// @Tags Superadmin
// @Accept json
// @Produce json
// @Param timezone body models.Timezone true "Timezone details"
// @Success 201 {object} utility.Response
// @Failure 400,422 {object} utility.Response
// @Router /superadmin/timezones [post]
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

// AddToLanguage godoc
// @Summary Add a new language
// @Description Add a new language to the system
// @Tags Superadmin
// @Accept json
// @Produce json
// @Param language body models.Language true "Language details"
// @Success 201 {object} utility.Response
// @Failure 400,422 {object} utility.Response
// @Router /superadmin/languages [post]
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
