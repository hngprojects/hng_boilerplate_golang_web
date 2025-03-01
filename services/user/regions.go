package user

import (
	"errors"
	"net/http"

	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func (s *userService) UpdateARegion(userData models.UserRegionTimezoneLanguage, userIDStr string, requesterID string) (*models.UserRegionTimezoneLanguage, int, error) {
	var (
		currentUser models.User
		regionData  models.UserRegionTimezoneLanguage
		theData     models.UserRegionTimezoneLanguage
	)

	currentUser, code, err := s.GetUser(requesterID)
	if err != nil {
		return &theData, code, err
	}

	_, code, err = s.GetUser(userIDStr)
	if err != nil {
		return &theData, code, err
	}

	pdb := inst.InitDB(s.db)
	isSuperAdmin := currentUser.CheckUserIsAdmin(pdb)
	if !isSuperAdmin && requesterID != userIDStr {
		return &theData, http.StatusForbidden, errors.New("user does not have permission to update this user")
	}

	if theData, err = regionData.GetUserRegionByID(pdb, userIDStr); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			userData.UserID = userIDStr
			userData.ID = utility.GenerateUUID()

			if err := userData.CreateUserRegion(pdb); err != nil {

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

		if err := theData.UpdateUserRegion(pdb); err != nil {
			return &theData, http.StatusBadRequest, err
		}
		return &theData, http.StatusOK, nil
	}

}

func (s *userService) GetUserRegion(userIDStr string, requesterID string) (*models.UserRegionTimezoneLanguage, int, error) {
	var (
		currentUser models.User
		regionData  models.UserRegionTimezoneLanguage
		theData     models.UserRegionTimezoneLanguage
	)

	currentUser, code, err := s.GetUser(requesterID)
	if err != nil {
		return &theData, code, err
	}

	_, code, err = s.GetUser(userIDStr)
	if err != nil {
		return &theData, code, err
	}
	pdb := inst.InitDB(s.db)

	isSuperAdmin := currentUser.CheckUserIsAdmin(pdb)
	if !isSuperAdmin && requesterID != userIDStr {
		return &theData, http.StatusForbidden, errors.New("user does not have permission to update this user")
	}

	if theData, err = regionData.GetUserRegionByID(pdb, userIDStr); err != nil {
		return &theData, http.StatusBadRequest, err
	}

	return &theData, http.StatusOK, nil
}
