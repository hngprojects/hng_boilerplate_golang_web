package seed

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/ping"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) Get1(c *gin.Context) {
	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	var user models.User
	if err := base.Db.Postgresql.Preload("Profile").Preload("Products").Preload("Organisations").Where("email = ?", "john@example.com").First(&user).Error; err != nil {
		fmt.Println(err)
		return
	}

	base.Logger.Info("user1 fetched successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "", user)

	c.JSON(http.StatusOK, rd)
}

func (base *Controller) Get2(c *gin.Context) {
	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	var user models.User
	if err := base.Db.Postgresql.Preload("Profile").Preload("Products").Preload("Organisations").Where("email = ?", "jane@example.com").First(&user).Error; err != nil {
		fmt.Println(err)
		return
	}

	base.Logger.Info("user2 fetched successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "", user)

	c.JSON(http.StatusOK, rd)
}
