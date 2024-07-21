package invite

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

// check emails limit
func CheckEmailsLimit(inviteReq models.InvitationRequest) bool {
	if len(inviteReq.Emails) > 50 {
		return true
	}
	return false
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
