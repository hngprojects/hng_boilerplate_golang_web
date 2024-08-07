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

func TestAddToTimezone(t *testing.T) {
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

	setup := func() (*gin.Engine, *auth.Controller) {
		router, saController := SetupSATestRouter()
		authController := auth.Controller{
			Db:        saController.Db,
			Validator: saController.Validator,
			Logger:    saController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Add to Timezone", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		timezone := models.Timezone{
			ID:          utility.GenerateUUID(),
			Timezone:    fmt.Sprintf("America/New_York-%s", utility.RandomString(10)),
			GmtOffset:   fmt.Sprintf("-05:00+%s", utility.RandomString(3)),
			Description: fmt.Sprintf("western -%s", utility.RandomString(3)),
		}

		jsonBody, _ := json.Marshal(timezone)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/timezones", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusCreated)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Timezone added successfully")
	})

	t.Run("Validation Error - Missing Fields", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		timezone := models.Timezone{
			ID:          utility.GenerateUUID(),
			Timezone:    "",
			GmtOffset:   fmt.Sprintf("-05:00+%s", utility.RandomString(5)),
			Description: fmt.Sprintf("western -%s", utility.RandomString(3)),
		}
		jsonBody, _ := json.Marshal(timezone)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/timezones", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})

	t.Run("Duplicate Timezone", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)
		theRandom := utility.RandomString(6)

		timezone := models.Timezone{
			Timezone:    fmt.Sprintf("UTC-%s", theRandom),
			GmtOffset:   fmt.Sprintf("-05:00+%s", utility.RandomString(5)),
			Description: fmt.Sprintf("western -%s", utility.RandomString(3)),
		}
		authController.Db.Postgresql.Create(&timezone)

		duplicateTimezone := models.Timezone{
			Timezone:    fmt.Sprintf("UTC-%s", theRandom),
			GmtOffset:   fmt.Sprintf("-05:00+%s", utility.RandomString(5)),
			Description: fmt.Sprintf("western -%s", utility.RandomString(3)),
		}
		jsonBody, _ := json.Marshal(duplicateTimezone)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/timezones", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusBadRequest)
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		router, _ := setup()

		timezone := models.Timezone{
			ID:          utility.GenerateUUID(),
			Timezone:    fmt.Sprintf("America/New_York-%s", utility.RandomString(18)),
			GmtOffset:   fmt.Sprintf("-05:00+%s", utility.RandomString(4)),
			Description: fmt.Sprintf("western -%s", utility.RandomString(3)),
		}
		jsonBody, _ := json.Marshal(timezone)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/timezones", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token is invalid!")
	})
}

func TestGetTimezones(t *testing.T) {
	setup := func() (*gin.Engine, *auth.Controller) {
		router, saController := SetupSATestRouter()
		authController := auth.Controller{
			Db:        saController.Db,
			Validator: saController.Validator,
			Logger:    saController.Logger,
		}

		return router, &authController
	}

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
	regularUser := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Regular User",
		Email:    fmt.Sprintf("user%v@qa.team", currUUID),
		Password: password,
		Role:     int(models.RoleIdentity.User),
	}

	timezone := models.Timezone{
		ID:          utility.GenerateUUID(),
		Timezone:    fmt.Sprintf("America/New_York-%s", utility.RandomString(10)),
		GmtOffset:   fmt.Sprintf("-05:00+%s", utility.RandomString(5)),
		Description: fmt.Sprintf("western -%s", utility.RandomString(3)),
	}

	db.Create(&adminUser)
	db.Create(&regularUser)
	db.Create(&timezone)

	t.Run("Successful Get Timezones for admin", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/timezones", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Timezones retrieved successfully")
	})

	t.Run("Successful Get Timezones for user", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/timezones", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Timezones retrieved successfully")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		router, _ := setup()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/timezones", nil)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token could not be found!")
	})
}
