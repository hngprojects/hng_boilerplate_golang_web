package test_users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := Setup()
	validate := validator.New()
	db := storage.Connection()

	// Setup user data for the test
	userData := models.CreateUserRequestModel{
		Email:       "test@example.com",
		Password:    "password",
		FirstName:   "First",
		LastName:    "Last",
		UserName:    "username",
		PhoneNumber: "1234567890",
	}

	SignupUser(t, gin.Default(), user.Controller{Db: db, Logger: logger, Validator: validate}, userData)

	loginData := models.LoginRequestModel{
		Email:    userData.Email,
		Password: userData.Password,
	}

	token := GetLoginToken(t, gin.Default(), user.Controller{Db: db, Logger: logger, Validator: validate}, loginData)

	tests := []struct {
		Name            string
		UserID          string
		Request         interface{}
		ExpectedCode    int
		ExpectedMessage string
	}{
		{
			Name:   "Successful response with valid id",
			UserID: utility.GenerateUUID(),
			Request: models.UpdateUserRequestModel{
				Name:        "Updated Name",
				PhoneNumber: "0987654321",
			},
			ExpectedCode:    http.StatusOK,
			ExpectedMessage: "User updated successfully",
		},
		{
			Name:   "Invalid userId",
			UserID: "invalid-id",
			Request: models.UpdateUserRequestModel{
				Name:        "Updated Name",
				PhoneNumber: "0987654321",
			},
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "Failed to parse request body",
		},
		{
			Name:   "Missing userId",
			UserID: "",
			Request: models.UpdateUserRequestModel{
				Name:        "Updated Name",
				PhoneNumber: "0987654321",
			},
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "Failed to parse request body",
		},
		{
			Name:   "Invalid request body",
			UserID: utility.GenerateUUID(),
			Request: map[string]interface{}{
				"name":        123, // Invalid type
				"phoneNumber": "0987654321",
			},
			ExpectedCode:    http.StatusUnprocessableEntity,
			ExpectedMessage: "Validation failed",
		},
	}

	uc := user.Controller{Db: db, Logger: logger, Validator: validate}

	for _, tt := range tests {
		r := gin.Default()
		r.PUT("/api/v1/users/:userId", uc.UpdateUser)

		t.Run(tt.Name, func(t *testing.T) {
			var buf bytes.Buffer

			err := json.NewEncoder(&buf).Encode(tt.Request)
			if err != nil {
				t.Fatal(err)
			}

			url := "/api/v1/users/" + tt.UserID
			req, err := http.NewRequest(http.MethodPut, url, &buf)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			hr := httptest.NewRecorder()
			r.ServeHTTP(hr, req)

			AssertStatusCode(t, hr.Code, tt.ExpectedCode)

			data := base.ParseResponse(hr)

			if tt.ExpectedMessage != "" {
				message := data["message"]
				if message != nil {
					AssertResponseMessage(t, message.(string), tt.ExpectedMessage)
				} else {
					AssertResponseMessage(t, "", tt.ExpectedMessage)
				}
			}
		})
	}
}
