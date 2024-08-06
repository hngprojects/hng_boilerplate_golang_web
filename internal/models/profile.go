package models

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

type Profile struct {
	ID             string         `gorm:"type:uuid;primary_key" json:"profile_id"`
	FirstName      string         `gorm:"column:first_name; type:text; not null" json:"first_name"`
	LastName       string         `gorm:"column:last_name; type:text;not null" json:"last_name"`
	Phone          string         `gorm:"type:varchar(255)" json:"phone"`
	AvatarURL      string         `gorm:"type:varchar(255)" json:"avatar_url"`
	Userid         string         `gorm:"type:uuid;" json:"user_id"`
	SecondaryEmail string         `gorm:"type:string;" json:"secondary_email"`
	CreatedAt      time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type UpdateProfileRequest struct {
	SecondaryEmail string `json:"secondary_email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Phone          string `json:"phone"`
}

func (p *Profile) UpdateProfile(db *gorm.DB, req UpdateProfileRequest, profId string) error {

	result, err := postgresql.UpdateFields(db, &p, req, profId)
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return errors.New("failed to update organisation")
	}

	return nil
}
