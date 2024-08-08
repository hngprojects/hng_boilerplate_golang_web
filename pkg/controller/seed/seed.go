package seed

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/seed"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) GetUser(c *gin.Context) {
	//get the user_id from the URL
	userIDStr := c.Param("user_id")

	user, err := seed.GetUser(userIDStr, base.Db.Postgresql)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), err, nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	base.Logger.Info("user fetched successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "", user)

	c.JSON(http.StatusOK, rd)
}
