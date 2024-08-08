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

func TestAddToLanguage(t *testing.T) {
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

	t.Run("Successful Add to Language", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		language := models.Language{
			Name: fmt.Sprintf("English-%s", utility.RandomString(10)),
			Code: fmt.Sprintf("en-%s", utility.RandomString(5)),
		}
		jsonBody, _ := json.Marshal(language)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/languages", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusCreated)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Language added successfully")
	})

	t.Run("Validation Error - Missing Fields", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		language := models.Language{
			Name: "",
			Code: fmt.Sprintf("en-%s", utility.RandomString(5)),
		}
		jsonBody, _ := json.Marshal(language)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/languages", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})

	t.Run("Duplicate Language", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)
		theRandom := utility.RandomString(6)

		language := models.Language{
			Name: fmt.Sprintf("English-%s", theRandom),
			Code: fmt.Sprintf("NA-%s", utility.RandomString(5)),
		}
		authController.Db.Postgresql.Create(&language)

		duplicateLanguage := models.Language{
			Name: fmt.Sprintf("English-%s", theRandom),
			Code: fmt.Sprintf("NA-%s", utility.RandomString(5)),
		}
		jsonBody, _ := json.Marshal(duplicateLanguage)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/languages", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusBadRequest)
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		router, _ := setup()

		language := models.Language{
			Name: fmt.Sprintf("Spanish-%s", utility.RandomString(8)),
			Code: fmt.Sprintf("sp-%s", utility.RandomString(5)),
		}
		jsonBody, _ := json.Marshal(language)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/languages", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token is invalid!")
	})
}

func TestGetLanguages(t *testing.T) {
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

	language := models.Language{
		Name: fmt.Sprintf("English-%s", utility.RandomString(10)),
		Code: fmt.Sprintf("en-%s", utility.RandomString(5)),
	}

	db.Create(&adminUser)
	db.Create(&regularUser)
	db.Create(&language)

	t.Run("Successful Get Languages for admin", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/languages", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Languages retrieved successfully")
	})

	t.Run("Successful Get Languages for user", func(t *testing.T) {
		router, authController := setup()
		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/languages", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Languages retrieved successfully")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		router, _ := setup()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/languages", nil)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token could not be found!")
	})
}
