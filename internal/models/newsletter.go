package models

import (
	"time"

	"github.com/google/uuid"
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
		n.ID = uuid.New().String()
	}
	return
}

func (NewsLetter) TableName() string {
	return "subscribers_table"
}
