package test_auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestRequestMagicLink(t *testing.T) {
	router, _ := SetupAuthTestRouter()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	theRole := models.RoleIdentity.SuperAdmin
	password, _ := utility.HashPassword("password")

	adminData := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "admin jane doe2",
		Email:    fmt.Sprintf("testadmin%v@qa.team", currUUID),
		Password: password,
		Role:     int(theRole),
	}
	db.Postgresql.Create(&adminData)

	t.Run("Successful Magic Link Request Email Sent", func(t *testing.T) {
		requestMagicLink := models.MagicLinkRequest{
			Email: adminData.Email,
		}
		reqBody, _ := json.Marshal(requestMagicLink)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/magick-link", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Magic link sent to email")
	})

}
