package squeeze

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/squeeze"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Logger    *utility.Logger
	Validator *validator.Validate
	ExtReq    request.ExternalRequest
}

func (base *Controller) Create(c *gin.Context) {
	var req models.SqueezeUserReq

	if err := c.ShouldBind(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	reqData, code, err := service.ValidateSqueezeUserRequest(req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	squeezeUser, err := service.CreateSqueeze(base.Db.Postgresql, base.ExtReq, reqData)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), "failed to submit your request", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("your request has been received. you will get a template shortly")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "your request has been received. you will get a template shortly", squeezeUser)
	c.JSON(http.StatusCreated, rd)

}
