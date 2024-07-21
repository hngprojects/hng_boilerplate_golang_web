package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/newsletter"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
)

func setupNewsLetterTestRouter() (*gin.Engine, *newsletter.Controller) {
	gin.SetMode(gin.TestMode)

	logger := Setup()
	db := storage.Connection()
	validator := validator.New()

	newsController := &newsletter.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupNewsLetterRoutes(r, newsController)
	return r, newsController
}

func SetupNewsLetterRoutes(r *gin.Engine, newsController *newsletter.Controller) {
	r.POST("/api/v1/newsletter", newsController.SubscribeNewsLetter)
}

func TestE2ENewsletterSubscription(t *testing.T) {
	router, _ := setupNewsLetterTestRouter()

	// Test POST /newsletter
	body := models.NewsLetter{
		Email: "e2e_test@example.com",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/api/v1/newsletter", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	AssertStatusCode(t, resp.Code, http.StatusCreated)

	response := ParseResponse(resp)
	AssertResponseMessage(t, response["message"].(string), "subscribed successfully")
}

func TestPostNewsletter_ValidateEmail(t *testing.T) {
	router, _ := setupNewsLetterTestRouter()

	body := models.NewsLetter{
		Email: "invalid-email",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := ParseResponse(resp)
	AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
	AssertResponseMessage(t, response["message"].(string), "Validation failed")
}

func TestPostNewsletter_CheckDuplicateEmail(t *testing.T) {
	router, newsController := setupNewsLetterTestRouter()

	db := newsController.Db.Postgresql
	db.Create(&models.NewsLetter{Email: "test@example.com"})

	body := models.NewsLetter{
		Email: "test@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := ParseResponse(resp)
	AssertStatusCode(t, resp.Code, http.StatusConflict)
	AssertResponseMessage(t, response["message"].(string), "Email already subscribed")
}

func TestPostNewsletter_SaveData(t *testing.T) {
	router, newsController := setupNewsLetterTestRouter()

	body := models.NewsLetter{
		Email: "test2@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/newsletter", bytes.NewBuffer(jsonBody))
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
	router, _ := setupNewsLetterTestRouter()

	body := models.NewsLetter{
		Email: "test3@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := ParseResponse(resp)
	AssertStatusCode(t, resp.Code, http.StatusCreated)
	AssertResponseMessage(t, response["message"].(string), "subscribed successfully")
}
