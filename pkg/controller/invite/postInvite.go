package invite

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	// "github.com/hngprojects/hng_boilerplate_golang_web/internal/config"

	// "github.com/hngprojects/hng_boilerplate_golang_web/external/mocks/send_invites"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/invite"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) PostInvite(c *gin.Context) {
	var inviteReq models.InvitationRequest

	if err := c.ShouldBindJSON(&inviteReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	//check if duplicate emails exist
	if invite.CheckDuplicateEmails(inviteReq) {
		rd := utility.BuildErrorResponse(http.StatusConflict, "error", "Duplicate emails found", nil, nil)
		c.JSON(http.StatusConflict, rd)
		return
	}

	// check emails limit
	if invite.CheckEmailsLimit(inviteReq) {
		rd := utility.BuildErrorResponse(http.StatusRequestEntityTooLarge, "error", "Payload too large; email limit exceeded", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err := base.Validator.Struct(&inviteReq)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// Validate org_id
	orgId, err := uuid.Parse(inviteReq.OrgID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Invalid org_id format", err, nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	// Check if org_id exists and return organization
	org, err := organisation.CheckOrgExists(orgId.String(), base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "organisation not found", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	///check if user from the claims is a member of the organisation
	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	// Check if user is a member of the organization
	isMember, err := organisation.CheckUserIsMemberOfOrg(userId, org.ID, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "User not a member of the organization", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	if !isMember {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "User not a member of the organization", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

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
		user, err := user.GetUserByEmail(email, base.Db.Postgresql)
		if err != nil {

			// Log error and skip user
			base.Logger.Error("error getting user by email", err)
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
		err = invitation.CreateInvitation(base.Db.Postgresql)
		if err != nil {
			// Log error and skip user
			base.Logger.Error("error creating invitation", err)
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
		// 	base.Logger.Error("error sending invitation", err)
		// 	continue
		// }

		base.Logger.Info("Invitations posted successfully")
		rd := utility.BuildSuccessResponse(http.StatusCreated, "Invitation(s) sent successfully", invitations)

		c.JSON(http.StatusCreated, rd)
	}
}
