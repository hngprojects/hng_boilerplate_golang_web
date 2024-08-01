package user

import (
	"github.com/gin-gonic/gin"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"net/http"
	"strconv"
)

// AssignRoleToUser godoc
// @Summary Assign role to user
// @Description Assign a new role to a user
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param role_id path int true "Role ID"
// @Success 200 {object} utility.Response
// @Failure 400,404 {object} utility.Response
// @Router /users/{user_id}/roles/{role_id} [put]
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
