package user

import (
	"github.com/gin-gonic/gin"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"net/http"
	"strconv"
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

	rd := utility.BuildSuccessResponse(http.StatusOK, "Role updated successfully", userData)
	c.JSON(http.StatusOK, rd)
}
