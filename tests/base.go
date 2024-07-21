package tests

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models/migrations"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/newsletter"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func SetupRoutes(router *gin.Engine, newsController *newsletter.Controller) {
	router.POST("/newsletter", newsController.SubscribeNewsLetter)
}

func Setup() *utility.Logger {
	logger := utility.NewLogger()
	config := config.Setup(logger, "../app")

	postgresql.ConnectToDatabase(logger, config.TestDatabase)
	db := storage.Connection()
	if config.TestDatabase.Migrate {
		migrations.RunAllMigrations(db)
	}
	return logger
}

func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}

func AssertStatusCode(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("handler returned wrong status code: got status %d expected status %d", got, expected)
	}
}

func AssertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("handler returned wrong message: got message: %q expected: %q", got, expected)
	}
}
func AssertBool(t *testing.T, got, expected bool) {
	if got != expected {
		t.Errorf("handler returned wrong boolean: got %v expected %v", got, expected)
	}
}

func setupTestRouter() (*gin.Engine, *newsletter.Controller) {
	gin.SetMode(gin.TestMode)

	logger := Setup()
	db := storage.Connection()
	validator := validator.New()
	extReq := request.ExternalRequest{}

	newsController := &newsletter.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
		ExtReq:    extReq,
	}

	r := gin.Default()
	SetupRoutes(r, newsController)
	return r, newsController
}
