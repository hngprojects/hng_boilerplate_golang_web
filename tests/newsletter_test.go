package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func TestPostNewsletter_ValidateEmail(t *testing.T) {
	router, _ := setupTestRouter()

	body := map[string]interface{}{
		"Email": "invalid-email",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := ParseResponse(resp)
	AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
	AssertResponseMessage(t, response["message"].(string), "Validation failed")
}

func TestPostNewsletter_CheckDuplicateEmail(t *testing.T) {
	router, newsController := setupTestRouter()

	db := newsController.Db.Postgresql
	db.Create(&models.NewsLetter{Email: "test@example.com"})

	body := map[string]interface{}{
		"Email": "test@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := ParseResponse(resp)
	AssertStatusCode(t, resp.Code, http.StatusConflict)
	AssertResponseMessage(t, response["message"].(string), "Email already subscribed")
}

func TestPostNewsletter_SaveData(t *testing.T) {
	router, newsController := setupTestRouter()

	body := map[string]interface{}{
		"Email": "test2@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := ParseResponse(resp)
	AssertStatusCode(t, resp.Code, http.StatusCreated)
	AssertResponseMessage(t, response["message"].(string), "subscribed successfully")

	var newsletter models.NewsLetter
	newsController.Db.Postgresql.First(&newsletter, "email = ?", "test2@example.com")
	if newsletter.Email != "test2@example.com" {
		t.Errorf("data not saved correctly to the database: expected email %s, got %s", "test2@example.com", newsletter.Email)
	}
}

func TestPostNewsletter_ResponseAndStatusCode(t *testing.T) {
	router, _ := setupTestRouter()

	body := map[string]interface{}{
		"Email": "test3@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := ParseResponse(resp)
	AssertStatusCode(t, resp.Code, http.StatusCreated)
	AssertResponseMessage(t, response["message"].(string), "subscribed successfully")
}
