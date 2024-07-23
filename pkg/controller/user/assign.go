package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/user"
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
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "role successfully replaced", userData)
	c.JSON(http.StatusOK, rd)
}
