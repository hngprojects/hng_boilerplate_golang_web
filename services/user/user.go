package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func ValidateCreateUserRequest(req models.CreateUserRequestModel, db *gorm.DB) (models.CreateUserRequestModel, error) {

	user := models.User{}
	profile := models.Profile{}

	// Check if the user email is valid or already exists

	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)
		formattedMail, checkBool := utility.EmailValid(req.Email)
		if !checkBool {
			return req, fmt.Errorf("email address is invalid")
		}
		req.Email = formattedMail
		exists := postgresql.CheckExists(db, &user, "email = ?", req.Email)
		if exists {
			return req, errors.New("user already exists with the given email")
		}
	}

	// Check if the user phone is valid, then format and check if already exists

	if req.PhoneNumber != "" {
		req.PhoneNumber = strings.ToLower(req.PhoneNumber)
		phone, _ := utility.PhoneValid(req.PhoneNumber)
		req.PhoneNumber = phone
		exists := postgresql.CheckExists(db, &profile, "phone = ?", req.PhoneNumber)
		if exists {
			return req, errors.New("user already exists with the given phone")
		}

	}

	return req, nil
}

func GetUser(userIDStr string, db *gorm.DB) (models.User, error) {
	var userResp models.User

	userResp, err := userResp.GetUserByID(db, userIDStr)
	if err != nil {
		return userResp, err
	}

	return userResp, nil
}

func CreateUser(req models.CreateUserRequestModel, db *gorm.DB) (gin.H, int, error) {

	var (
		email        = strings.ToLower(req.Email)
		firstName    = strings.Title(strings.ToLower(req.FirstName))
		lastName     = strings.Title(strings.ToLower(req.LastName))
		username     = strings.ToLower(req.UserName)
		phoneNumber  = req.PhoneNumber
		password     = req.Password
		responseData gin.H
	)

	password, err := utility.HashPassword(req.Password)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	user := models.User{
		ID:       utility.GenerateUUID(),
		Name:     username,
		Email:    email,
		Password: password,
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phoneNumber,
		},
	}

	err = user.CreateUser(db)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	token, expiry, err := middleware.CreateToken(user)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	responseData = gin.H{
		"email":        user.Email,
		"username":     user.Name,
		"first_name":   user.Profile.FirstName,
		"last_name":    user.Profile.LastName,
		"phone":        user.Profile.Phone,
		"expires_in":   expiry,
		"access_token": token,
	}

	return responseData, http.StatusCreated, nil
}

func LoginUser(req models.LoginRequestModel, db *gorm.DB) (gin.H, int, error) {

	var (
		user         = models.User{}
		responseData gin.H
	)

	// Check if the user email exists
	exists := postgresql.CheckExists(db, &user, "email = ?", req.Email)
	if !exists {
		return responseData, 400, fmt.Errorf("invalid credentials")
	}

	if !utility.CompareHash(req.Password, user.Password) {
		return responseData, 400, fmt.Errorf("invalid credentials")
	}

	userData, err := user.GetUserByID(db, user.ID)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("unable to fetch user " + err.Error())
	}

	token, expiry, err := middleware.CreateToken(userData)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	responseData = gin.H{
		"email":        userData.Email,
		"username":     userData.Name,
		"first_name":   userData.Profile.FirstName,
		"last_name":    userData.Profile.LastName,
		"phone":        userData.Profile.Phone,
		"expires_in":   expiry,
		"access_token": token,
	}

	return responseData, http.StatusCreated, nil
}

func GetUserByEmail(email string, db *gorm.DB)(models.User, error){
	var user models.User

	user, err := user.GetUserByEmail(db, email);
	if err != nil {
		return user, err
	}
	return user, nil
}