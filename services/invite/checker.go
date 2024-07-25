package invite

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func CheckerValidator(c *gin.Context, base *storage.Database, inviteReq models.InvitationCreateReq, userId string) error {
	_, err := organisation.CheckOrgExists(inviteReq.OrganisationID, base.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Invalid Organisation ID", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return err
	}

	isAdmin, err := CheckUserIsAdmin(base.Postgresql, userId, inviteReq.OrganisationID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to check if user is an admin", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return err
	}
	if !isAdmin {
		rd := utility.BuildErrorResponse(http.StatusForbidden, "error", "User is not an admin of the organisation", nil, nil)
		c.JSON(http.StatusForbidden, rd)
		return errors.New("User is not an admin of the organisation")
	}
	return nil
}

func CheckerPostInvite(c *gin.Context,base *storage.Database,inviteReq models.InvitationRequest) (models.Organisation, error) {
	var org models.Organisation

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return org, errors.New("unable to get user claims")
	}
	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	if CheckDuplicateEmails(inviteReq) {
		rd := utility.BuildErrorResponse(http.StatusConflict, "error", "Duplicate emails found", nil, nil)
		c.JSON(http.StatusConflict, rd)
		return org, errors.New("Duplicate emails found")
	}

	if CheckEmailsLimit(inviteReq) {
		rd := utility.BuildErrorResponse(http.StatusRequestEntityTooLarge, "error", "Payload too large; email limit exceeded", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return org, errors.New("Payload too large; email limit exceeded")
	}

	orgId, err := uuid.Parse(inviteReq.OrgID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Invalid org_id format", err, nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return org, err
	}

	orgResp, err := organisation.CheckOrgExists(orgId.String(), base.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "organisation not found", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return org, err
	}

	isMember, err := organisation.CheckUserIsMemberOfOrg(userId, orgResp.ID, base.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "User not a member of the organization", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return org, err
	}

	if !isMember {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "User not a member of the organization", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return org, errors.New("User not a member of the organization")
	}

	return org, nil
}

func IteratorPostInvite(c *gin.Context, inviteReq models.InvitationRequest, base *storage.Database,logger *utility.Logger, org models.Organisation) {
	invitations := []map[string]interface{}{}

	for _, email := range inviteReq.Emails {
		if email == "" {
			invitations = append(invitations,
				map[string]interface{}{
					"error":   "invalid request",
					"message": "email address cannot be empty",
				},
			)
			msg := fmt.Sprintf("missing email field %s", email)
			rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", msg, nil, nil)
			c.JSON(http.StatusUnprocessableEntity, rd)
			continue
		}

		if _, valid := utility.EmailValid(email); !valid {
			invitations = append(invitations,
				map[string]interface{}{
					"error":   "invalid request",
					"message": fmt.Sprintf("email address %s not valid", email),
				},
			)
			msg := fmt.Sprintf("invalid email format for %s", email)
			rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", msg, nil, nil)
			c.JSON(http.StatusUnprocessableEntity, rd)
			continue
		}

		user, err := user.GetUserByEmail(email, base.Postgresql)
		if err != nil {

			logger.Error("error getting user by email", err)
			continue
		}

		invitation := models.Invitation{
			ID:             utility.GenerateUUID(),
			UserID:         user.ID,
			OrganisationID: org.ID,
			Email:          email,
			CreatedAt:      time.Now(),
			ExpiresAt:      time.Now().Add(time.Hour * 24),
			IsValid:        true,
		}

		err = invitation.CreateInvitation(base.Postgresql)
		if err != nil {
			logger.Error("error creating invitation", err)
			continue
		}

		invitations = append(
			invitations,
			map[string]interface{}{
				"email":        email,
				"organization": org.Name,
				"expires_at":   invitation.ExpiresAt.Format(time.RFC3339),
			},
		)

		logger.Info("Invitations posted successfully")
		rd := utility.BuildSuccessResponse(http.StatusCreated, "Invitation(s) sent successfully", invitations)

		c.JSON(http.StatusCreated, rd)
	}
}
