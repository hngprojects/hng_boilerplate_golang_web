package waitlist

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/waitlist"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	DB        *storage.Database
	Logger    *utility.Logger
	Validator *validator.Validate
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

	data, code, err := waitlist.SignupWaitlistUserService(base.DB.Postgresql, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(code, "waitlist signup successful", data)
	c.JSON(code, rd)
}
