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

func TestUpdateUserDataPrivacy(t *testing.T) {

	_, userController := SetupUsersTestRouter()
	db := userController.Db.Postgresql
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

	dataPrivacy := models.DataPrivacySettings{
		ID:                    utility.GenerateUUID(),
		UserID:                regularUser.ID,
		ProfileVisibility:     false,
		ShareDataWithPartners: false,
		ReceiveEmailUpdates:   false,
		Enable2FA:             true,
		UseDataEncryption:     false,
		AllowAnalytics:        true,
		PersonalizedAds:       true,
	}

	db.Create(&adminUser)
	db.Create(&regularUser)
	db.Create(&dataPrivacy)

	setup := func() (*gin.Engine, *auth.Controller) {
		router, userController := SetupUsersTestRouter()
		authController := auth.Controller{
			Db:        userController.Db,
			Validator: userController.Validator,
			Logger:    userController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Update User Data Privacy for admin", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.DataPrivacySettings{
			ProfileVisibility:     true,
			ShareDataWithPartners: true,
			ReceiveEmailUpdates:   true,
			Enable2FA:             false,
			UseDataEncryption:     true,
			AllowAnalytics:        false,
			PersonalizedAds:       false,
		}
		body, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/data-privacy-settings", regularUser.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["status"].(string), "success")
		tests.AssertResponseMessage(t, response["message"].(string), "User data privacy settings updated successfully")
	})

	t.Run("Successful Update User Data Privacy for regular user", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.DataPrivacySettings{
			ProfileVisibility:     true,
			ShareDataWithPartners: true,
			ReceiveEmailUpdates:   true,
			Enable2FA:             false,
			UseDataEncryption:     true,
			AllowAnalytics:        false,
			PersonalizedAds:       false,
		}
		body, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/data-privacy-settings", regularUser.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["status"].(string), "success")
		tests.AssertResponseMessage(t, response["message"].(string), "User data privacy settings updated successfully")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		router, _ := setup()

		updateData := models.DataPrivacySettings{
			ProfileVisibility:     true,
			ShareDataWithPartners: true,
			ReceiveEmailUpdates:   true,
			Enable2FA:             false,
			UseDataEncryption:     true,
			AllowAnalytics:        false,
			PersonalizedAds:       false,
		}
		body, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/data-privacy-settings", regularUser.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token could not be found!")
	})

	t.Run("User ID not found", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateData := models.DataPrivacySettings{
			ProfileVisibility:     true,
			ShareDataWithPartners: true,
			ReceiveEmailUpdates:   true,
			Enable2FA:             false,
			UseDataEncryption:     true,
			AllowAnalytics:        false,
			PersonalizedAds:       false,
		}
		body, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s/data-privacy-settings", utility.GenerateUUID()), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusNotFound)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["status"].(string), "error")
		tests.AssertResponseMessage(t, response["message"].(string), "user not found")
	})
}
