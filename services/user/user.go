package user

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
)

func (s *userService) GetUser(userIDStr string) (models.User, int, error) {
	var userResp models.User

	pdb := inst.InitDB(s.db)
	userResp, err := userResp.GetUserByID(pdb, userIDStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userResp, http.StatusNotFound, errors.New("user not found")
		}
		return userResp, http.StatusBadRequest, err
	}
	return userResp, http.StatusOK, nil
}

func (s *userService) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	pdb := inst.InitDB(s.db)
	user, err := user.GetUserByEmail(pdb, email)

	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *userService) GetAUser(userIDStr string, requesterID string) (*models.User, int, error) {
	var userResp models.User
	pdb := inst.InitDB(s.db)
	// Fetch the requesting user from the database
	requester, code, err := s.GetUser(requesterID)
	if err != nil {
		return nil, code, err
	}

	isSuperAdmin := requester.CheckUserIsAdmin(pdb)
	if isSuperAdmin {
		userResp, err = userResp.GetUserByID(pdb, userIDStr)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &userResp, http.StatusNotFound, errors.New("user not found")
			}
			return &userResp, http.StatusBadRequest, err
		}
	} else {
		userResp, err = userResp.GetUserByIDsAdmin(pdb, userIDStr, requesterID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &userResp, http.StatusNotFound, errors.New("user not found")
			}
			return &userResp, http.StatusBadRequest, err
		}
	}

	return &userResp, http.StatusOK, nil
}

func (s *userService) GetAUserOrganisation(UserID string, requesterID string) (*[]models.Organisation, int, error) {
	var (
		orgData models.Organisation
		orgResp []models.Organisation
	)
	pdb := inst.InitDB(s.db)

	user, code, err := s.GetUser(requesterID)
	if err != nil {
		return nil, code, err
	}

	isSuperAdmin := user.CheckUserIsAdmin(pdb)
	if isSuperAdmin {
		orgResp, err = orgData.GetOrganisationsByUserID(s.db, UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &orgResp, http.StatusNotFound, errors.New("user not found")
			}
			return &orgResp, http.StatusBadRequest, err
		}
		orgResp, err = orgData.GetOrganisationsByUserIDs(pdb, UserID, requesterID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &orgResp, http.StatusNotFound, errors.New("user not found")
			}
			return &orgResp, http.StatusBadRequest, err
		}
	}

	return &orgResp, http.StatusOK, nil
}

func (s *userService) DeleteAUser(userIDStr string, requesterID string) (int, error) {
	var (
		currentUser models.User
		targetUser  models.User
	)

	currentUser, code, err := s.GetUser(requesterID)
	if err != nil {
		return code, err
	}

	targetUser, code, err = s.GetUser(userIDStr)
	if err != nil {
		return code, err
	}

	pdb := inst.InitDB(s.db)
	isSuperAdmin := currentUser.CheckUserIsAdmin(pdb)
	if isSuperAdmin || requesterID == userIDStr {

		if err := targetUser.DeleteAUser(pdb); err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		return http.StatusForbidden, errors.New("user does not have permission to delete this user")
	}

	return http.StatusOK, nil
}

func (s *userService) UpdateAUser(userData models.UpdateUserRequestModel, userIDStr string, requesterID string) (*models.User, int, error) {
	var (
		currentUser models.User
		targetUser  models.User
	)

	currentUser, code, err := s.GetUser(requesterID)
	if err != nil {
		return &targetUser, code, err
	}

	targetUser, code, err = s.GetUser(userIDStr)
	if err != nil {
		return &targetUser, code, err
	}
	pdb := inst.InitDB(s.db)
	isSuperAdmin := currentUser.CheckUserIsAdmin(pdb)
	if isSuperAdmin || requesterID == userIDStr {

		targetUser.Name = userData.UserName
		targetUser.Profile.FirstName = userData.FirstName
		targetUser.Profile.LastName = userData.LastName
		targetUser.Profile.Phone = userData.PhoneNumber

		err = targetUser.Update(pdb)
		if err != nil {
			return &targetUser, http.StatusInternalServerError, err
		}

	} else {
		return &targetUser, http.StatusForbidden, errors.New("user does not have permission to update this user")
	}

	return &targetUser, http.StatusOK, nil
}
func (s *userService) GetAllUsers(pagination database.Pagination) ([]models.User, *database.PaginationResponse, int, error) {

	var users []models.User

	pdb := inst.InitDB(s.db)
	paginationResponse, err := pdb.SelectAllFromDbOrderByPaginated("created_at", "desc", "", pagination, &users, "deleted_at IS NULL")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users, nil, http.StatusNoContent, nil
		}
		return users, nil, http.StatusBadRequest, err

	}

	return users, &paginationResponse, http.StatusOK, nil

}
