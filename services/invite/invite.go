package invite

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
	"strings"
	"time"
)

// check if user is an admin
func CheckUserIsAdmin(db *gorm.DB, user_id string, org_id string) (bool, error) {
	//use the org_id to check for an existing organisation
	var org models.Organisation

	orgResp, err := org.GetOrgByID(db, org_id)
	if err != nil {
		return false, err
	}
	return orgResp.OwnerID == user_id, nil
}

// check emails limit
func CheckEmailsLimit(inviteReq models.InvitationRequest) bool {
	return len(inviteReq.Emails) > 5 // limit to 5 emails for testing
}

// check duplicate emails
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
func ExtractTokenFromInvitationLink(invitationLink string) string {
	splitLink := strings.Split(invitationLink, "/")
	return splitLink[len(splitLink)-1]
}

func GetInvitationDetails(token string, db *gorm.DB) (error, models.Invitation) {
	// Check if the invitation token exists in the database
	exists := postgresql.CheckExists(db, &models.Invitation{}, "token = ?", token)
	// If it does, return the invitation details
	if exists {
		var invitation models.Invitation
		postgresql.SelectOneFromDb(db, &invitation, "token = ?", token)
		return nil, invitation

	}
	return errors.New("Invalid invitation link format"), models.Invitation{}

}
func AcceptInvitationLink(token string, db *gorm.DB) (error, models.Invitation) {

	err, invitation := GetInvitationDetails(token, db)
	if err != nil {
		return err, models.Invitation{}
	}
	if invitation.ExpiresAt.Before(time.Now()) {
		return errors.New("Expired Invitation Link"), models.Invitation{}
	}
	if !invitation.IsValid {
		return errors.New("Invitation Link already used"), models.Invitation{}
	}
	if invitation.OrganisationID == "" {
		return errors.New("Organization not found"), models.Invitation{}
	}

	return nil, invitation
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
