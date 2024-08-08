package user

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

func GetUser(userIDStr string, db *gorm.DB) (models.User, int, error) {
	var userResp models.User

	userResp, err := userResp.GetUserByID(db, userIDStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userResp, http.StatusNotFound, errors.New("user not found")
		}
		return userResp, http.StatusBadRequest, err
	}
	return userResp, http.StatusOK, nil
}

func GetUserByEmail(email string, db *gorm.DB) (models.User, error) {
	var user models.User

	user, err := user.GetUserByEmail(db, email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetAUser(userIDStr string, db *gorm.DB, c *gin.Context) (*models.User, int, error) {
	var userResp models.User

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	userID, ok := userId.(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	user, code, err := GetUser(userID, db)
	if err != nil {
		return nil, code, err
	}

	isSuperAdmin := user.CheckUserIsAdmin(db)
	if isSuperAdmin {
		userResp, err = userResp.GetUserByID(db, userIDStr)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &userResp, http.StatusNotFound, errors.New("user not found")
			}
			return &userResp, http.StatusBadRequest, err
		}
	} else {
		userResp, err = userResp.GetUserByIDsAdmin(db, userIDStr, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &userResp, http.StatusNotFound, errors.New("user not found")
			}
			return &userResp, http.StatusBadRequest, err
		}
	}

	return &userResp, http.StatusOK, nil
}

func GetAUserOrganisation(userIDStr string, db *gorm.DB, c *gin.Context) (*[]models.Organisation, int, error) {
	var (
		orgData models.Organisation
		orgResp []models.Organisation
	)

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	userID, ok := userId.(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	user, code, err := GetUser(userID, db)
	if err != nil {
		return nil, code, err
	}

	isSuperAdmin := user.CheckUserIsAdmin(db)
	if isSuperAdmin {
		orgResp, err = orgData.GetOrganisationsByUserID(db, userIDStr)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &orgResp, http.StatusNotFound, errors.New("user not found")
			}
			return &orgResp, http.StatusBadRequest, err
		}
	} else {
		orgResp, err = orgData.GetOrganisationsByUserIDs(db, userIDStr, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &orgResp, http.StatusNotFound, errors.New("user not found")
			}
			return &orgResp, http.StatusBadRequest, err
		}
	}

	return &orgResp, http.StatusOK, nil
}

func DeleteAUser(userIDStr string, db *gorm.DB, c *gin.Context) (int, error) {
	var (
		currentUser models.User
		targetUser  models.User
	)

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return http.StatusNotFound, err
	}

	currentUserID, ok := userId.(string)
	if !ok {
		return http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	currentUser, code, err := GetUser(currentUserID, db)
	if err != nil {
		return code, err
	}

	targetUser, code, err = GetUser(userIDStr, db)
	if err != nil {
		return code, err
	}

	isSuperAdmin := currentUser.CheckUserIsAdmin(db)
	if isSuperAdmin || currentUserID == userIDStr {

		if err := targetUser.DeleteAUser(db); err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		return http.StatusForbidden, errors.New("user does not have permission to delete this user")
	}

	return http.StatusOK, nil
}

func UpdateAUser(userData models.UpdateUserRequestModel, userIDStr string, db *gorm.DB, c *gin.Context) (*models.User, int, error) {
	var (
		currentUser models.User
		targetUser  models.User
	)

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return &targetUser, http.StatusNotFound, err
	}

	currentUserID, ok := userId.(string)
	if !ok {
		return &targetUser, http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	currentUser, code, err := GetUser(currentUserID, db)
	if err != nil {
		return &targetUser, code, err
	}

	targetUser, code, err = GetUser(userIDStr, db)
	if err != nil {
		return &targetUser, code, err
	}

	isSuperAdmin := currentUser.CheckUserIsAdmin(db)
	if isSuperAdmin || currentUserID == userIDStr {

		targetUser.Name = userData.UserName
		targetUser.Profile.FirstName = userData.FirstName
		targetUser.Profile.LastName = userData.LastName
		targetUser.Profile.Phone = userData.PhoneNumber

		err = targetUser.Update(db)
		if err != nil {
			return &targetUser, http.StatusInternalServerError, err
		}

	} else {
		return &targetUser, http.StatusForbidden, errors.New("user does not have permission to update this user")
	}

	return &targetUser, http.StatusOK, nil
}

func GetAllUsers(c *gin.Context, db *gorm.DB) ([]models.User, *postgresql.PaginationResponse, int, error) {

	var users []models.User
	pagination := postgresql.GetPagination(c)

	paginationResponse, err := postgresql.SelectAllFromDbOrderByPaginated(db, "created_at", "desc", pagination, &users, "deleted_at IS NULL")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users, nil, http.StatusNoContent, nil
		}
		return users, nil, http.StatusBadRequest, err

	}

	return users, &paginationResponse, http.StatusOK, nil

}
