package test_contact

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

func TestAddToContactUs(t *testing.T) {
	_, contactController := SetupContactTestRouter()
	db := contactController.Db.Postgresql
	currUUID := utility.GenerateUUID()
	password, _ := utility.HashPassword("password")

	regularUser := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Admin User",
		Email:    fmt.Sprintf("admin%v@qa.team", currUUID),
		Password: password,
		Role:     int(models.RoleIdentity.User),
	}

	db.Create(&regularUser)

	setup := func() (*gin.Engine, *auth.Controller) {
		router, contactController := SetupContactTestRouter()
		authController := auth.Controller{
			Db:        contactController.Db,
			Validator: contactController.Validator,
			Logger:    contactController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Create Contact Us", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		contactData := models.ContactUs{
			Name:    "John Doe",
			Email:   "johndoe@example.com",
			Subject: "</br><i>Inquiry about services3 with html",
			Message: "I would like to know more about your services3.",
		}
		payload, _ := json.Marshal(contactData)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/contact", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusCreated)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Message sent successfully")

		var createdContact models.ContactUs
		db.Last(&createdContact)
		if createdContact.Email != contactData.Email {
			t.Errorf("Expected contact email %s, but got %s", contactData.Email, createdContact.Email)
		}
	})

	t.Run("Missing Fields - Bad Request", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		contactData := models.ContactUs{
			Name: "John Doe",
		}
		payload, _ := json.Marshal(contactData)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/contact", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusBadRequest)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Failed to parse request body")
	})

	t.Run("Invalid Field Values - Unprocessable Entity", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		contactData := models.ContactUs{
			Name:    "John Doe",
			Email:   "invalid_email",
			Subject: "test subject",
			Message: "message test",
		}
		payload, _ := json.Marshal(contactData)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/contact", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})

}
