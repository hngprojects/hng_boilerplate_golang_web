package service

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/actions"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/actions/names"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func GetAllContactUs(c *gin.Context, db *gorm.DB) ([]models.ContactUs, *database.PaginationResponse, int, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)

	var contact models.ContactUs

	contacts, paginationResponse, err := contact.FetchAllContactUs(pdb, c)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contacts, nil, http.StatusNoContent, nil
		}
		return contacts, nil, http.StatusBadRequest, err

	}

	return contacts, &paginationResponse, http.StatusOK, nil

}

func AddToContactUs(contact *models.ContactUs, db *gorm.DB) error {
	// instance of Postgresql db
	pdb := inst.InitDB(db)

	contact.Message = utility.CleanStringInput(contact.Message)

	if err := contact.CreateContactUs(pdb); err != nil {
		return err
	}

	msgReq := models.ContactUs{
		Email:   contact.Email,
		Name:    contact.Name,
		Message: contact.Message,
	}

	err := actions.AddNotificationToQueue(storage.DB.Redis.RedisDb(), names.SendContactUsMail, msgReq)
	if err != nil {
		return err
	}

	return nil
}

func DeleteContactUs(ID string, db *gorm.DB) (int, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)
	var (
		contact models.ContactUs
	)

	contact, err := contact.GetContactUsById(pdb, ID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err := contact.DeleteContactUs(pdb); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func GetContactUsById(ID string, db *gorm.DB) (*models.ContactUs, error) {

	// instance of Postgresql db
	pdb := inst.InitDB(db)
	var contact models.ContactUs

	contactData, err := contact.GetContactUsById(pdb, ID)
	if err != nil {
		return nil, err
	}

	return &contactData, nil
}

func GetContactUsByEmail(email string, db *gorm.DB) (*[]models.ContactUs, error) {

	// instance of Postgresql db
	pdb := inst.InitDB(db)
	var contact models.ContactUs

	contactData, err := contact.GetContactUsByEmail(pdb, email)
	if err != nil {
		return nil, err
	}

	return &contactData, nil
}
