package invite

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/mocks/send_invites"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/seed"
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
	OrgID  string   `json:"org_id" validate:"required,uuid"`
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
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// Validate org_id
	orgId, err := uuid.Parse(inviteReq.OrgID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid org_id parsed", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// Check if org_id exists and return organization
	org, err := seed.CheckOrgExists(orgId.String(), base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Organisation not found", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}


	// Store invitations
	invitations := []map[string]interface{}{}

	// Loop through emails and create invitation
	for _, email := range inviteReq.Emails {
		// Check if email is valid
		if !utility.EmailValid(email) {
			invitations = append(invitations, 
				map[string]interface{}{
					"error": "invalid request",
					"message": fmt.Sprintf("email address %s not valid", email),
				},
			)
			// Log error and skip user
			base.Logger.Error(email, "is not a valid email")
			continue
		}

		// Check if user with email exists and get user
		user, err := seed.GetUserByEmail(email, base.Db.Postgresql)
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
		err = invitation.CreateInvitation(base.Db.Postgresql, invitation)
		if err != nil {
			// Log error and skip user
			base.Logger.Error("error creating invitation", err)
			continue
		}

		//send invitation
		token, err := send_invites.MockSendInvite(invitation.Email, org.Name, invitation.ExpiresAt)
		if err != nil {
			// Log error and skip user
			base.Logger.Error("error sending invitation", err)
			continue
		}

		// Append invitation to invitations
		invitations = append(
			invitations,
			map[string]interface{}{
				"email":        email,
				"organization": org.Name,
				"expires_at":   invitation.ExpiresAt.Format(time.RFC3339),
				"token":        token,
			},
		)
	}

	base.Logger.Info("Invitations posted successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Invitation(s) sent successfully", invitations)

	c.JSON(http.StatusCreated, rd)
}
