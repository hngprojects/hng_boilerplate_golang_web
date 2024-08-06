package test_contact

import (
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

func TestGetAllContactUs(t *testing.T) {
	_, contactController := SetupContactTestRouter()
	db := contactController.Db.Postgresql
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

	db.Create(&adminUser)
	db.Create(&regularUser)

	contact1 := models.ContactUs{
		ID:      utility.GenerateUUID(),
		Email:   fmt.Sprintf("contact%v@qa.team", currUUID),
		Subject: fmt.Sprintf("First subject - %v ", currUUID),
		Message: fmt.Sprintf("First test message - %v ", currUUID),
	}
	contact2 := models.ContactUs{
		ID:      utility.GenerateUUID(),
		Email:   fmt.Sprintf("contact2%v@qa.team", currUUID),
		Subject: fmt.Sprintf("Second subject - %v ", currUUID),
		Message: fmt.Sprintf("Second test message - %v ", currUUID),
	}

	db.Create(&contact1)
	db.Create(&contact2)

	setup := func() (*gin.Engine, *auth.Controller) {
		router, contactController := SetupContactTestRouter()
		authController := auth.Controller{
			Db:        contactController.Db,
			Validator: contactController.Validator,
			Logger:    contactController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Get All Contact Us", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/contact", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Messages retrieved successfully")
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		router, _ := setup()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/contact", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token is invalid!")
	})

	t.Run("Forbidden Access - Regular User Trying to Get All", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/contact", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "role not authorized!")
	})
}

func TestGetContactUsById(t *testing.T) {
	_, contactController := SetupContactTestRouter()
	db := contactController.Db.Postgresql
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

	db.Create(&adminUser)
	db.Create(&regularUser)

	contact := models.ContactUs{
		ID:      utility.GenerateUUID(),
		Email:   fmt.Sprintf("contact%v@qa.team", currUUID),
		Subject: fmt.Sprintf("Test subject - %v ", currUUID),
		Message: fmt.Sprintf("Test message - %v ", currUUID),
	}

	db.Create(&contact)

	setup := func() (*gin.Engine, *auth.Controller) {
		router, contactController := SetupContactTestRouter()
		authController := auth.Controller{
			Db:        contactController.Db,
			Validator: contactController.Validator,
			Logger:    contactController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Get Contact Us By ID", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/contact/id/%s", contact.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Message retrieved successfully")
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		router, _ := setup()

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/contact/id/%s", contact.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token is invalid!")
	})

	t.Run("Forbidden Access - Regular User Trying to Get By ID", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/contact/id/%s", contact.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "role not authorized!")
	})
}

func TestGetContactUsByEmail(t *testing.T) {
	_, contactController := SetupContactTestRouter()
	db := contactController.Db.Postgresql
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

	db.Create(&adminUser)
	db.Create(&regularUser)

	contact := models.ContactUs{
		ID:      utility.GenerateUUID(),
		Email:   fmt.Sprintf("contact%v@qa.team", currUUID),
		Subject: fmt.Sprintf("Test subject - %v ", currUUID),
		Message: fmt.Sprintf("Test message - %v ", currUUID),
	}

	db.Create(&contact)

	setup := func() (*gin.Engine, *auth.Controller) {
		router, contactController := SetupContactTestRouter()
		authController := auth.Controller{
			Db:        contactController.Db,
			Validator: contactController.Validator,
			Logger:    contactController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Get Contact Us By Email", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/contact/email/%s", contact.Email), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Messages retrieved successfully")
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		router, _ := setup()

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/contact/email/%s", contact.Email), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token is invalid!")
	})

	t.Run("Forbidden Access - Regular User Trying to Get By Email", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/contact/email/%s", contact.Email), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "role not authorized!")
	})
}
