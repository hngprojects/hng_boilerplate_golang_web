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
		return errors.New("User is not an admin of the organisation")
	}
	return nil
}

func CheckerPostInvite(c *gin.Context,base *storage.Database,inviteReq models.InvitationRequest) (models.Organisation, error) {
	var org models.Organisation

	//check if user is jwt authenticated
	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return org, errors.New("unable to get user claims")
	}
	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	//check if duplicate emails exist
	if CheckDuplicateEmails(inviteReq) {
		rd := utility.BuildErrorResponse(http.StatusConflict, "error", "Duplicate emails found", nil, nil)
		c.JSON(http.StatusConflict, rd)
		return org, errors.New("Duplicate emails found")
	}

	// check emails limit
	if CheckEmailsLimit(inviteReq) {
		rd := utility.BuildErrorResponse(http.StatusRequestEntityTooLarge, "error", "Payload too large; email limit exceeded", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return org, errors.New("Payload too large; email limit exceeded")
	}

	// Validate org_id
	orgId, err := uuid.Parse(inviteReq.OrgID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Invalid org_id format", err, nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return org, err
	}

	// Check if org_id exists and return organization
	orgResp, err := organisation.CheckOrgExists(orgId.String(), base.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "organisation not found", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return org, err
	}

	// Check if user is a member of the organization
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
	// Store invitations
	invitations := []map[string]interface{}{}

	// Loop through emails and create invitation
	for _, email := range inviteReq.Emails {
		//check if email is an empty string
		if email == "" {
			invitations = append(invitations,
				map[string]interface{}{
					"error":   "invalid request",
					"message": "email address cannot be empty",
				},
			)
			// Log error and skip user
			msg := fmt.Sprintf("missing email field %s", email)
			rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", msg, nil, nil)
			c.JSON(http.StatusUnprocessableEntity, rd)
			continue
		}

		// Check if email is valid
		if _, valid := utility.EmailValid(email); !valid {
			invitations = append(invitations,
				map[string]interface{}{
					"error":   "invalid request",
					"message": fmt.Sprintf("email address %s not valid", email),
				},
			)
			// Log error and skip user
			msg := fmt.Sprintf("invalid email format for %s", email)
			rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", msg, nil, nil)
			c.JSON(http.StatusUnprocessableEntity, rd)
			continue
		}

		// Check if user with email exists and get user
		user, err := user.GetUserByEmail(email, base.Postgresql)
		if err != nil {

			// Log error and skip user
			logger.Error("error getting user by email", err)
			continue
		}

		// Create invitation
		invitation := models.Invitation{
			ID:             utility.GenerateUUID(),
			UserID:         user.ID,
			OrganisationID: org.ID,
			Email:          email,
			CreatedAt:      time.Now(),
			ExpiresAt:      time.Now().Add(time.Hour * 24),
			IsValid:        true,
		}

		// Store invitation in db
		err = invitation.CreateInvitation(base.Postgresql)
		if err != nil {
			// Log error and skip user
			logger.Error("error creating invitation", err)
			continue
		}

		// Append invitation to invitations
		invitations = append(
			invitations,
			map[string]interface{}{
				"email":        email,
				"organization": org.Name,
				"expires_at":   invitation.ExpiresAt.Format(time.RFC3339),
				// "token":        token,
			},
		)

		// // Send email
		// subject := "You're Invited to Join " + org.Name
		// senderEmail := "micahshallom@gmail.com"
		// to := []string{"hamzasaidu34@gmail.com"}
		// domain := config.Config.Mail.Domain
		// apikey := config.Config.Mail.APIKey

		// inviteErr := send_invites.MockSendInvite(domain, apikey, senderEmail, to, subject)
		// if inviteErr != nil {
		// 	// Log error and skip user
		// 	logger.Error("error sending invitation", err)
		// 	continue
		// }

		logger.Info("Invitations posted successfully")
		rd := utility.BuildSuccessResponse(http.StatusCreated, "Invitation(s) sent successfully", invitations)

		c.JSON(http.StatusCreated, rd)
	}
}
