package test_superadmin

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
)

func TestUpdateTimezone(t *testing.T) {
	_, saController := SetupSATestRouter()
	db := saController.Db.Postgresql

	currUUID := utility.GenerateUUID()
	password, _ := utility.HashPassword("password")

	adminUser := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Admin User",
		Email:    fmt.Sprintf("admin%v@qa.team", currUUID),
		Password: password,
		Role:     int(models.RoleIdentity.SuperAdmin),
	}

	db.Create(&adminUser)

	timezone := models.Timezone{
		ID:          utility.GenerateUUID(),
		Timezone:    fmt.Sprintf("America/New_York-%s", utility.RandomString(10)),
		GmtOffset:   fmt.Sprintf("-05:00+%s", utility.RandomString(3)),
		Description: fmt.Sprintf("western -%s", utility.RandomString(3)),
	}
	db.Create(&timezone)

	setup := func() (*gin.Engine, *auth.Controller) {
		router, saController := SetupSATestRouter()
		authController := auth.Controller{
			Db:        saController.Db,
			Validator: saController.Validator,
			Logger:    saController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Update Timezone", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateTimezone := models.Timezone{
			Timezone:    fmt.Sprintf("America/Los_Angeles-%s", utility.RandomString(10)),
			GmtOffset:   fmt.Sprintf("-08:00+%s", utility.RandomString(3)),
			Description: fmt.Sprintf("pacific -%s", utility.RandomString(3)),
		}
		jsonBody, _ := json.Marshal(updateTimezone)

		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/timezones/%s", timezone.ID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Timezone updated successfully")
	})

	t.Run("Timezone Not Found", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateTimezone := models.Timezone{
			Timezone:    fmt.Sprintf("America/Los_Angeles-%s", utility.RandomString(10)),
			GmtOffset:   fmt.Sprintf("-08:00+%s", utility.RandomString(3)),
			Description: fmt.Sprintf("pacific -%s", utility.RandomString(3)),
		}
		jsonBody, _ := json.Marshal(updateTimezone)

		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/timezones/%s", utility.GenerateUUID()), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusNotFound)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "record not found")
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		router, _ := setup()

		updateTimezone := models.Timezone{
			Timezone:    fmt.Sprintf("America/Los_Angeles-%s", utility.RandomString(18)),
			GmtOffset:   fmt.Sprintf("-08:00+%s", utility.RandomString(3)),
			Description: fmt.Sprintf("pacific -%s", utility.RandomString(3)),
		}
		jsonBody, _ := json.Marshal(updateTimezone)

		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/timezones/%s", timezone.ID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token is invalid!")
	})
}
