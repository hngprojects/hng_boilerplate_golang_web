package organisation

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func ValidateCreateOrgRequest(req models.CreateOrgRequestModel, db *gorm.DB) (models.CreateOrgRequestModel, int, error) {

	pdb := inst.InitDB(db)
	org := models.Organisation{}

	// Check if the organization already exists

	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)
		formattedMail, checkBool := utility.EmailValid(req.Email)
		if !checkBool {
			return req, http.StatusUnprocessableEntity, fmt.Errorf("email address is invalid")
		}
		req.Email = formattedMail
		exists := pdb.CheckExists(&org, "email = ?", req.Email)
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

	pdb := inst.InitDB(db)
	err := org.CreateOrganisation(pdb)

	if err != nil {
		return nil, err
	}

	var user models.User

	user, err = user.GetUserByID(pdb, userId)

	if err != nil {
		return nil, err
	}

	err = user.AddUserToOrganisation(pdb, &user, []interface{}{&org})

	if err != nil {
		return nil, err
	}

	return &org, nil
}

func GetOrganisation(orgId string, userId string, db *gorm.DB) (*models.Organisation, error) {
	var org models.Organisation
	pdb := inst.InitDB(db)
	org, err := org.CheckOrgExists(orgId, pdb)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organisation not found")
		}
		return nil, err
	}

	isMember, err := org.CheckUserIsMemberOfOrg(userId, orgId, pdb)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user not authorised to retrieve this organisation")
	}

	return &org, nil
}

func UpdateOrganisation(orgId string, userId string, updateReq models.UpdateOrgRequestModel, db *gorm.DB) (*models.Organisation, error) {
	var org models.Organisation
	pdb := inst.InitDB(db)
	org, err := org.CheckOrgExists(orgId, pdb)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organisation not found")
		}
		return nil, err
	}

	isMember, err := org.CheckUserIsMemberOfOrg(userId, orgId, pdb)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user not authorised to update this organisation")
	}

	if updateReq.Email != "" && updateReq.Email != org.Email {
		updateReq.Email = strings.ToLower(updateReq.Email)
		formattedMail, checkBool := utility.EmailValid(updateReq.Email)
		if !checkBool {
			return nil, errors.New("email address is invalid")
		}
		updateReq.Email = formattedMail
		exists := pdb.CheckExists(&org, "email = ?", updateReq.Email)
		if exists {
			return nil, errors.New("organisation already exists with the given email")
		}
	}

	return org.Update(pdb, updateReq, orgId)
}

func DeleteOrganisation(orgId string, userId string, db *gorm.DB) error {
	var org models.Organisation
	pdb := inst.InitDB(db)
	org, err := org.CheckOrgExists(orgId, pdb)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organisation not found")
		}
		return err
	}

	isMember, err := org.CheckUserIsMemberOfOrg(userId, orgId, pdb)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("user not authorised to delete this organisation")
	}

	return org.Delete(pdb)
}

func AddUserToOrganisation(orgId string, req models.AddUserToOrgRequestModel, db *gorm.DB) error {
	var user models.User
	var org models.Organisation
	pdb := inst.InitDB(db)
	org, err := org.CheckOrgExists(orgId, pdb)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organisation not found")
		}
		return err
	}

	user, err = user.GetUserByID(pdb, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	isMember, err := org.CheckUserIsMemberOfOrg(req.UserId, orgId, pdb)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("user already added to organisation")
	}

	err = user.AddUserToOrganisation(pdb, &user, []interface{}{&org})

	if err != nil {
		return err
	}

	return nil

}

func GetUsersInOrganisation(orgId string, userId string, db *gorm.DB, c *gin.Context) ([]models.UserInOrgResponse, database.PaginationResponse, error) {
	var org models.Organisation
	pdb := inst.InitDB(db)
	_, err := org.CheckOrgExists(orgId, pdb)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.PaginationResponse{}, errors.New("organisation not found")
		}
		return nil, database.PaginationResponse{}, err
	}

	isMember, err := org.CheckUserIsMemberOfOrg(userId, orgId, pdb)
	if err != nil {
		return nil, database.PaginationResponse{}, err
	}
	if !isMember {
		return nil, database.PaginationResponse{}, errors.New("user does not have access to the organisation")
	}

	users, paginationResponse, err := org.GetUsersInOrganisation(c, pdb, orgId)

	if err != nil {
		return nil, database.PaginationResponse{}, err
	}

	return users, paginationResponse, nil
}
