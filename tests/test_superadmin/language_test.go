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
		}
		authController.Db.Postgresql.Create(&language)

		duplicateLanguage := models.Language{
			Name: fmt.Sprintf("English-%s", theRandom),
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