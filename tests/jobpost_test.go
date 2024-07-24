package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models/migrations"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/router"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	logger := utility.NewLogger()
	config := config.Setup(logger, "../app")
	postgresql.ConnectToDatabase(logger, config.TestDatabase)
	db := storage.Connection()
	if config.TestDatabase.Migrate {
		migrations.RunAllMigrations(db)
	}
	r := gin.Default()
	validate := validator.New()
	return router.JobPost(r, "/api/v1", validate, db, logger)
}

func getResponseData(w *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return nil, err
	}
	if data, ok := response["data"].(map[string]interface{}); ok {
		return data, nil
	}
	return nil, fmt.Errorf("response data is not in expected format")
}

func TestCreateJobPost(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()

	body := map[string]interface{}{
		"title":        "Software Engineer",
		"description":  "Develop and maintain software applications",
		"location":     "Remote",
		"salary":       100000,
		"job_type":     "Full-Time",
		"company_name": "TechCorp",
	}

	jsonValue, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/jobs", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	response, err := getResponseData(w)
	assert.NoError(t, err)
	assert.Equal(t, "Software Engineer", response["title"])
	assert.Equal(t, "Develop and maintain software applications", response["description"])
}

func TestFetchAllJobPost(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/jobs", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, data)
}

func TestFetchJobPostById(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()

	// First, create a job post to fetch later
	body := map[string]interface{}{
		"title":        "Software Engineer",
		"description":  "Develop and maintain software applications",
		"location":     "Remote",
		"salary":       100000,
		"job_type":     "Full-Time",
		"company_name": "TechCorp",
	}
	jsonValue, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/jobs", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	createResponse, err := getResponseData(w)
	assert.NoError(t, err)
	jobPostID, ok := createResponse["ID"].(string)
	assert.True(t, ok)

	// Now fetch the created job post by ID
	req, _ = http.NewRequest("GET", "/api/v1/jobs/"+jobPostID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	response, err := getResponseData(w)
	assert.NoError(t, err)
	assert.Equal(t, "Software Engineer", response["title"])

	// Fetch a non-existing job post
	req, _ = http.NewRequest("GET", "/api/v1/jobs/invalid-id", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}