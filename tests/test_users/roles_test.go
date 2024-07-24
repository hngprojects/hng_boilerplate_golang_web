package test_users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestUpdateUserRole(t *testing.T) {
	router, userController := SetupUsersTestRouter()
	db := userController.Db.Postgresql
	currUUID := utility.GenerateUUID()
	theRole := models.RoleIdentity.SuperAdmin
	userRole := models.RoleIdentity.User
	password, _ := utility.HashPassword("password")

	adminData := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "admin jane doe",
		Email:    fmt.Sprintf("testadmin%v@qa.team", currUUID),
		Password: password,
		Role:     int(theRole),
	}
	db.Create(&adminData)

	userData := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "user jane doe",
		Email:    fmt.Sprintf("testuser%v@qa.team", currUUID),
		Password: password,
		Role:     int(userRole),
	}
	db.Create(&userData)

	loginData := models.LoginRequestModel{
		Email:    adminData.Email,
		Password: "password",
	}

	auth := auth.Controller{Db: userController.Db, Validator: userController.Validator, Logger: userController.Logger}
	token := tests.GetLoginToken(t, router, auth, loginData)

	t.Run("Successful Update", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/roles/%d", userData.ID, theRole), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Role updated successfully")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/roles/%d", userData.ID, theRole), nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token could not be found!")
		tests.AssertResponseMessage(t, response["error"].(string), "Unauthorized")
	})

	t.Run("User Not Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/roles/%d", currUUID, theRole), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusNotFound)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "invalid user")
	})

	t.Run("Role Not Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/roles/%d", userData.ID, 999), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusNotFound)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "invalid role")
	})

}
