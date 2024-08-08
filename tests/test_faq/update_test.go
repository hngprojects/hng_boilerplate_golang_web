package test_faq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestUpdateFaq(t *testing.T) {
	setup := func() (*gin.Engine, *auth.Controller) {
		router, faqController := SetupFAQTestRouter()
		authController := auth.Controller{
			Db:        faqController.Db,
			Validator: faqController.Validator,
			Logger:    faqController.Logger,
		}
		return router, &authController
	}

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

	faq := models.FAQ{
		ID:        utility.GenerateUUID(),
		Question:  fmt.Sprintf("What is the purpose of this %s FAQ?", utility.RandomString(9)),
		Answer:    "To provide answers to frequently asked questions.",
		Category:  "Policies",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db.Create(&adminUser)
	db.Create(&faq)

	t.Run("Successful Update FAQ", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateFaq := models.UpdateFAQ{
			Question: fmt.Sprintf("Update of this %s FAQ?", utility.RandomString(8)),
			Answer:   "Updated answer",
			Category: "Policies",
		}
		jsonBody, _ := json.Marshal(updateFaq)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/faq/%s", faq.ID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "FAQ updated successfully")
	})

	t.Run("Validation Error - Missing Fields", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateFaq := models.UpdateFAQ{
			Question: "",
			Answer:   "",
		}
		jsonBody, _ := json.Marshal(updateFaq)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/faq/%s", faq.ID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})

	t.Run("FAQ Not Found", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		updateFaq := models.UpdateFAQ{
			Question: fmt.Sprintf("the purpose of this %s FAQ?", utility.RandomString(10)),
			Answer:   "Updated answer",
			Category: "Policies1",
		}
		jsonBody, _ := json.Marshal(updateFaq)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/faq/%s", currUUID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusBadRequest)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "record not found")
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		router, _ := setup()

		updateFaq := models.UpdateFAQ{
			Question: fmt.Sprintf("What is the purpose of this %s FAQ?", utility.RandomString(15)),
			Answer:   "Updated answer",
		}
		jsonBody, _ := json.Marshal(updateFaq)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/faq/%s", faq.ID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token is invalid!")
	})
}
