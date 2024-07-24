package models

import (
	"time"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)
type Invitation struct {
	ID             string         `gorm:"primaryKey;type:uuid" json:"id"`
	Email          string         `gorm:"unique;not null" json:"email" validate:"required,email"`
	OrganizationID string         `gorm:"type:uuid;not null" json:"organization_id"`
	Organisation   Organisation   `gorm:"foreignKey:OrganizationID" json:"organization"`
	IsValid        bool           `gorm:"not null;default:true" json:"is_valid"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to set a UUID before creating a new Invitation
func (i *Invitation) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == "" {
		i.ID = utility.GenerateUUID()
	}
	return
}

func (i *Invitation) CreateInvitation(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, i)
	if err != nil {
		return err
	}
	return nil
}