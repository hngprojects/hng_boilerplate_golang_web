package seed

import (
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func GetUser(userIDStr string, db *gorm.DB) (models.User, error) {
	var userResp models.User

	userResp, err := userResp.GetUserByID(db, userIDStr)
	if err != nil {
		return userResp, err
	}

	return userResp, nil
}

func CheckOrgExists(orgId string, db *gorm.DB) (models.Organisation, error) {
	var org models.Organisation

	org, err := org.GetOrgByID(db, orgId)
	if err != nil {
		return org, err
	}

	return org, nil
}

func GetUserByEmail(email string, db *gorm.DB)(models.User, error){
	var user models.User

	user, err := user.GetUserByEmail(db, email);
	if err != nil {
		return user, err
	}
	return user, nil
}