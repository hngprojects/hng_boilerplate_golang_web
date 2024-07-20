package models

import (
	"errors"
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

var ErrWaitlistUserExist = errors.New("waitlist user exists")

type WaitlistUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWaitlistUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
}

func (w *WaitlistUser) CreateWaitlistUser(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, w)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return ErrWaitlistUserExist
		}
	}

	return err
}
