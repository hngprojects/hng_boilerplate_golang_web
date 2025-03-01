package user

import (
	"errors"
	"net/http"

	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func (s *userService) GetUserDataPrivacySettings(userIDStr string, requesterID string) (*models.DataPrivacySettings, int, error) {
	var (
		currentUser models.User
		privacyData models.DataPrivacySettings
		theData     models.DataPrivacySettings
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
		return &theData, http.StatusForbidden, errors.New("user does not have permission to view this user's privacy settings")
	}

	if theData, err = privacyData.GetUserDataPrivacySettingsByID(pdb, userIDStr); err != nil {
		if err.Error() == "record not found" {
			theModel := models.DataPrivacySettings{
				UserID: userIDStr,
			}
			err := theModel.CreateDataPrivacySettings(pdb)
			if err != nil {
				return nil, http.StatusBadRequest, err
			}

			if theData, err = privacyData.GetUserDataPrivacySettingsByID(pdb, userIDStr); err != nil {
				return nil, http.StatusBadRequest, err
			}
			return &theData, http.StatusCreated, nil
		}
		return &theData, http.StatusBadRequest, err
	}

	return &theData, http.StatusOK, nil
}

func (s *userService) UpdateUserDataPrivacySettings(userData models.DataPrivacySettings, userIDStr string, requesterID string) (*models.DataPrivacySettings, int, error) {
	var (
		currentUser models.User
		privacyData models.DataPrivacySettings
		theData     models.DataPrivacySettings
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

	if theData, err = privacyData.GetUserDataPrivacySettingsByID(pdb, userIDStr); err != nil {
		return &theData, http.StatusBadRequest, err
	} else {

		theData.AllowAnalytics = userData.AllowAnalytics
		theData.Enable2FA = userData.Enable2FA
		theData.PersonalizedAds = userData.PersonalizedAds
		theData.ProfileVisibility = userData.ProfileVisibility
		theData.UseDataEncryption = userData.UseDataEncryption
		theData.ShareDataWithPartners = userData.ShareDataWithPartners
		theData.ReceiveEmailUpdates = userData.ReceiveEmailUpdates

		if err := theData.UpdateDataPrivacySettings(pdb); err != nil {
			return &theData, http.StatusBadRequest, err
		}
		return &theData, http.StatusOK, nil
	}

}
