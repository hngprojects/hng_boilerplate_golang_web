package invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/invite"

	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func (base *Controller) CreateInvite(c *gin.Context) {
	var inviteReq models.InvitationCreateReq

	if err := c.ShouldBindJSON(&inviteReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	err := base.Validator.Struct(&inviteReq)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	_, valid := utility.EmailValid(inviteReq.Email)
	if !valid {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid email format", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// call checker validator to check if user is an admin of the organisation and if organisation exists
	_, statusCode, msg, err := invite.CheckerValidator(base.Db, inviteReq, userId, base.Logger)
	if err != nil {
		rd := utility.BuildErrorResponse(statusCode, "error", msg, err, nil)
		c.JSON(statusCode, rd)
		return
	}

	// generate token, save to db and return invitation link
	inviteLink, err := invite.InvitationLinkGenerator(c, base.Db, inviteReq, userId)
	if err != nil {
		return
	}

	mapData := map[string]string{"invitation_link": inviteLink}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "Invitation created successfully", mapData)
	c.JSON(http.StatusCreated, rd)
}
