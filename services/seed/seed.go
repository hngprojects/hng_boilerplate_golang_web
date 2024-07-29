package seed

import (
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func GetUser(userIDStr string, db *gorm.DB) ([]models.User, error) {
	var (
		user models.User
		userResp []models.User
	)


	userResp, err := user.GetSeedUsers(db)
	if err != nil {
		return userResp, err
	}

	return userResp, nil
}
