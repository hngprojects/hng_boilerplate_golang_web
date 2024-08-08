package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func (base *Controller) GetUserDataPrivacySettings(c *gin.Context) {
	var (
		userID = c.Param("user_id")
	)

	respData, code, err := service.GetUserDataPrivacySettings(userID, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "User data privacy settings retrieved successfully", respData)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) UpdateUserDataPrivacySettings(c *gin.Context) {
	var (
		userID = c.Param("user_id")
		req    = models.DataPrivacySettings{}
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

	respData, code, err := service.UpdateUserDataPrivacySettings(req, userID, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("user data privacy settings updated successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "User data privacy settings updated successfully", respData)
	c.JSON(http.StatusOK, rd)

}
