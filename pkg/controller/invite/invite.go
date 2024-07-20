package invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

type InvitationRequest struct {
	Emails []string `json:"emails" validate:"required"`
	OrgID  string   `json:"org_id" binding:"uuid"`
}



func (base *Controller) PostInvite(c *gin.Context) {

	var inviteReq InvitationRequest

	if err := c.ShouldBindJSON(&inviteReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err := base.Validator.Struct(&inviteReq)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest,"error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// if err != nil {
	// 	rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), err, nil)
	// 	c.JSON(http.StatusInternalServerError, rd)
	// 	return
	// }

	base.Logger.Info("invite posted successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "", "invite url")

	c.JSON(http.StatusOK, rd)
}
