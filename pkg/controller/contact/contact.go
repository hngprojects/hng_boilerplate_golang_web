package contact

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/contact"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) GetAllContactUs(c *gin.Context) {

	contactsData, paginationResponse, code, err := service.GetAllContactUs(c, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Messages retrieved successfully", contactsData, paginationResponse)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetContactUsById(c *gin.Context) {

	var (
		reqID = c.Param("id")
	)

	contactData, err := service.GetContactUsById(reqID, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Message retrieved successfully", contactData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetContactUsByEmail(c *gin.Context) {

	var (
		reqEmail = c.Param("email")
	)

	contactData, err := service.GetContactUsByEmail(reqEmail, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Messages retrieved successfully", contactData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) DeleteContactUs(c *gin.Context) {

	var (
		reqID = c.Param("id")
	)

	code, err := service.DeleteContactUs(reqID, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Message deleted successfully", nil)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) AddToContactUs(c *gin.Context) {
	var (
		req = models.ContactUs{}
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

	err = service.AddToContactUs(&req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, rd)

		return
	}

	base.Logger.Info("message sent successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Message sent successfully", nil)
	c.JSON(http.StatusCreated, rd)

}
