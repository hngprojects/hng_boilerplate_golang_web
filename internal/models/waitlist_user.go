package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

var ErrWaitlistUserExist = errors.New("waitlist user exists")

type WaitlistUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWaitlistUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func (w *WaitlistUser) CreateWaitlistUser(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, w)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrWaitlistUserExist
		}
	}

	return err
}

func (w *WaitlistUser) GetWaitlistUserByEmail(db *gorm.DB) (int, error) {
	err, nerr := postgresql.SelectOneFromDb(db, &w, "email = ?", w.Email)
	if nerr != nil {
		return http.StatusBadRequest, nerr 
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
