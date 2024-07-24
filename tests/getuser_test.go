package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)



func TestGetUserById(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}
	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	userController := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, userController, userSignUpData)

	token := GetLoginToken(t, r, userController, loginData)
	userID := currUUID

	tests := []struct {
		Name         string
		UserId       string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Invalid User ID Format",
			UserId:       "invalid-id-erttt",
			ExpectedCode: http.StatusBadRequest,
			Message:      "invalid user id format",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Successful retrieval of user info",
			UserId:       userID,
			ExpectedCode: http.StatusOK,
			Message:      "User retrieved successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "User Not Found",
			UserId:       utility.GenerateUUID(),
			ExpectedCode: http.StatusNotFound,
			Message:      "user not found",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "User Not Authorized",
			UserId:       userID,
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	orgUrl := r.Group("/api/v1", middleware.Authorize())
	{
		orgUrl.GET("/users/:userid", userController.GetUserByID)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", test.UserId), nil)
			if err != nil {
				t.Fatal(err)
			}

			for key, value := range test.Headers {
				req.Header.Set(key, value)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.ExpectedCode {
				t.Errorf("Expected status code %d, got %d", test.ExpectedCode, rr.Code)
			}

			if rr.Code == http.StatusNoContent {
				// 204 No Content, no body to check
				return
			}

			var data map[string]interface{}
			if err := json.NewDecoder(rr.Body).Decode(&data); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			code, ok := data["status_code"].(float64)
			if !ok {
				t.Fatalf("Expected status_code to be float64, got %T", data["status_code"])
			}
			if int(code) != test.ExpectedCode {
				t.Errorf("Expected status code %d, got %d", test.ExpectedCode, int(code))
			}

			if test.Message != "" {
				message, ok := data["message"].(string)
				if !ok {
					t.Fatalf("Expected message to be string, got %T", data["message"])
				}
				if message != test.Message {
					t.Errorf("Expected message '%s', got '%s'", test.Message, message)
				}
			}
		})
	}
}
