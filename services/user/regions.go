package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func UpdateARegion(userData models.UserRegionTimezoneLanguage, userIDStr string,
	db *gorm.DB, c *gin.Context) (*models.UserRegionTimezoneLanguage, int, error) {
	var (
		currentUser models.User
		regionData  models.UserRegionTimezoneLanguage
		theData     models.UserRegionTimezoneLanguage
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

	if theData, err = regionData.GetUserRegionByID(db, userIDStr); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			userData.UserID = userIDStr
			userData.ID = utility.GenerateUUID()

			if err := userData.CreateUserRegion(db); err != nil {

				return nil, http.StatusBadRequest, err
			}
			return &userData, http.StatusOK, nil

		} else {

			return &theData, http.StatusBadRequest, err
		}
	} else {

		theData.LanguageID = userData.LanguageID
		theData.RegionID = userData.RegionID
		theData.TimezoneID = userData.TimezoneID

		if err := theData.UpdateUserRegion(db); err != nil {
			return &theData, http.StatusBadRequest, err
		}
		return &theData, http.StatusOK, nil
	}

}

func GetUserRegion(userIDStr string,
	db *gorm.DB, c *gin.Context) (*models.UserRegionTimezoneLanguage, int, error) {
	var (
		currentUser models.User
		regionData  models.UserRegionTimezoneLanguage
		theData     models.UserRegionTimezoneLanguage
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

	if theData, err = regionData.GetUserRegionByID(db, userIDStr); err != nil {
		return &theData, http.StatusBadRequest, err
	}

	return &theData, http.StatusOK, nil
}
