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

type mapper map[string]interface{}

func CheckerValidator(base *storage.Database, inviteReq models.InvitationCreateReq, userId string, logger *utility.Logger) (models.Organisation, int, string, error) {
	//check if organisation exists
	orgResp, err := organisation.CheckOrgExists(inviteReq.OrganisationID, base.Postgresql)
	if err != nil {
		return orgResp, http.StatusNotFound, "Invalid Organisation ID", err
	}

	isAdmin, err := CheckUserIsAdmin(base.Postgresql, userId, inviteReq.OrganisationID)
	if err != nil {
		return orgResp, http.StatusInternalServerError, "Internal server error", err
	}
	if !isAdmin {
		return orgResp, http.StatusUnauthorized, "User is not an admin of the organisation", errors.New("User is not an admin of the organisation")
	}
	return orgResp, http.StatusOK, "", nil
}

func CheckerPostInvite(base *storage.Database, inviteReq models.InvitationRequest, userId string) (models.Organisation, int, string, error) {
	var org models.Organisation

	// Check if duplicate emails exist
	if CheckDuplicateEmails(inviteReq) {
		return org, http.StatusConflict, "Duplicate emails found", errors.New("Duplicate emails found")
	}

	// Check emails limit
	if CheckEmailsLimit(inviteReq) {
		return org, http.StatusRequestEntityTooLarge, "Payload too large; email limit exceeded", errors.New("Payload too large; email limit exceeded")
	}

	// Validate org_id
	orgId, err := uuid.Parse(inviteReq.OrgID)
	if err != nil {
		return org, http.StatusUnprocessableEntity, "Invalid org_id format", err
	}

	// Check if org_id exists and return organization
	orgResp, err := organisation.CheckOrgExists(orgId.String(), base.Postgresql)
	if err != nil {
		return org, http.StatusNotFound, "organisation not found", err
	}

	orgResp, statusCode, msg, err := CheckerValidator(
		base,
		models.InvitationCreateReq{
			OrganisationID: orgResp.ID,
			Email:          "",
		},
		userId,
		utility.NewLogger(),
	)
	if err != nil {
		return org, statusCode, msg, err
	}

	// Check if user is a member of the organization
	isMember, err := organisation.CheckUserIsMemberOfOrg(userId, orgResp.ID, base.Postgresql)
	if err != nil {
		return org, http.StatusNotFound, "User not a member of the organization", err
	}
	if !isMember {
		return org, http.StatusNotFound, "User not a member of the organization", errors.New("User not a member of the organization")
	}

	return orgResp, http.StatusOK, "", nil
}

func IteratorPostInvite(c *gin.Context, inviteReq models.InvitationRequest, base *storage.Database, logger *utility.Logger, org models.Organisation) (int, string, []mapper) {
	var invitations []mapper
	var inviteErrors []mapper

	if len(inviteReq.Emails) == 0 {
		return http.StatusBadRequest, "No emails provided", nil
	}

	// Loop through emails and create invitation
	for _, email := range inviteReq.Emails {
		if email == "" {
			inviteErrors = append(
				inviteErrors,
				map[string]interface{}{
					"error": "Email address cannot be empty"},
			)
			continue
		}

		if _, valid := utility.EmailValid(email); !valid {
			inviteErrors = append(
				inviteErrors,
				map[string]interface{}{
					"error": fmt.Sprintf("Invalid email address: %s", email)},
			)
			continue
		}

		user, err := user.GetUserByEmail(email, base.Postgresql)
		if err != nil {
			inviteErrors = append(
				inviteErrors,
				map[string]interface{}{
					"error": fmt.Sprintf("error getting user by email: %s", email),
				},
			)
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
			inviteErrors = append(
				inviteErrors,
				map[string]interface{}{
					"error": fmt.Sprintf("error creating invitation for email: %s", email),
				},
			)

			continue
		}

		// Send email
		err = SendEmail(email, org.Name, invitation.ExpiresAt.Format(time.RFC3339))
		if err != nil {
			inviteErrors = append(
				inviteErrors,
				map[string]interface{}{
					"error": fmt.Sprintf("error sending email to: %s", email),
				},
			)
			continue
		}

		invitations = append(invitations, map[string]interface{}{
			"email":        email,
			"organization": org.Name,
			"expires_at":   invitation.ExpiresAt.Format(time.RFC3339),
		})
	}

	if len(inviteErrors) > 0 {
		rd := utility.BuildSuccessResponse(
			http.StatusOK,
			fmt.Sprintf("%d invitations sent successfully", len(invitations)),
			invitations,
		)
		c.JSON(http.StatusOK, rd)

		return http.StatusBadRequest, fmt.Sprintf("%d invitations failed", len(inviteErrors)), inviteErrors
	}

	return http.StatusCreated, "Invitation(s) sent successfully", invitations
}

// write a dummy sending email function
func SendEmail(email string, orgName string, expiresAt string) error {
	return nil
}

func CheckerPostInvite(c *gin.Context, base *storage.Database, inviteReq models.InvitationRequest) (models.Organisation, error) {
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

func IteratorPostInvite(c *gin.Context, inviteReq models.InvitationRequest, base *storage.Database, logger *utility.Logger, org models.Organisation) {
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
