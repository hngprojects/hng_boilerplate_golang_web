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
)

func TestUpdateUserRegion(t *testing.T) {
	setup := func() (*gin.Engine, *auth.Controller) {
		router, saController := SetupUsersTestRouter()
		authController := auth.Controller{
			Db:        saController.Db,
			Validator: saController.Validator,
			Logger:    saController.Logger,
		}

		return router, &authController
	}

	_, saController := SetupUsersTestRouter()
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

	region := models.Region{
		ID:   utility.GenerateUUID(),
		Name: fmt.Sprintf("Region-%s", utility.RandomString(10)),
		Code: fmt.Sprintf("RG-%s", utility.RandomString(5)),
	}

	timezone := models.Timezone{
		ID:          utility.GenerateUUID(),
		Timezone:    fmt.Sprintf("Timezone-%s", utility.RandomString(10)),
		GmtOffset:   fmt.Sprintf("UTC+%s", utility.RandomString(3)),
		Description: fmt.Sprintf("western -%s", utility.RandomString(3)),
	}

	language := models.Language{
		ID:   utility.GenerateUUID(),
		Name: fmt.Sprintf("English-%s", utility.RandomString(10)),
		Code: fmt.Sprintf("en-%s", utility.RandomString(5)),
	}

	db.Create(&adminUser)
	db.Create(&regularUser)
	db.Create(&region)
	db.Create(&timezone)
	db.Create(&language)

	t.Run("Successful Update User Region by Admin", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.UserRegionTimezoneLanguage{
			RegionID:   region.ID,
			TimezoneID: timezone.ID,
			LanguageID: language.ID,
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/regions", regularUser.ID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "User info updated successfully")
	})

	t.Run("Validation Failed", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.UserRegionTimezoneLanguage{
			RegionID:   utility.GenerateUUID(),
			TimezoneID: "",
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/regions", regularUser.ID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})

	t.Run("Wrong User ID", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.UserRegionTimezoneLanguage{
			RegionID:   region.ID,
			TimezoneID: timezone.ID,
			LanguageID: language.ID,
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/regions", utility.GenerateUUID()), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusNotFound)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "user not found")
	})

	t.Run("Normal User Trying to Access", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.UserRegionTimezoneLanguage{
			RegionID:   region.ID,
			TimezoneID: timezone.ID,
			LanguageID: language.ID,
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/regions", adminUser.ID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "role not authorized!")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		router, _ := setup()

		updateData := models.UserRegionTimezoneLanguage{
			RegionID:   region.ID,
			TimezoneID: timezone.ID,
			LanguageID: language.ID,
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/regions", regularUser.ID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token could not be found!")
	})
}
