package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func CreateProviderUser(req models.CreateUserRequestModel, db *gorm.DB) (gin.H, int, error) {

	var (
		email        = strings.ToLower(req.Email)
		firstName    = strings.Title(strings.ToLower(req.FirstName))
		lastName     = strings.Title(strings.ToLower(req.LastName))
		username     = strings.ToLower(req.UserName)
		phoneNumber  = req.PhoneNumber
		password     = req.Password
		responseData gin.H
		user         models.User
	)

	// check if user already exists
	req, err := ValidateCreateUserRequest(req, db)
	if err != nil && errors.Is(err, errors.New("user already exists with the given email")) {

		exists := postgresql.CheckExists(db, &user, "email = ?", email)
		if !exists {
			return responseData, http.StatusNotFound, fmt.Errorf("user not found")
		}

	} else {

		user = models.User{
			ID:       utility.GenerateUUID(),
			Name:     username,
			Email:    email,
			Password: password,
			Role:     int(models.RoleIdentity.User),
			Profile: models.Profile{
				ID:        utility.GenerateUUID(),
				FirstName: firstName,
				LastName:  lastName,
				Phone:     phoneNumber,
			},
		}
		err := user.CreateUser(db)
		if err != nil {
			return responseData, http.StatusInternalServerError, err
		}
	}

	tokenData, err := middleware.CreateToken(user)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	tokens := map[string]string{
		"access_token": tokenData.AccessToken,
		"exp":          strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
	}

	access_token := models.AccessToken{ID: tokenData.AccessUuid, OwnerID: user.ID}

	err = access_token.CreateAccessToken(db, tokens)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	responseData = gin.H{
		"email":        user.Email,
		"username":     user.Name,
		"first_name":   user.Profile.FirstName,
		"last_name":    user.Profile.LastName,
		"phone":        user.Profile.Phone,
		"role":         models.UserRoleName,
		"expires_in":   tokenData.ExpiresAt.Unix(),
		"access_token": tokenData.AccessToken,
	}

	return responseData, http.StatusCreated, nil
}
