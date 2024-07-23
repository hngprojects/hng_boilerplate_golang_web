package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Invitation struct {
	ID             string       `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	UserID         string       `gorm:"type:uuid;"`
	OrganisationID string       `gorm:"type:uuid;"`
	Organisation   Organisation `gorm:"foreignKey:OrganisationID"`
	Token          string       `gorm:"type:varchar(255);"`
	CreatedAt      time.Time
	ExpiresAt      time.Time
	IsValid        bool
	Email          string `gorm:"type:varchar(100);"`
}

type InvitationRequest struct {
	Emails []string `json:"emails" validate:"required"`
	OrgID  string   `json:"org_id" validate:"required,uuid"`
}

type InvitationCreateReq struct {
	OrganisationID string `json:"organisation_id" validate:"required,uuid"`
	Email          string `json:"email" validate:"required,email"`
}

func (i *Invitation) CreateInvitation(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &i)
	if err != nil {
		return err
	}
	return nil
}
