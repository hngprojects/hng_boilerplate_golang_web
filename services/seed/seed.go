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
