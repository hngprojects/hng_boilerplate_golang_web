package invite

import (
	"crypto/rand"
	"encoding/hex"

	"strings"

	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
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
