package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func (base *Controller) CompleteUserAuth(c *gin.Context) {

	userResp, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusServiceUnavailable, "error", err.Error(), err, nil)
		c.JSON(http.StatusServiceUnavailable, rd)
		return
	}

	userReq := models.CreateUserRequestModel{
		UserName:  userResp.Name,
		Email:     userResp.Email,
		FirstName: userResp.FirstName,
		LastName:  userResp.LastName,
	}

	respData, code, err := auth.CreateProviderUser(userReq, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, respData)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("user sign in successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "user sign in successfully", respData)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ProviderLogin(c *gin.Context) {

	var userResp goth.User

	userResp, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err == nil {

		userReq := models.CreateUserRequestModel{
			UserName:  userResp.Name,
			Email:     userResp.Email,
			FirstName: userResp.FirstName,
			LastName:  userResp.LastName,
		}

		respData, code, err := auth.CreateProviderUser(userReq, base.Db.Postgresql)
		if err != nil {
			rd := utility.BuildErrorResponse(code, "error", err.Error(), err, respData)
			c.JSON(code, rd)
			return
		}

		base.Logger.Info("user sign in successfully")

		rd := utility.BuildSuccessResponse(http.StatusOK, "user sign in successfully", respData)
		c.JSON(http.StatusOK, rd)
	} else {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}

}

func (base *Controller) ProviderLogout(c *gin.Context) {

	err := gothic.Logout(c.Writer, c.Request)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("user logged out successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "user logged out successfully", nil)
	c.JSON(http.StatusOK, rd)

}
