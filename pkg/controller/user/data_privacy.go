package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
