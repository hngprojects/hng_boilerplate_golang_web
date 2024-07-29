package organisation

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

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

func GetOrganisation(orgId string, userId string, db *gorm.DB) (*models.Organisation, error) {
	org, err := CheckOrgExists(orgId, db)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organisation not found")
		}
		return nil, err
	}

	isMember, err := CheckUserIsMemberOfOrg(userId, orgId, db)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user not authorised to retrieve this organisation")
	}

	return &org, nil
}

func UpdateOrganisation(orgId string, userId string, updateReq models.UpdateOrgRequestModel, db *gorm.DB) (*models.Organisation, error) {
	org, err := CheckOrgExists(orgId, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organisation not found")
		}
		return nil, err
	}

	isMember, err := CheckUserIsMemberOfOrg(userId, orgId, db)
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
		exists := postgresql.CheckExists(db, &org, "email = ?", updateReq.Email)
		if exists {
			return nil, errors.New("organisation already exists with the given email")
		}
	}

	// Update the organization fields
	if updateReq.Name != "" {
		org.Name = updateReq.Name
	}
	if updateReq.Description != "" {
		org.Description = updateReq.Description
	}
	if updateReq.Email != "" {
		org.Email = updateReq.Email
	}
	if updateReq.State != "" {
		org.State = updateReq.State
	}
	if updateReq.Industry != "" {
		org.Industry = updateReq.Industry
	}
	if updateReq.Type != "" {
		org.Type = updateReq.Type
	}
	if updateReq.Address != "" {
		org.Address = updateReq.Address
	}
	if updateReq.Country != "" {
		org.Country = updateReq.Country
	}

	return org.Update(db)
}

func DeleteOrganisation(orgId string, userId string, db *gorm.DB) error {
	org, err := CheckOrgExists(orgId, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organisation not found")
		}
		return err
	}

	isMember, err := CheckUserIsMemberOfOrg(userId, orgId, db)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("user not authorised to delete this organisation")
	}

	return org.Delete(db)
}

func AddUserToOrganisation(orgId string, req models.AddUserToOrgRequestModel, db *gorm.DB) error {
	var user models.User
	org, err := CheckOrgExists(orgId, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organisation not found")
		}
		return err
	}

	user, err = user.GetUserByID(db, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	isMember, err := CheckUserIsMemberOfOrg(req.UserId, orgId, db)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("user already added to organisation")
	}

	err = user.AddUserToOrganisation(db, &user, []interface{}{&org})

	if err != nil {
		return err
	}

	return nil

}

func GetUsersInOrganisation(orgId string, userId string, db *gorm.DB, c *gin.Context) ([]UserResponse, postgresql.PaginationResponse, error) {
	var users []UserResponse
	_, err := CheckOrgExists(orgId, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, postgresql.PaginationResponse{}, errors.New("organisation not found")
		}
		return nil, postgresql.PaginationResponse{}, err
	}

	isMember, err := CheckUserIsMemberOfOrg(userId, orgId, db)
	if err != nil {
		return nil, postgresql.PaginationResponse{}, err
	}
	if !isMember {
		return nil, postgresql.PaginationResponse{}, errors.New("user does not have access to the organisation")
	}

	pagination := postgresql.GetPagination(c)

	offset := (pagination.Page - 1) * pagination.Limit

	if err := db.Table("users").
		Select("users.id, users.email, profiles.phone as phone_number , users.name").
		Joins("JOIN user_organisations ON user_organisations.user_id = users.id").
		Joins("JOIN profiles ON profiles.userid = users.id").
		Where("user_organisations.organisation_id = ?", orgId).
		Offset(offset).
		Limit(pagination.Limit).
		Find(&users).Error; err != nil {
		return nil, postgresql.PaginationResponse{}, err
	}

	var totalUsers int64
	if err := db.Table("users").
		Joins("JOIN user_organisations ON user_organisations.user_id = users.id").
		Joins("JOIN profiles ON profiles.userid = users.id").
		Where("user_organisations.organisation_id = ?", orgId).
		Count(&totalUsers).Error; err != nil {
		return nil, postgresql.PaginationResponse{}, err
	}

	totalPages := int(math.Ceil(float64(totalUsers) / float64(pagination.Limit)))
	paginationResponse := postgresql.PaginationResponse{
		CurrentPage:     pagination.Page,
		PageCount:       pagination.Limit,
		TotalPagesCount: totalPages,
	}

	return users, paginationResponse, nil
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

func CheckOrgExists(orgId string, db *gorm.DB) (models.Organisation, error) {
	var org models.Organisation

	org, err := org.GetOrgByID(db, orgId)
	if err != nil {
		return org, err
	}

	return org, nil
}
