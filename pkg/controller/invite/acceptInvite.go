package invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/invite"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func (base *Controller) PostAcceptInvite(c *gin.Context) {
	// accept invite logic here
	var inviteReq models.InvitationAcceptReq
	claims, exists := c.Get("userClaims")
	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)
	err := c.ShouldBind(&inviteReq)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	err = base.Validator.Struct(&inviteReq)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}
	invitationToken := invite.ExtractTokenFromInvitationLink(inviteReq.InvitationLink)
	invitation, msg, err := invite.AcceptInvitationLink(userId, invitationToken, base.Db.Postgresql)
	if err != nil {
		base.Logger.Error("Failed to accept invitation link", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", msg, err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	// add user to organisation
	///check if user from the claims is a member of the organisation
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	err = invite.AddUserToOrganisation(base.Db.Postgresql, invitation.OrganisationID, userId)
	if err != nil {
		base.Logger.Error("Failed to add user to organisation", err)
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "A server error occurred", nil, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "Invitation accepted successfully", nil)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetAcceptInvite(c *gin.Context) {
	// get accept invite logic here
	invitationToken := c.Param("t")
	claims, exists := c.Get("userClaims")
	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)
	invitation, msg, err := invite.AcceptInvitationLink(userId, invitationToken, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(
			http.StatusBadRequest,
			"error",
			msg,
			err,
			nil,
		)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	// add user to organisation
	///check if user from the claims is a member of the organisation
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	err = invite.AddUserToOrganisation(base.Db.Postgresql, invitation.OrganisationID, userId)
	if err != nil {
		base.Logger.Error("Failed to add user to organisation", err)
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "A server error occurred", nil, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "Invitation accepted successfully", nil)
	c.JSON(http.StatusOK, rd)
}
