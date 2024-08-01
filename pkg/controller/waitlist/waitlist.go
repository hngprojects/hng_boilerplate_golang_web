package waitlist

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/waitlist"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"net/http"
)

type Controller struct {
	DB        *storage.Database
	Logger    *utility.Logger
	Validator *validator.Validate
}

func (base *Controller) GetWaitLists(c *gin.Context) {
	waitlistData, paginationResponse, code, err := service.GetWaitLists(c, base.DB.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "Waitlist retrieved successfully", waitlistData, paginationResponse)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) Create(c *gin.Context) {
	var (
		req = models.CreateWaitlistUserRequest{}
	)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		v := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, v)
		return
	}
	err = base.Validator.Struct(&req)
	if err != nil {
		v := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "The given data was invalid", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, v)
		return
	}
	data, code, err := service.SignupWaitlistUserService(base.DB.Postgresql, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}
	rd := utility.BuildSuccessResponse(code, "waitlist signup successful", data)
	c.JSON(code, rd)
}
