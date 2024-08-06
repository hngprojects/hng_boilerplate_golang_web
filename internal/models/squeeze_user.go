package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type SqueezeUser struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	FirstName string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string    `gorm:"type:varchar(100);not null" json:"last_name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

type SqueezeUserReq struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

func (s *SqueezeUser) Create(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &s)

	if err != nil {
		return err
	}

	return nil
}
