package organisation

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

	return req, 0,  nil
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


func CheckOrgExists(orgId string, db *gorm.DB) (models.Organisation, error) {
	var org models.Organisation

	org, err := org.GetOrgByID(db, orgId)
	if err != nil {
		return org, err
	}

	return org, nil
}

func CheckUserIsMemberOfOrg(userId string, orgId string, db *gorm.DB) (bool, error) {
	var org models.Organisation
	var user models.User

	org, err := org.GetOrgByID(db, orgId)
	if err != nil {
		return false, err
	}

	user, err = user.GetUserByID(db, userId)
	if err != nil {
		return false, err
	}

	for _, org := range user.Organisations {
		if org.ID == orgId {
			return true, nil
		}
	}

	return false, nil
}

