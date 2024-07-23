package invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/invite"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func (base *Controller) CreateInvite(c *gin.Context) {
	var inviteReq models.InvitationCreateReq

	if err := c.ShouldBindJSON(&inviteReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	//validate request using default validator
	err := base.Validator.Struct(&inviteReq)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	//check if email format is correct
	_, valid := utility.EmailValid(inviteReq.Email)
	if !valid {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid email format", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	//check if organisation exists
	_, err = organisation.CheckOrgExists(inviteReq.OrganisationID, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Invalid Organisation ID", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	//generate token
	token, err := invite.GenerateInvitationToken()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to generate token", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	//save invitation
	err = invite.SaveInvitation(base.Db.Postgresql, token, inviteReq)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to save invitation", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	//generate invitation link
	invitationLink := invite.GenerateInvitationLink("http://localhost:8019", token)

	rd := utility.BuildSuccessResponse(http.StatusCreated, "success", "Invitation created successfully", invitationLink, nil)
	c.JSON(http.StatusCreated, rd)
}
