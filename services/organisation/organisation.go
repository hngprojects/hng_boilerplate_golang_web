package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func ValidateCreateOrgRequest(req models.CreateOrgRequestModel, db *gorm.DB) (models.CreateOrgRequestModel, int, error) {

	org := models.Organisation{}

	// Check if the organization already exists

	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)
		formattedMail, checkBool := utility.EmailValid(req.Email)
		if !checkBool {
			return req, http.StatusUnprocessableEntity, fmt.Errorf("email address is invalid")
		}
		req.Email = formattedMail
		exists := postgresql.CheckExists(db, &org, "email = ?", req.Email)
		if exists {
			return req, http.StatusBadRequest, errors.New("organization already exists with the given email")
		}
	}

	return req, 0, nil
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

	var user models.User

	user, err = user.GetUserByID(db, userId)

	if err != nil {
		return nil, err
	}

	err = user.AddUserToOrganisation(db, &user, []interface{}{&org})

	if err != nil {
		return nil, err
	}

	return &org, nil
}

func GetOrganisationById(orgId string, userId string, db *gorm.DB) (models.Organisation, error) {
	var org models.Organisation

	org, err := org.GetActiveOrganisationById(db, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Organisation{}, errors.New("organisation not found")
		}
		return models.Organisation{}, err
	}

	if org.OwnerID != userId {
		return models.Organisation{}, errors.New("user not authorised to retrieve this organisation")
	}

	return org, nil
}

func UpdateOrganisation(orgId string, userId string, updateReq models.CreateOrgRequestModel, db *gorm.DB) (*models.Organisation, error) {
	var org models.Organisation

	org, err := org.GetActiveOrganisationById(db, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organisation not found")
		}
		return nil, err
	}

	if org.OwnerID != userId {
		return nil, errors.New("user not authorised to update this organisation")
	}

	if updateReq.Email != "" && updateReq.Email != org.Email {
		updateReq.Email = strings.ToLower(updateReq.Email)
		formattedMail, checkBool := utility.EmailValid(updateReq.Email)
		if !checkBool {
			return nil, errors.New("email address is invalid")
		}
		updateReq.Email = formattedMail
		exists := postgresql.CheckExists(db, &org, "email = ?", updateReq.Email)
		if exists {
			return nil, errors.New("organisation already exists with the given email")
		}
	}

	org.Name = updateReq.Name
    org.Description = updateReq.Description
    org.Email = updateReq.Email
    org.State = updateReq.State
    org.Industry = updateReq.Industry
    org.Type = updateReq.Type
    org.Address = updateReq.Address
    org.Country = updateReq.Country

	return org.Update(db)
}

func DeleteOrganisation(orgId string, userId string, db *gorm.DB) error {
	var org models.Organisation

	org, err := org.GetActiveOrganisationById(db, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organisation not found")
		}
		return err
	}

	if org.OwnerID != userId {
		return errors.New("user not authorised to delete this organisation")
	}

	return org.Delete(db)
}
