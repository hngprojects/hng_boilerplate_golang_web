package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SqueezeUser struct {
	ID             string         `gorm:"type:uuid;primary_key" json:"id"`
	Email          string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	FirstName      string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName       string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Phone          string         `gorm:"type:varchar(20);not null" json:"phone"`
	Location       string         `gorm:"type:varchar(255);not null" json:"location"`
	JobTitle       string         `gorm:"type:varchar(100);not null" json:"job_title"`
	Company        string         `gorm:"type:varchar(255);not null" json:"company"`
	Interests      pq.StringArray `gorm:"type:text[]" json:"interests"`
	ReferralSource string         `gorm:"type:varchar(100);not null" json:"referral_source"`
	CreatedAt      time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

type SqueezeUserReq struct {
	Email          string   `json:"email" validate:"required,email"`
	FirstName      string   `json:"first_name" validate:"required"`
	LastName       string   `json:"last_name" validate:"required"`
	Phone          string   `json:"phone" validate:"required"`
	Location       string   `json:"location" validate:"required"`
	JobTitle       string   `json:"job_title" validate:"required"`
	Company        string   `json:"company" validate:"required"`
	Interests      []string `json:"interests" validate:"required"`
	ReferralSource string   `json:"referral_source" validate:"required"`
}

func (s *SqueezeUser) Create(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &s)

	if err != nil {
		return err
	}

	return nil
}
