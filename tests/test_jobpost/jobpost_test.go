package test_jobpost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/jobpost"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
)


func setupJobPostRouter() (*gin.Engine, *jobpost.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tests.Setup()
	db := storage.Connection()
	validator := validator.New()

	jobpostCtrl := &jobpost.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupJopPostRoutes(r, jobpostCtrl)
	return r, jobpostCtrl
}

func SetupJopPostRoutes(r *gin.Engine, jobpostCtrl *jobpost.Controller ) {
	r.POST("/api/v1/jobs", jobpostCtrl.CreateJobPost )
	r.GET("/api/v1/jobs", jobpostCtrl.FetchAllJobPost)
	r.GET("/api/v1/jobs/:ID", jobpostCtrl.FetchJobPostByID)
}

func TestCreateJobPost(t *testing.T) {
	router, _ := setupJobPostRouter()
	body := models.JobPost{
		Title:        "Software Engineer(Test)",
		Description:  "Develop and maintain software applications",
		Location:     "Remote",
		Salary:       100000,
		JobType:      "Full-Time",
		CompanyName:  "TechCorp",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/jobs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := tests.ParseResponse(resp)
	tests.AssertStatusCode(t, resp.Code, http.StatusCreated)
	if titleValue, ok := response["data"].(map[string]interface{})["title"].(string); ok {
  		tests.AssertResponseMessage(t, titleValue, "Software Engineer(Test)")
	} else {
  	t.Errorf("Expected title field in response data, but got nil")
	}
	if descriptionValue, ok := response["data"].(map[string]interface{})["description"].(string); ok {
  		tests.AssertResponseMessage(t, descriptionValue, "Develop and maintain software applications")
	} else {
  	t.Errorf("Expected description field in response data, but got nil")
	}	
	if company_nameValue, ok := response["data"].(map[string]interface{})["company_name"].(string); ok {
  		tests.AssertResponseMessage(t, company_nameValue, "TechCorp")
	} else {
  	t.Errorf("Expected company_name field in response data, but got nil")
	}
} 

func TestFetchAllJobPost(t *testing.T) {
	router, _ := setupJobPostRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/jobs", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := tests.ParseResponse(resp)
	tests.AssertStatusCode(t, resp.Code, http.StatusOK)

	data, ok := response["data"].([]interface{})
	if !ok {
		t.Errorf("Expected response data to be an array, got %T", response["data"])
	}

	if len(data) == 0 {
		t.Errorf("Expected response data to be non-empty")
	}
}

func TestFetchJobPostByID(t *testing.T) {
	router, _ := setupJobPostRouter()

	body := models.JobPost{
		Title:        "Backend dev needed",
		Description:  "blah blah blah",
		Location:     "Remote",
		Salary:       1000,
		JobType:      "Health",
		CompanyName:  "BMW",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/jobs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := tests.ParseResponse(resp)
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract job post data from creation response. Response: %v", response)
	}

	jobPostID, ok := data["ID"].(string)
	if !ok {
		t.Fatalf("Failed to extract job post ID from creation response. Response: %v", response)
	}

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/jobs/"+jobPostID, nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req)

	bodyBytes, err := io.ReadAll(resp2.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Printf("Response body from fetch job post by ID: %s\n", string(bodyBytes))

	tests.AssertStatusCode(t, resp2.Code, http.StatusOK)
}

