package invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func InvitationLinkGenerator(c *gin.Context, base *storage.Database, inviteReq models.InvitationCreateReq, userId string) (string, error) {
	token, err := GenerateInvitationToken()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to generate token", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return "", err
	}

	err = SaveInvitation(base.Postgresql, userId, token, inviteReq)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to save invitation", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return "", err
	}

	//generate invitation link
	invitationLink := GenerateInvitationLink("http://localhost:8019/api/v1", token)

	return invitationLink, nil
}