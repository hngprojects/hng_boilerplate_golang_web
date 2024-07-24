package invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func CheckerValidator(c *gin.Context, base *storage.Database, inviteReq models.InvitationCreateReq, userId string) (error) {
	//check if organisation exists
	_, err := organisation.CheckOrgExists(inviteReq.OrganisationID, base.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Invalid Organisation ID", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return err
	}

	//check if user is an admin of the organisation
	isAdmin, err := CheckUserIsAdmin(base.Postgresql, userId, inviteReq.OrganisationID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to check if user is an admin", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return err
	}
	if !isAdmin {
		rd := utility.BuildErrorResponse(http.StatusForbidden, "error", "User is not an admin of the organisation", nil, nil)
		c.JSON(http.StatusForbidden, rd)
		return err
	}
	return nil
}
