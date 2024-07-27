package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
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

func PasswordReset(userEmail string, db *gorm.DB) (string, int, error) {

	var (
		user      = models.User{}
		passReset = models.PasswordReset{}
		config    = config.GetConfig()
	)

	resetExist, err := passReset.GetPasswordResetByEmail(db, userEmail)
	if err != nil {
		return "error", http.StatusUnauthorized, err
	}

	if resetExist != nil {
		if err := resetExist.DeletePasswordReset(db); err != nil {
			return "error", http.StatusInternalServerError, err
		}
	}

	exists := postgresql.CheckExists(db, &user, "email = ?", userEmail)
	if !exists {
		return "error", http.StatusNotFound, fmt.Errorf("user not found")
	}

	resetToken := utility.GenerateUUID()
	reset := models.PasswordReset{
		ID:        utility.GenerateUUID(),
		Email:     userEmail,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(time.Duration(config.App.ResetPasswordDuration) * time.Minute),
	}

	err = reset.CreatePasswordReset(db)
	if err != nil {
		return "error", http.StatusInternalServerError, err
	}

	// Send email with the reset link (e.g., http://example.com/reset-password?token=resetToken)
	//SendBackgroundEmail(user.Email, resetToken, "reset password")

	return "success", http.StatusOK, nil
}

func VerifyPasswordResetToken(req models.ResetPasswordRequestModel, db *gorm.DB) (*models.User, int, error) {

	var (
		user      = models.User{}
		passReset = models.PasswordReset{}
	)

	resetExist, err := passReset.GetPasswordResetByToken(db, req.Token)
	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("invalid or expired token")
	}

	userDataExist, err := user.GetUserByEmail(db, resetExist.Email)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	hashedPassword, err := utility.HashPassword(req.NewPassword)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	userDataExist.Password = hashedPassword
	err = userDataExist.Update(db)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := resetExist.DeletePasswordReset(db); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Send email with the reset link (e.g., http://example.com/reset-password?token=resetToken)
	//SendBackgroundEmail(user.Email, "reset password verify")

	return &userDataExist, http.StatusOK, nil

}
