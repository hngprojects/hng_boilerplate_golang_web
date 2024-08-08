package test_faq

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

func TestAddToFaq(t *testing.T) {
	_, newsController := SetupFAQTestRouter()
	db := newsController.Db.Postgresql

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
		router, newsController := SetupFAQTestRouter()
		authController := auth.Controller{
			Db:        newsController.Db,
			Validator: newsController.Validator,
			Logger:    newsController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Add to FAQ", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		faq := models.FAQ{
			Question: fmt.Sprintf("What is the purpose of this %s FAQ?", utility.RandomString(6)),
			Answer:   "To provide answers to frequently asked questions.",
			Category: "Policies",
		}
		jsonBody, _ := json.Marshal(faq)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/faq", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusCreated)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "FAQ added successfully")
	})

	t.Run("Validation Error - Missing Fields", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		faq := models.FAQ{
			Question: fmt.Sprintf("What is the purpose a %s FAQ?", utility.RandomString(5)),
			Answer:   "",
			Category: "Policies",
		}
		jsonBody, _ := json.Marshal(faq)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/faq", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})

	t.Run("Duplicate Question", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)
		theRandom := utility.RandomString(6)

		faq := models.FAQ{
			Question: fmt.Sprintf("What is the purpose of this %s FAQ?", theRandom),
			Answer:   "This is a duplicate question test.",
			Category: "Policies",
		}
		authController.Db.Postgresql.Create(&faq)

		duplicateFaq := models.FAQ{
			Question: fmt.Sprintf("What is the purpose of this %s FAQ?", theRandom),
			Answer:   "To provide answers to frequently asked questions.",
			Category: "Policies",
		}
		jsonBody, _ := json.Marshal(duplicateFaq)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/faq", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusBadRequest)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "question exists")
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		router, _ := setup()

		faq := models.FAQ{
			Question: fmt.Sprintf("What is the purpose of this %s FAQ?", utility.RandomString(7)),
			Answer:   "To provide answers to frequently asked questions.",
			Category: "Policies",
		}
		jsonBody, _ := json.Marshal(faq)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/faq", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token is invalid!")
	})
}
