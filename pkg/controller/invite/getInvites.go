package invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/invite"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

// GetInvites godoc
// @Summary Get all invitations
// @Description Retrieve all invitations for a super admin user
// @Tags invite
// @Produce json
// @Success 200 {object} utility.Response "Invitations Successfully retrieved"
// @Failure 400 {object} utility.Response "unable to get user claims"
// @Failure 403 {object} utility.Response "User is not a super admin"
// @Failure 500 {object} utility.Response "Internal server error"
// @Router /invites [get]
func (base *Controller) GetInvites(c *gin.Context) {

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	user, code, err := user.GetUser(userId, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	isSuperAdmin := user.CheckUserIsAdmin(base.Db.Postgresql)
	if !isSuperAdmin {
		rd := utility.BuildErrorResponse(http.StatusForbidden, "error", "User is not a super admin", nil, nil)
		c.JSON(http.StatusForbidden, rd)
		return
	}

	invitations, err := invite.GetInvitations(user, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("invitations fetched successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "Invitations Successfully retrieved", invitations)

	c.JSON(http.StatusOK, rd)
}
