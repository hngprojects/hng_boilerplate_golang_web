package service

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func GetFaq(c *gin.Context, db *gorm.DB) ([]models.FAQ, *postgresql.PaginationResponse, int, error) {

	var faq models.FAQ

	faqs, paginationResponse, err := faq.FetchAllFaq(db, c)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return faqs, nil, http.StatusNoContent, nil
		}
		return faqs, nil, http.StatusBadRequest, err

	}

	return faqs, &paginationResponse, http.StatusOK, nil

}

func AddToFaq(faq *models.FAQ, db *gorm.DB) error {

	if postgresql.CheckExists(db, faq, "question = ?", faq.Question) {
		return errors.New("question exists")
	}

	if err := faq.CreateFaq(db); err != nil {
		return err
	}

	return nil
}

func DeleteFaq(ID string, db *gorm.DB) (int, error) {
	var (
		faq models.FAQ
	)

	faq, err := faq.GetFaqById(db, ID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err := faq.DeleteFaq(db); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func UpdateFaq(req models.UpdateFAQ, ID string, db *gorm.DB) (models.FAQ, int, error) {
	var (
		faq models.FAQ
	)

	faq, err := faq.GetFaqById(db, ID)
	if err != nil {
		return faq, http.StatusBadRequest, err
	}

	faq.Question = req.Question
	faq.Answer = req.Answer

	if err := faq.UpdateFaq(db); err != nil {
		return faq, http.StatusBadRequest, err
	}

	return faq, http.StatusOK, nil
}
