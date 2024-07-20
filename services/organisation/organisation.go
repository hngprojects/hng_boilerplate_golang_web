package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func ValidateCreateOrgRequest(req models.CreateOrgRequestModel, db *gorm.DB) (models.CreateOrgRequestModel, error) {

	org := models.Organisation{}

	// Check if the organization already exists
	exists := postgresql.CheckExists(db, &org, "email = ?", req.Email)
	if exists {
		return req, errors.New("organization already exists with the given email")
	}

	return req, nil
}

func CreateOrganisation(req models.CreateOrgRequestModel, db *gorm.DB, userId string) (*models.Organisation, error) {

	org := models.Organisation{
		ID:          utility.GenerateUUID(),
		Name:        strings.ToLower(req.Name),
		Description: strings.ToLower(req.Description),
		Email:       strings.ToLower(req.Email),
		State:       strings.ToLower(req.State),
		Industry:    strings.ToLower(req.Industry),
		Type:        strings.ToLower(req.Type),
		OwnerID:     userId,
		Address:     strings.ToLower(req.Address),
		Country:     strings.ToLower(req.Country),
	}

	err := org.CreateOrganisation(db)

	if err != nil {
		return nil, err
	}

	user, err := models.GetUserByID(db, userId)

	if err != nil {
		return nil, err
	}

	err = models.AddUserToOrganisation(db, &user, []interface{}{&org})

	if err != nil {
		return nil, err
	}

	return &org, nil
}
