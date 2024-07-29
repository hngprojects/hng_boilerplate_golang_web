package invite

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"strings"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func CheckUserIsAdmin(db *gorm.DB, user_id string, org_id string) (bool, error) {
	var org models.Organisation

	orgResp, err := org.GetOrgByID(db, org_id)
	if err != nil {
		return false, err
	}
	return orgResp.OwnerID == user_id, nil
}

func CheckEmailsLimit(inviteReq models.InvitationRequest) bool {
	return len(inviteReq.Emails) > 5 // limit to 5 emails for testing
}

func CheckDuplicateEmails(inviteReq models.InvitationRequest) bool {
	emailsMap := make(map[string]bool)
	for _, email := range inviteReq.Emails {
		if _, ok := emailsMap[email]; ok {
			return true
		}
		emailsMap[email] = true
	}
	return false
}

func GenerateInvitationToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateInvitationLink(baseurl, token string) string {
	return baseurl + "/invite/accept/" + token
}

func SaveInvitation(db *gorm.DB, user_id string, token string, req models.InvitationCreateReq) error {
	var (
		email = strings.ToLower(req.Email)
	)

	invitation := models.Invitation{
		ID:             utility.GenerateUUID(),
		UserID:         user_id,
		Token:          token,
		Email:          email,
		OrganisationID: req.OrganisationID,
		IsValid:        true,
	}

	err := invitation.CreateInvitation(db)
	if err != nil {
		return err
	}
	return nil
}

func GetInvitations(user models.User, db *gorm.DB) ([]models.InvitationResponse, error) {
	var invitation models.Invitation
	var invResp []models.InvitationResponse

	invitations, err := invitation.GetInvitationsByID(db, user.ID)
	if err != nil {
		return invResp, err
	}

	for _, inv := range invitations {
		var status string
		switch inv.IsValid {
		case true:
			status = "active"
		default:
			status = "expired"
		}

		invResp = append(invResp, models.InvitationResponse{
			Email:       inv.Email,
			OrgID:       inv.OrganisationID,
			Status:      status,
			InviteToken: inv.Token,
			Sent_At:     inv.CreatedAt,
			Expires_At:  inv.ExpiresAt,
		})
	}
	return invResp, nil
}

func ExtractTokenFromInvitationLink(invitationLink string) string {
	splitLink := strings.Split(invitationLink, "/")
	return splitLink[len(splitLink)-1]
}

func GetInvitationDetails(token string, db *gorm.DB) (models.Invitation, error) {
	var invitation models.Invitation
	// Check if the invitation token exists in the database
	exists := postgresql.CheckExists(db, &invitation, "token = ?", token)
	// If it does, return the invitation details
	if exists {
		postgresql.SelectOneFromDb(db, &invitation, "token = ?", token)
		return invitation, nil
	}
	return invitation, errors.New("Invalid invitation link format")
}

func AcceptInvitationLink(user_id string, token string, db *gorm.DB) (models.Invitation, string, error) {
	var invitation models.Invitation

	invitation, err := GetInvitationDetails(token, db)
	if err != nil {
		return invitation,"Error getting invitation details", err
	}
	if invitation.ExpiresAt.Before(time.Now()) {
		return invitation, "Invitation link expired", errors.New("Invitation link expired")
	}
	if !invitation.IsValid {
		return invitation, "Invitation link is invalid", errors.New("Invitation link is invalid")
	}
	if invitation.OrganisationID == "" {
		return invitation, "Invalid organisation ID", errors.New("Invalid organisation ID")
	}

	//query the user, get the email and check if the email of the user is the same as the email in the invitation
	var user models.User
	postgresql.SelectOneFromDb(db, &user, "id = ?", user_id)
	if user.Email != invitation.Email {
		return invitation, "Invalid invitation link", errors.New("Invalid invitation link")
	}

	// Set the invitation to invalid and save it to the database
	invitation.IsValid = false
	_, err = postgresql.SaveAllFields(db, &invitation)
	if err != nil {
		return invitation, "Error saving invitation", err
	}

	return invitation, "Invitation link accepted successfully", nil
}


func AddUserToOrganisation(db *gorm.DB, orgID string, userId string) error {
	var user models.User

	user, err := user.GetUserByID(db, userId)
	if err != nil {
		return err
	}
	var org models.Organisation
	org, err = org.GetOrgByID(db, orgID)
	if err != nil {
		return err
	}
	err = user.AddUserToOrganisation(db, &user, []interface{}{&org})
	if err != nil {
		return err
	}
	return nil
}
