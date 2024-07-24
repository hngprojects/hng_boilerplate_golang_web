package models

import "time"

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

type InvitationAcceptReq struct {
	InvitationLink string `json:"invitation_link" validate:"required"`
}
