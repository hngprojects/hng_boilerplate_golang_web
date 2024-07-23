package invite

import (
	"crypto/rand"
	"encoding/hex"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

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

func SaveInvitation(db *gorm.DB, token string, req models.InvitationCreateReq) error {
	var invitation models.Invitation

	invitation.Token = token
	invitation.Email = req.Email
	invitation.OrganisationID = req.OrganisationID
	invitation.IsValid = true

	err := invitation.CreateInvitation(db)
	if err != nil {
		return err
	}
	return nil
}
