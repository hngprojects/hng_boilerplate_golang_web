package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"gorm.io/gorm"
)

func GetUserDataPrivacySettings(userIDStr string,
	db *gorm.DB, c *gin.Context) (*models.DataPrivacySettings, int, error) {
	var (
		currentUser models.User
		privacyData models.DataPrivacySettings
		theData     models.DataPrivacySettings
	)

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return &theData, http.StatusNotFound, err
	}

	currentUserID, ok := userId.(string)
	if !ok {
		return &theData, http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	currentUser, code, err := GetUser(currentUserID, db)
	if err != nil {
		return &theData, code, err
	}

	_, code, err = GetUser(userIDStr, db)
	if err != nil {
		return &theData, code, err
	}

	isSuperAdmin := currentUser.CheckUserIsAdmin(db)
	if !isSuperAdmin && currentUserID != userIDStr {
		return &theData, http.StatusForbidden, errors.New("user does not have permission to view this user's privacy settings")
	}

	if theData, err = privacyData.GetUserDataPrivacySettingsByID(db, userIDStr); err != nil {
		if err.Error() == "record not found" {
			theModel := models.DataPrivacySettings{
				UserID: userIDStr,
			}
			err := theModel.CreateDataPrivacySettings(db)
			if err != nil {
				return nil, http.StatusBadRequest, err
			}

			if theData, err = privacyData.GetUserDataPrivacySettingsByID(db, userIDStr); err != nil {
				return nil, http.StatusBadRequest, err
			}
			return &theData, http.StatusCreated, nil
		}
		return &theData, http.StatusBadRequest, err
	}

	return &theData, http.StatusOK, nil
}

func UpdateUserDataPrivacySettings(userData models.DataPrivacySettings, userIDStr string,
	db *gorm.DB, c *gin.Context) (*models.DataPrivacySettings, int, error) {
	var (
		currentUser models.User
		privacyData models.DataPrivacySettings
		theData     models.DataPrivacySettings
	)

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return &theData, http.StatusNotFound, err
	}

	currentUserID, ok := userId.(string)
	if !ok {
		return &theData, http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	currentUser, code, err := GetUser(currentUserID, db)
	if err != nil {
		return &theData, code, err
	}

	_, code, err = GetUser(userIDStr, db)
	if err != nil {
		return &theData, code, err
	}

	isSuperAdmin := currentUser.CheckUserIsAdmin(db)
	if !isSuperAdmin && currentUserID != userIDStr {
		return &theData, http.StatusForbidden, errors.New("user does not have permission to update this user")
	}

	if theData, err = privacyData.GetUserDataPrivacySettingsByID(db, userIDStr); err != nil {
		return &theData, http.StatusBadRequest, err
	} else {

		theData.AllowAnalytics = userData.AllowAnalytics
		theData.Enable2FA = userData.Enable2FA
		theData.PersonalizedAds = userData.PersonalizedAds
		theData.ProfileVisibility = userData.ProfileVisibility
		theData.UseDataEncryption = userData.UseDataEncryption
		theData.ShareDataWithPartners = userData.ShareDataWithPartners
		theData.ReceiveEmailUpdates = userData.ReceiveEmailUpdates

		if err := theData.UpdateDataPrivacySettings(db); err != nil {
			return &theData, http.StatusBadRequest, err
		}
		return &theData, http.StatusOK, nil
	}

}
