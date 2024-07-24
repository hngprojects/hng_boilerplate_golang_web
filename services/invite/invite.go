package invite

import (
	"errors"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
	"strings"
	"time"
)

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
