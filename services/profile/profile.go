package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func UpdateProfile(req models.UpdateProfileRequest, userId string, db *gorm.DB) (gin.H, int, error) {

	var (
		user    models.User
		profile models.Profile
	)
	pdb := inst.InitDB(db)

	profileId, err := user.GetProfileID(pdb, userId)

	if err != nil {
		return gin.H{}, http.StatusNotFound, err
	}

	err = profile.UpdateProfile(pdb, req, profileId)

	if err != nil {
		return gin.H{}, http.StatusInternalServerError, err
	}

	responseData := gin.H{
		"message": "Profile updated successfully",
	}
	return responseData, http.StatusOK, nil
}
