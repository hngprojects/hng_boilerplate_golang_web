package invite

import (
	"crypto/rand"
	"encoding/hex"

	"strings"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
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
