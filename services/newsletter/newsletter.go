package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"gorm.io/gorm"
)

func GetNewsletters(c *gin.Context, db *gorm.DB) ([]models.NewsLetter, *database.PaginationResponse, int, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)

	var newsletter models.NewsLetter

	newsLetters, paginationResponse, err := newsletter.FetchAllNewsLetter(pdb, c)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return newsLetters, nil, http.StatusNoContent, nil
		}
		return newsLetters, nil, http.StatusBadRequest, err

	}

	return newsLetters, &paginationResponse, http.StatusOK, nil

}

func GetDeletedNewsletters(c *gin.Context, db *gorm.DB) ([]models.NewsLetter, *database.PaginationResponse, int, error) {

	// instance of Postgresql db
	pdb := inst.InitDB(db)
	var newsletter models.NewsLetter

	delNewsLetters, paginationResponse, err := newsletter.FetchAllDeletedNewsLetter(pdb, c)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return delNewsLetters, nil, http.StatusNoContent, nil
		}
		return delNewsLetters, nil, http.StatusBadRequest, err

	}

	return delNewsLetters, &paginationResponse, http.StatusOK, nil

}

func NewsLetterSubscribe(newsletter *models.NewsLetter, db *gorm.DB) error {

	// instance of Postgresql db
	pdb := inst.InitDB(db)
	if pdb.CheckExists(newsletter, "email = ?", newsletter.Email) {
		return models.ErrEmailAlreadySubscribed
	}

	newsletter.Email = strings.ToLower(newsletter.Email)

	if err := newsletter.CreateNewsLetter(pdb); err != nil {
		return err
	}

	return nil
}

func DeleteNewsLetter(ID string, db *gorm.DB, c *gin.Context) (int, error) {
	var (
		newsLetter models.NewsLetter
	)
	// instance of Postgresql db
	pdb := inst.InitDB(db)

	newsLetter, err := newsLetter.GetNewsLetterById(pdb, ID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err := newsLetter.DeleteNewsLetter(pdb); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func RestoreDeletedNewsLetter(ID string, db *gorm.DB, c *gin.Context) (int, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)
	var (
		newsLetter models.NewsLetter
	)

	newsLetter, err := newsLetter.GetDeletedNewsLetterById(pdb, ID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if !newsLetter.DeletedAt.Valid {
		return http.StatusBadRequest, errors.New("newsletter email is not soft-deleted")
	}
	newsLetter.DeletedAt = gorm.DeletedAt{}

	if err := newsLetter.UpdateNewsLetter(pdb); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
