package newsletter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) SubscribeNewsLetter(c *gin.Context) {
	var (
		req = models.NewsLetter{}
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

	if postgresql.CheckExists(base.Db.Postgresql, &req, "email = ?", req.Email) {
		rd := utility.BuildErrorResponse(http.StatusConflict, "error", "Email already subscribed", nil, nil)
		c.JSON(http.StatusConflict, rd)
		return
	}

	if err := postgresql.CreateOneRecord(base.Db.Postgresql, &req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to subscribe", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("subscribed successfully")

	rd := utility.BuildSuccessResponse2(http.StatusCreated, "subscribed successfully")
	c.JSON(http.StatusCreated, rd)

}
