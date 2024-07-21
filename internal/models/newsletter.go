package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type NewsLetter struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email" validate:"required,email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (n *NewsLetter) BeforeCreate(tx *gorm.DB) (err error) {

	if n.ID == "" {
		n.ID = utility.GenerateUUID()
	}
	return
}

func (c *NewsLetter) CreateNewsLetter(db *gorm.DB) error {

	err := postgresql.CreateOneRecord(db, &c)

	if err != nil {
		return err
	}

	return nil
}
