package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

// @Summary Update user region
// @Description Update a user's region, timezone, and language
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body models.UserRegionTimezoneLanguage true "User region update request"
// @Success 200 {object} utility.Response "User info updated successfully"
// @Failure 400 {object} utility.Response "Failed to parse request body"
// @Failure 422 {object} utility.Response "Validation failed"
// @Router /users/{user_id}/region [put]
func (base *Controller) UpdateUserRegion(c *gin.Context) {
	var (
		userID = c.Param("user_id")
		req    = models.UserRegionTimezoneLanguage{}
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

	respData, code, err := service.UpdateARegion(req, userID, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("user info updated successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "User info updated successfully", respData)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetUserRegion(c *gin.Context) {
	var (
		userID = c.Param("user_id")
	)

	respData, code, err := service.GetUserRegion(userID, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "User region retrieved successfully", respData)
	c.JSON(http.StatusOK, rd)

}
