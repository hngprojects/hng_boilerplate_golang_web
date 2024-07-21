package service

import (
	"errors"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

var ErrEmailAlreadySubscribed = errors.New("email already subscribed")

func NewsLetterSubscribe(newsletter *models.NewsLetter, db *gorm.DB) error {

	if postgresql.CheckExists(db, newsletter, "email = ?", newsletter.Email) {
		return ErrEmailAlreadySubscribed
	}

	if err := postgresql.CreateOneRecord(db, newsletter); err != nil {
		return err
	}

	return nil
}
