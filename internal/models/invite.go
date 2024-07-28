package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Invitation struct {
	ID             string       `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	UserID         string       `gorm:"type:uuid;" json:"user_id"`
	OrganisationID string       `gorm:"type:uuid;" json:"organisation_id"`
	Organisation   Organisation `gorm:"foreignKey:OrganisationID"`
	Token          string       `gorm:"type:varchar(255);" json:"token"`
	CreatedAt      time.Time    `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	ExpiresAt      time.Time    `gorm:"column:expires_at; not null" json:"expires_at"`
	IsValid        bool         `gorm:"type:boolean;default:true" json:"is_valid"`
	Email          string       `gorm:"type:varchar(100);" json:"email"`
}

type InvitationRequest struct {
	Emails []string `json:"emails" validate:"required"`
	OrgID  string   `json:"org_id" validate:"required,uuid"`
}

type InvitationResponse struct {
	Email       string    `json:"email"`
	OrgID       string    `json:"org_id"`
	Status      string    `json:"status"`
	InviteToken string    `json:"invite_token"`
	Sent_At     time.Time `json:"sent_at"`
	Expires_At  time.Time `json:"expires_at"`
}

type InvitationCreateReq struct {
	OrganisationID string `json:"organisation_id" validate:"required,uuid"`
	Email          string `json:"email" validate:"required,email"`
}

func (i *Invitation) CreateInvitation(db *gorm.DB) error {
	//set the expiration time to 24 hours
	i.ExpiresAt = time.Now().Add(24 * time.Hour)

	err := postgresql.CreateOneRecord(db, &i)
	if err != nil {
		return err
	}
	return nil
}

func (i *Invitation) GetInvitationsByID(db *gorm.DB, user_id string) ([]Invitation, error) {
	//get all invitations with the user_id
	var invitations []Invitation

	err := postgresql.SelectAllFromDb(db.Preload("Organisation"), "", &invitations, "user_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

type InvitationAcceptReq struct {
	InvitationLink string `json:"invitation_link" validate:"required"`
}