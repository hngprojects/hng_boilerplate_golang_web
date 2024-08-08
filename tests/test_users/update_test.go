package test_users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAUser(t *testing.T) {
	_, userController := SetupUsersTestRouter()
	db := userController.Db.Postgresql
	currUUID := utility.GenerateUUID()
	password, _ := utility.HashPassword("password")

	adminUser := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Admin User",
		Email:    fmt.Sprintf("admin%v@qa.team", currUUID),
		Password: password,
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: "Admin",
			LastName:  "Doe user",
			Phone:     "09876543211",
			AvatarURL: "http://example.com/avatar2.jpg",
		},
		Role: int(models.RoleIdentity.SuperAdmin),
	}
	regularUser := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Regular User",
		Email:    fmt.Sprintf("user%v@qa.team", currUUID),
		Password: password,
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: "Regualr",
			LastName:  "User",
			Phone:     "0987654321",
			AvatarURL: "http://example.com/avatar2.jpg",
		},
		Role: int(models.RoleIdentity.User),
	}

	db.Create(&adminUser)
	db.Create(&regularUser)

	// This function ensures a fresh setup for each test case
	setup := func() (*gin.Engine, *auth.Controller) {
		router, userController := SetupUsersTestRouter()
		authController := auth.Controller{
			Db:        userController.Db,
			Validator: userController.Validator,
			Logger:    userController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Update User by Super Admin", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.UpdateUserRequestModel{
			FirstName:   "UpdatedFirstName",
			LastName:    "UpdatedLastName",
			UserName:    "UpdatedUserName",
			PhoneNumber: "1234567890",
		}
		updateDataJSON, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", regularUser.ID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "User info updated successfully")
	})

	t.Run("Update user who doesn't exist", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.UpdateUserRequestModel{
			FirstName:   "UpdatedFirstName",
			LastName:    "UpdatedLastName",
			UserName:    "UpdatedUserName",
			PhoneNumber: "1234567890",
		}
		updateDataJSON, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", utility.GenerateUUID()), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "user not found")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		router, _ := setup()

		updateData := models.UpdateUserRequestModel{
			FirstName:   "UpdatedFirstName",
			LastName:    "UpdatedLastName",
			UserName:    "UpdatedUserName",
			PhoneNumber: "1234567890",
		}
		updateDataJSON, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", regularUser.ID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token could not be found!")
	})

	t.Run("Field Validation", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.UpdateUserRequestModel{
			FirstName:   "UpdatedFirstName",
			LastName:    "UpdatedLastName",
			PhoneNumber: "1234567890",
		}
		updateDataJSON, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", regularUser.ID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})

	t.Run("Unauthorized Update User by Regular User", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.UpdateUserRequestModel{
			FirstName:   "UpdatedFirstName",
			LastName:    "UpdatedLastName",
			UserName:    "UpdatedUserName",
			PhoneNumber: "1234567890",
		}
		updateDataJSON, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", adminUser.ID), bytes.NewBuffer(updateDataJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusForbidden, resp.Code)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "user does not have permission to update this user")
	})
}
