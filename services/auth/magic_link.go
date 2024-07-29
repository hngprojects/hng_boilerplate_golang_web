package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func MagicLinkRequest(userEmail string, db *gorm.DB) (string, int, error) {

	var (
		user      = models.User{}
		magicLink = models.MagicLink{}
		config    = config.GetConfig()
	)

	magicExist, err := magicLink.GetMagicLinkByEmail(db, userEmail)
	if err != nil {
		return "error", http.StatusUnauthorized, err
	}

	if magicExist != nil {
		if err := magicExist.DeleteMagicLink(db); err != nil {
			return "error", http.StatusInternalServerError, err
		}
	}

	exists := postgresql.CheckExists(db, &user, "email = ?", userEmail)
	if !exists {
		return "error", http.StatusNotFound, fmt.Errorf("user not found")
	}

	requestToken := utility.GenerateUUID()
	magic := models.MagicLink{
		ID:        utility.GenerateUUID(),
		Email:     userEmail,
		Token:     requestToken,
		ExpiresAt: time.Now().Add(time.Duration(config.App.MagicLinkDuration) * time.Minute),
	}

	err = magic.CreateMagicLink(db)
	if err != nil {
		return "error", http.StatusInternalServerError, err
	}

	// Send email with the reset link (e.g., http://example.com/reset-password?token=resetToken)
	//SendBackgroundEmail(user.Email, resetToken, "magic link")

	return "success", http.StatusOK, nil
}

func VerifyMagicLinkToken(req models.VerifyMagicLinkRequest, db *gorm.DB) (gin.H, int, error) {

	var (
		user         = models.User{}
		responseData gin.H
		magicLink    = models.MagicLink{}
	)

	magicExist, err := magicLink.GetMagicLinkByToken(db, req.Token)
	if err != nil {
		return responseData, http.StatusUnauthorized, errors.New("invalid or expired token")
	}

	exists := postgresql.CheckExists(db, &user, "email = ?", magicExist.Email)
	if !exists {
		return responseData, http.StatusBadRequest, errors.New("invalid credentials")
	}

	userData, err := user.GetUserByEmail(db, magicExist.Email)
	if err != nil {
		return responseData, http.StatusInternalServerError, errors.New("unable to fetch user")
	}

	tokenData, err := middleware.CreateToken(user)
	if err != nil {
		return responseData, http.StatusInternalServerError, errors.New("error saving token")
	}

	tokens := map[string]string{
		"access_token": tokenData.AccessToken,
		"exp":          strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
	}

	access_token := models.AccessToken{ID: tokenData.AccessUuid, OwnerID: user.ID}

	err = access_token.CreateAccessToken(db, tokens)

	if err != nil {
		return responseData, http.StatusInternalServerError, errors.New("error saving token")
	}

	if err := magicExist.DeleteMagicLink(db); err != nil {
		return responseData, http.StatusInternalServerError, err
	}

	responseData = gin.H{
		"email":        userData.Email,
		"username":     userData.Name,
		"first_name":   userData.Profile.FirstName,
		"last_name":    userData.Profile.LastName,
		"phone":        userData.Profile.Phone,
		"role":         userData.Role,
		"expires_in":   tokenData.ExpiresAt.Unix(),
		"access_token": tokenData.AccessToken,
		"created_at":   userData.CreatedAt,
		"updated_at":   userData.UpdatedAt,
	}

	return responseData, http.StatusOK, nil
}
