package seed

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func GetUser(c *gin.Context, db *gorm.DB) (models.User, error) {

	var userResp models.User

	// Get the user_id from the URL
	userIDStr := c.Param("user_id")

	// convert the string id to integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return userResp, err
	}

	userEmail := ""

	if userID == 1 {
		userEmail = "john@example.com"
	} else if userID == 2 {
		userEmail = "jane@example.com"
	} else {
		return userResp, errors.New("user id does not exist")
	}

	if err := db.Preload("Profile").Preload("Products").Preload("Organisations").Where("email = ?", userEmail).First(&userResp).Error; err != nil {
		return userResp, err
	}

	return userResp, nil

}
