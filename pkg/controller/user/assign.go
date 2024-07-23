package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hngprojects/hng_boilerplate_golang_web/services/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func (base *Controller) AssignRoleToUser(c *gin.Context) {

	userID := c.Param("user_id")
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userData, err := service.ReplaceUserRole(userID, roleID, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), nil, nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	respData, code, err := auth.LoginUser(req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("user login successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "user login successfully", respData)
	c.JSON(http.StatusOK, rd)
}
