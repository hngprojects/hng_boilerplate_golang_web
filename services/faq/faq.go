package service

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"gorm.io/gorm"
)

func GetFaq(c *gin.Context, db *gorm.DB) ([]models.FAQ, *database.PaginationResponse, int, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)

	var faq models.FAQ

	faqs, paginationResponse, err := faq.FetchAllFaq(pdb, c)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return faqs, nil, http.StatusNoContent, nil
		}
		return faqs, nil, http.StatusBadRequest, err

	}

	return faqs, &paginationResponse, http.StatusOK, nil

}

func AddToFaq(faq *models.FAQ, db *gorm.DB) error {

	// instance of Postgresql db
	pdb := inst.InitDB(db)
	if pdb.CheckExists(faq, "question = ?", faq.Question) {
		return errors.New("question exists")
	}

	if err := faq.CreateFaq(pdb); err != nil {
		return err
	}

	return nil
}

func DeleteFaq(ID string, db *gorm.DB) (int, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)
	var (
		faq models.FAQ
	)

	faq, err := faq.GetFaqById(pdb, ID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err := faq.DeleteFaq(pdb); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func UpdateFaq(req models.UpdateFAQ, ID string, db *gorm.DB) (models.FAQ, int, error) {
	var (
		faq models.FAQ
	)

	// instance of Postgresql db
	pdb := inst.InitDB(db)

	faq, err := faq.GetFaqById(pdb, ID)
	if err != nil {
		return faq, http.StatusBadRequest, err
	}

	faq.Question = req.Question
	faq.Answer = req.Answer

	if err := faq.UpdateFaq(pdb); err != nil {
		return faq, http.StatusBadRequest, err
	}

	return faq, http.StatusOK, nil
}
