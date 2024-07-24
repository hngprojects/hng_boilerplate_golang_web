package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func UpdateUserPassword(c *gin.Context, req models.ChangePasswordRequestModel, db *gorm.DB) (*models.User, int, error) {

	user := models.User{}

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	userID, ok := userId.(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	userDataExist, err := user.GetUserByID(db, userID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("unable to fetch user " + err.Error())
	}

	if !utility.CompareHash(req.OldPassword, userDataExist.Password) {
		return nil, http.StatusBadRequest, fmt.Errorf("old password is incorrect")
	}

	if req.OldPassword == req.NewPassword {
		return nil, http.StatusConflict, errors.New("new password cannot be the same as the old password")
	}

	hashedPassword, err := utility.HashPassword(req.NewPassword)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	userDataExist.Password = hashedPassword
	err = userDataExist.Update(db)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &userDataExist, http.StatusOK, nil
}
