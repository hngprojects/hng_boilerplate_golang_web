package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func GetNewsletters(c *gin.Context, db *gorm.DB) ([]models.NewsLetter, *postgresql.PaginationResponse, int, error) {

	var newsletter models.NewsLetter

	newsLetters, paginationResponse, err := newsletter.FetchAllNewsLetter(db, c)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return newsLetters, nil, http.StatusNoContent, nil
		}
		return newsLetters, nil, http.StatusBadRequest, err

	}

	return newsLetters, &paginationResponse, http.StatusOK, nil

}

func GetDeletedNewsletters(c *gin.Context, db *gorm.DB) ([]models.NewsLetter, *postgresql.PaginationResponse, int, error) {

	var newsletter models.NewsLetter

	delNewsLetters, paginationResponse, err := newsletter.FetchAllDeletedNewsLetter(db, c)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return delNewsLetters, nil, http.StatusNoContent, nil
		}
		return delNewsLetters, nil, http.StatusBadRequest, err

	}

	return delNewsLetters, &paginationResponse, http.StatusOK, nil

}

func NewsLetterSubscribe(newsletter *models.NewsLetter, db *gorm.DB) error {

	if postgresql.CheckExists(db, newsletter, "email = ?", newsletter.Email) {
		return models.ErrEmailAlreadySubscribed
	}

	newsletter.Email = strings.ToLower(newsletter.Email)

	if err := newsletter.CreateNewsLetter(db); err != nil {
		return err
	}

	return nil
}

func DeleteNewsLetter(ID string, db *gorm.DB, c *gin.Context) (int, error) {
	var (
		newsLetter models.NewsLetter
	)

	newsLetter, err := newsLetter.GetNewsLetterById(db, ID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err := newsLetter.DeleteNewsLetter(db); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func RestoreDeletedNewsLetter(ID string, db *gorm.DB, c *gin.Context) (int, error) {
	var (
		newsLetter models.NewsLetter
	)

	newsLetter, err := newsLetter.GetDeletedNewsLetterById(db, ID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if !newsLetter.DeletedAt.Valid {
		return http.StatusBadRequest, errors.New("newsletter email is not soft-deleted")
	}
	newsLetter.DeletedAt = gorm.DeletedAt{}

	if err := newsLetter.UpdateNewsLetter(db); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
