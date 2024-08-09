package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Testimonial struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

type TestimonialReq struct {
	Name    string `json:"name" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (t *Testimonial) Create(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)

	if err != nil {
		return err
	}

	return nil
}
