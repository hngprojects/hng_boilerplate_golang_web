package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/actions"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/actions/names"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type AuthService interface {
	ValidateCreateUserRequest(req models.CreateUserRequestModel) (models.CreateUserRequestModel, error)
	GetUser(userIDStr string) (models.User, error)
	CreateUser(req models.CreateUserRequestModel) (gin.H, int, error)
	CreateAdmin(req models.CreateUserRequestModel) (gin.H, int, error)
	LoginUser(req models.LoginRequestModel) (gin.H, int, error)
	LogoutUser(access_uuid, owner_id string) (gin.H, int, error)
}

// authService implements AuthService.
type authService struct {
	db database.DatabaseManager
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db: inst.InitDB(db)}
}

func (s *authService)ValidateCreateUserRequest(req models.CreateUserRequestModel) (models.CreateUserRequestModel, error) {
	// instance of Postgresql db
	
	user := models.User{}
	profile := models.Profile{}

	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)
		formattedMail, checkBool := utility.EmailValid(req.Email)
		if !checkBool {
			return req, fmt.Errorf("email address is invalid")
		}
		req.Email = formattedMail
		exists := s.db.CheckExists(&user, "email = ?", req.Email)
		if exists {
			return req, errors.New("user already exists with the given email")
		}
	}

	if req.PhoneNumber != "" {
		req.PhoneNumber = strings.ToLower(req.PhoneNumber)
		phone, _ := utility.PhoneValid(req.PhoneNumber)
		req.PhoneNumber = phone
		exists := s.db.CheckExists(&profile, "phone = ?", req.PhoneNumber)
		if exists {
			return req, errors.New("user already exists with the given phone")
		}

	}

	return req, nil
}

func (s *authService)GetUser(userIDStr string) (models.User, error) {
	// instance of Postgresql db
	
	var userResp models.User
	userResp, err := userResp.GetUserByID(s.db, userIDStr)
	if err != nil {
		return userResp, err
	}

	return userResp, nil
}

func (s *authService)CreateUser(req models.CreateUserRequestModel) (gin.H, int, error) {

	// instance of Postgresql db
	
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
		Role:     int(models.RoleIdentity.User),
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phoneNumber,
		},
	}

	err = user.CreateUser(s.db)
	if err != nil {
		return nil, http.StatusInternalServerError, err
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

	err = access_token.CreateAccessToken(s.db, tokens)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	responseData = gin.H{
		"user": map[string]string{
			"id":         user.ID,
			"email":      user.Email,
			"username":   user.Name,
			"first_name": user.Profile.FirstName,
			"last_name":  user.Profile.LastName,
			"fullname":   user.Profile.FirstName + " " + user.Profile.LastName,
			"phone":      user.Profile.Phone,
			"role":       strconv.Itoa(user.Role),
			"expires_in": strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
			"created_at": strconv.Itoa(int(user.CreatedAt.Unix())),
			"updated_at": strconv.Itoa(int(user.UpdatedAt.Unix())),
		},
		"access_token": tokenData.AccessToken,
	}

	resetReq := models.SendWelcomeMail{
		Email: user.Email,
	}

	err = actions.AddNotificationToQueue(storage.DB.Redis.RedisDb(), names.SendWelcomeMail, resetReq)
	if err != nil {
		return responseData, http.StatusInternalServerError, err
	}

	return responseData, http.StatusCreated, nil
}

func (s *authService)CreateAdmin(req models.CreateUserRequestModel) (gin.H, int, error) {

	// instance of Postgresql db
	
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
		Role:     int(models.RoleIdentity.SuperAdmin),
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phoneNumber,
		},
	}

	err = user.CreateUser(s.db)
	if err != nil {
		return nil, http.StatusInternalServerError, err
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

	err = access_token.CreateAccessToken(s.db, tokens)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	responseData = gin.H{
		"user": map[string]string{
			"id":         user.ID,
			"email":      user.Email,
			"username":   user.Name,
			"first_name": user.Profile.FirstName,
			"last_name":  user.Profile.LastName,
			"fullname":   user.Profile.FirstName + " " + user.Profile.LastName,
			"phone":      user.Profile.Phone,
			"role":       strconv.Itoa(user.Role),
			"expires_in": strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
			"created_at": strconv.Itoa(int(user.CreatedAt.Unix())),
			"updated_at": strconv.Itoa(int(user.UpdatedAt.Unix())),
		},
		"access_token": tokenData.AccessToken,
	}

	return responseData, http.StatusCreated, nil
}

func (s *authService)LoginUser(req models.LoginRequestModel) (gin.H, int, error) {

	// instance of Postgresql db
	
	var (
		user         = models.User{}
		responseData gin.H
	)

	// Check if the user email exists
	exists := s.db.CheckExists(&user, "email = ?", req.Email)
	if !exists {
		return responseData, 400, fmt.Errorf("invalid credentials")
	}

	if !utility.CompareHash(req.Password, user.Password) {
		return responseData, 400, fmt.Errorf("invalid credentials")
	}

	userData, err := user.GetUserByID(s.db, user.ID)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("unable to fetch user " + err.Error())
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

	err = access_token.CreateAccessToken(s.db, tokens)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	responseData = gin.H{

		"user": map[string]string{
			"id":         userData.ID,
			"email":      userData.Email,
			"username":   userData.Name,
			"first_name": userData.Profile.FirstName,
			"last_name":  userData.Profile.LastName,
			"fullname":   userData.Profile.FirstName + " " + userData.Profile.LastName,
			"phone":      userData.Profile.Phone,
			"role":       strconv.Itoa(userData.Role),
			"expires_in": strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
			"created_at": strconv.Itoa(int(userData.CreatedAt.Unix())),
			"updated_at": strconv.Itoa(int(userData.UpdatedAt.Unix())),
		},
		"access_token": tokenData.AccessToken,
	}

	return responseData, http.StatusOK, nil
}

func (s *authService)LogoutUser(access_uuid, owner_id string) (gin.H, int, error) {

	// instance of Postgresql db
	
	var (
		responseData gin.H
	)

	access_token := models.AccessToken{ID: access_uuid, OwnerID: owner_id}

	// revoke user access_token to invalidate session
	err := access_token.RevokeAccessToken(s.db)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error revoking user session: " + err.Error())
	}

	responseData = gin.H{}

	return responseData, http.StatusOK, nil
}
