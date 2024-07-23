package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type DefaultIdentity struct {
	User       int
	SuperAdmin int
}

var RoleIdentity = DefaultIdentity{
	User:       1,
	SuperAdmin: 2,
}

type Role struct {
	ID          int            `gorm:"primaryKey;type:int" json:"id"`
	Name        string         `gorm:"unique;not null;type:varchar(20)" json:"name" validate:"required"`
	Description string         `gorm:"unique;not null" json:"description" validate:"required"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Role) CreateRole(db *gorm.DB) error {

	err := postgresql.CreateOneRecord(db, &r)

	if err != nil {
		return err
	}

	return nil
}
