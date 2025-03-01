package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func (base *Controller) UpdateUserRegion(c *gin.Context) {
	var (
		userID = c.Param("user_id")
		req    = models.UserRegionTimezoneLanguage{}
	)

	requesterID, err := authHelper(c, *base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", err.Error(), nil, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	err = c.ShouldBind(&req)
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

	respData, code, err := base.UserService.UpdateARegion(req, userID, requesterID)
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

	requesterID, err := authHelper(c, *base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", err.Error(), nil, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	respData, code, err := base.UserService.GetUserRegion(userID, requesterID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "User region retrieved successfully", respData)
	c.JSON(http.StatusOK, rd)

}
