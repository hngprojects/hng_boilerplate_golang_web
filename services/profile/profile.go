package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func UpdateProfile(req models.UpdateProfileRequest, userId string, db *gorm.DB) (gin.H, int, error) {

	var (
		user    models.User
		profile models.Profile
	)

	profileId, err := user.GetProfileID(db, userId)

	if err != nil {
		return gin.H{}, http.StatusNotFound, err
	}

	err = profile.UpdateProfile(db, req, profileId)

	if err != nil {
		return gin.H{}, http.StatusInternalServerError, err
	}

	responseData := gin.H{
		"message": "Profile updated successfully",
	}
	return responseData, http.StatusOK, nil
}
