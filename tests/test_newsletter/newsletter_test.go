package test_newsletter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/newsletter"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func setupNewsLetterTestRouter() (*gin.Engine, *newsletter.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tst.Setup()
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
	currUUID := utility.GenerateUUID()
	body := models.NewsLetter{
		Email: fmt.Sprintf("testuser%v@qa.team", currUUID),
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

	tst.AssertStatusCode(t, resp.Code, http.StatusCreated)

	response := tst.ParseResponse(resp)
	tst.AssertResponseMessage(t, response["message"].(string), "subscribed successfully")
}

func TestPostNewsletter_ValidateEmail(t *testing.T) {
	router, _ := setupNewsLetterTestRouter()

	currUUID := utility.GenerateUUID()
	body := models.NewsLetter{
		Email: fmt.Sprintf("testuser%v@qa", currUUID),
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := tst.ParseResponse(resp)
	tst.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
	tst.AssertResponseMessage(t, response["message"].(string), "Validation failed")
}

func TestPostNewsletter_CheckDuplicateEmail(t *testing.T) {
	router, newsController := setupNewsLetterTestRouter()

	currUUID := utility.GenerateUUID()

	db := newsController.Db.Postgresql
	db.Create(&models.NewsLetter{Email: fmt.Sprintf("testuser%v@qa.team", currUUID)})

	body := models.NewsLetter{
		Email: fmt.Sprintf("testuser%v@qa.team", currUUID),
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := tst.ParseResponse(resp)
	tst.AssertStatusCode(t, resp.Code, http.StatusConflict)
	tst.AssertResponseMessage(t, response["message"].(string), "Email already subscribed")
}

func TestPostNewsletter_SaveData(t *testing.T) {
	router, newsController := setupNewsLetterTestRouter()

	currUUID := utility.GenerateUUID()
	body := models.NewsLetter{
		Email: fmt.Sprintf("testuser%v@qa.team", currUUID),
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := tst.ParseResponse(resp)
	tst.AssertStatusCode(t, resp.Code, http.StatusCreated)
	tst.AssertResponseMessage(t, response["message"].(string), "subscribed successfully")

	var newsletter models.NewsLetter
	newsController.Db.Postgresql.First(&newsletter, "email = ?", fmt.Sprintf("testuser%v@qa.team", currUUID))
	if newsletter.Email != fmt.Sprintf("testuser%v@qa.team", currUUID) {
		t.Errorf("data not saved correctly to the database: expected email %s, got %s", fmt.Sprintf("testuser%v@qa.team", currUUID), newsletter.Email)
	}
}

func TestPostNewsletter_ResponseAndStatusCode(t *testing.T) {
	router, _ := setupNewsLetterTestRouter()

	currUUID := utility.GenerateUUID()
	body := models.NewsLetter{
		Email: fmt.Sprintf("testuser%v@gmail.com", currUUID),
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/newsletter", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := tst.ParseResponse(resp)
	tst.AssertStatusCode(t, resp.Code, http.StatusCreated)
	tst.AssertResponseMessage(t, response["message"].(string), "subscribed successfully")
}
