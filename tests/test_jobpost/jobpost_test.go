package test_jobpost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/jobpost"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestJobPostCreate(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/jobs"}
	currUUID := utility.GenerateUUID()
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}
	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData, false)

	token := tst.GetLoginToken(t, r, auth, loginData)

	tests := []struct {
		Name         string
		RequestBody  models.CreateJobPostModel
		ExpectedCode int
		Message      string
		ErrorField   string
		ErrorMessage string
		Headers      map[string]string
	}{
		{
			Name: "Successful job post creation",
			RequestBody: models.CreateJobPostModel{
				Title:               "Software Engineer Intern",
				SalaryRange:         "below_30k",
				JobType:             "internship",
				Location:            "San Francisco, CA",
				Deadline:            time.Now().AddDate(0, 1, 0),
				JobMode:             "remote",
				ExperienceLevel:     "0-2 years",
				Benefits:            "Flexible hours, Remote work, Health insurance",
				CompanyName:         "Tech Innovators",
				Description:         "We are looking for a passionate Software Engineer Intern to join our team. You will be working on exciting projects and gain hands-on experience.",
				KeyResponsibilities: "Develop and maintain web applications, Collaborate with the team on various projects, Participate in code reviews",
				Qualifications:      "Ability to work solo, Bachelor degree",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Job created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "Invalid job type",
			RequestBody: models.CreateJobPostModel{
				Title:               "Software Engineer Intern",
				SalaryRange:         "below_30k",
				JobType:             "",
				Location:            "San Francisco, CA",
				Deadline:            time.Now().AddDate(0, 1, 0),
				JobMode:             "remote",
				ExperienceLevel:     "0-2 years",
				Benefits:            "Flexible hours, Remote work, Health insurance",
				CompanyName:         "Tech Innovators",
				Description:         "We are looking for a passionate Software Engineer Intern to join our team. You will be working on exciting projects and gain hands-on experience.",
				KeyResponsibilities: "Develop and maintain web applications, Collaborate with the team on various projects, Participate in code reviews",
				Qualifications:      "Ability to work solo, Bachelor degree",
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
			ErrorField:   "CreateJobPostModel.JobType",
			ErrorMessage: "JobType is a required field",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "User unauthorized",
			RequestBody: models.CreateJobPostModel{
				Title:               "Software Engineer Intern",
				SalaryRange:         "below_30k",
				JobType:             "internship",
				Location:            "San Francisco, CA",
				Deadline:            time.Now().AddDate(0, 1, 0),
				JobMode:            "remote",
				ExperienceLevel:     "2 years",
				Benefits:            "Flexible hours, Remote work, Health insurance",
				CompanyName:         "Tech Innovators",
				Description:         "We are looking for a passionate Software Engineer Intern to join our team. You will be working on exciting projects and gain hands-on experience.",
				KeyResponsibilities: "Develop and maintain web applications, Collaborate with the team on various projects, Participate in code reviews",
				Qualifications:      "Ability to work solo, Bachelor degree",
			},
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
		},
	}

	jobPostController := jobpost.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		jobUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql))
		{
			jobUrl.POST("/jobs", jobPostController.CreateJobPost)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["status_code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

				if test.ErrorField != "" {
					errorData := data["error"].(map[string]interface{})
					errorMessage := errorData[test.ErrorField].(string)
					tst.AssertResponseMessage(t, errorMessage, test.ErrorMessage)
				}
			}
		})
	}
}

func TestFetchAllJobPost(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/jobs"

	jobPostController := jobpost.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r := gin.Default()
	jobUrl := r.Group("/api/v1")
	{
		jobUrl.GET("/jobs", jobPostController.FetchAllJobPost)
	}

	tests := []struct {
		name         string
		expectedCode int
		message      string
	}{
		{
			name:         "Fetch all job posts",
			expectedCode: http.StatusOK,
			message:      "Job listings retrieved successfully.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, requestURI, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.expectedCode)

			data := tst.ParseResponse(rr)
			message := data["message"].(string)
			tst.AssertResponseMessage(t, message, test.message)
		})
	}
}

func TestFetchJobPostById(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/jobs"
	currUUID := utility.GenerateUUID()
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}
	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	authController := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, authController, userSignUpData, false)

	token := tst.GetLoginToken(t, r, authController, loginData)

	jobPostData := models.CreateJobPostModel{
		Title:               "Software Engineer Intern",
		SalaryRange:         "below_30k",
		JobType:             "internship",
		Location:            "San Francisco, CA",
		Deadline:            time.Now().AddDate(0, 1, 0),
		JobMode:             "remote",
		ExperienceLevel:     "2 years",
		Benefits:            "Flexible hours, Remote work, Health insurance",
		CompanyName:         "Tech Innovators",
		Description:         "We are looking for a passionate Software Engineer Intern to join our team. You will be working on exciting projects and gain hands-on experience.",
		KeyResponsibilities: "Develop and maintain web applications, Collaborate with the team on various projects, Participate in code reviews",
		Qualifications:      "Ability to work solo, Bachelor degree",
	}

	jobPostController := jobpost.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r.POST("/api/v1/jobs", jobPostController.CreateJobPost)

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(jobPostData)

	req, err := http.NewRequest(http.MethodPost, requestURI, &b)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	tst.AssertStatusCode(t, rr.Code, http.StatusCreated)

	data := tst.ParseResponse(rr)
	jobPostID := data["data"].(map[string]interface{})["id"].(string)

	tests := []struct {
		name         string
		expectedCode int
		message      string
		jobPostID    string
	}{
		{
			name:         "Fetch job post by ID",
			expectedCode: http.StatusOK,
			message:      "Job post retrieved successfully",
			jobPostID:    jobPostID,
		},
		{
			name:         "Invalid uuid format",
			expectedCode: http.StatusBadRequest,
			message:      "Invalid ID format",
			jobPostID:    "invalidIDFormat",
		},
	}

	r.GET("/api/v1/jobs/:job_id", jobPostController.FetchJobPostByID)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/jobs/%s", test.jobPostID), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.expectedCode)

			data := tst.ParseResponse(rr)
			message := data["message"].(string)
			tst.AssertResponseMessage(t, message, test.message)
		})
	}
}

func TestUpdateJobPostById(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/jobs"
	currUUID := utility.GenerateUUID()
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}
	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	authController := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, authController, userSignUpData, false)

	token := tst.GetLoginToken(t, r, authController, loginData)

	jobPostData := models.CreateJobPostModel{
		Title:               "Software Engineer Intern",
		SalaryRange:         "below_30k",
		JobType:             "internship",
		Location:            "San Francisco, CA",
		Deadline:            time.Now().AddDate(0, 1, 0),
		JobMode:             "remote",
		ExperienceLevel:     "2 years",
		Benefits:            "Flexible hours, Remote work, Health insurance",
		CompanyName:         "Tech Innovators",
		Description:         "We are looking for a passionate Software Engineer Intern to join our team. You will be working on exciting projects and gain hands-on experience.",
		KeyResponsibilities: "Develop and maintain web applications, Collaborate with the team on various projects, Participate in code reviews",
		Qualifications:      "Ability to work solo, Bachelor degree",
	}

	jobPostController := jobpost.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r.POST("/api/v1/jobs", jobPostController.CreateJobPost)
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(jobPostData)

	req, err := http.NewRequest(http.MethodPost, requestURI, &b)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	tst.AssertStatusCode(t, rr.Code, http.StatusCreated)

	data := tst.ParseResponse(rr)
	jobPostID := data["data"].(map[string]interface{})["id"].(string)

	updatedJobPostData := models.JobPost{
		Title:               "Updated Software Engineer Intern",
		SalaryRange:         "below_30k",
		JobType:             "full-time",
		Location:            "San Francisco, CA",
		Deadline:            time.Now().AddDate(0, 1, 0),
		JobMode:             "hybrid",
		ExperienceLevel:     "3 years",
		Benefits:            "Flexible hours, Health insurance, Stock options",
		CompanyName:         "Tech Innovators Inc.",
		Description:         "We are looking for a passionate Software Engineer Intern to join our team. You'll work on exciting projects and gain valuable experience.",
		KeyResponsibilities: "Develop and maintain web applications, Collaborate with the team on various projects, Participate in code reviews",
		Qualifications:      "Ability to work in a team, Bachelor degree in Computer Science",
	}

	tests := []struct {
		name         string
		expectedCode int
		message      string
		jobPostID    string
		updateData   models.JobPost
	}{
		{
			name:         "Update job post",
			expectedCode: http.StatusOK,
			message:      "Job post updated successfully",
			jobPostID:    jobPostID,
			updateData:   updatedJobPostData,
		},
	}

	r.PATCH("/api/v1/jobs/:job_id", jobPostController.UpdateJobPostByID)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.updateData)

			req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/jobs/%s", test.jobPostID), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.expectedCode)
		})
	}
}

func TestDeleteJobPostById(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/jobs"
	currUUID := utility.GenerateUUID()
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}
	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	authController := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, authController, userSignUpData, false)

	token := tst.GetLoginToken(t, r, authController, loginData)

	jobPostData := models.CreateJobPostModel{
		Title:               "Software Engineer Intern",
		SalaryRange:         "below_30k",
		JobType:             "internship",
		Location:            "San Francisco, CA",
		Deadline:            time.Now().AddDate(0, 1, 0),
		JobMode:             "remote",
		ExperienceLevel:     "2 years",
		Benefits:            "Flexible hours, Remote work, Health insurance",
		CompanyName:         "Tech Innovators",
		Description:         "We are looking for a passionate Software Engineer Intern to join our team. You will be working on exciting projects and gain hands-on experience.",
		KeyResponsibilities: "Develop and maintain web applications, Collaborate with the team on various projects, Participate in code reviews",
		Qualifications:      "Ability to work solo, Bachelor degree",
	}

	jobPostController := jobpost.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r.POST("/api/v1/jobs", jobPostController.CreateJobPost)
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(jobPostData)

	req, err := http.NewRequest(http.MethodPost, requestURI, &b)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	tst.AssertStatusCode(t, rr.Code, http.StatusCreated)

	data := tst.ParseResponse(rr)
	jobPostID := data["data"].(map[string]interface{})["id"].(string)

	tests := []struct {
		name         string
		expectedCode int
		jobPostID    string
	}{
		{
			name:         "Delete existing job post",
			expectedCode: http.StatusNoContent,
			jobPostID:    jobPostID,
		},
		{
			name:         "Delete non-existent job post",
			expectedCode: http.StatusNotFound,
			jobPostID:    jobPostID,
		},
	}

	r.DELETE("/api/v1/jobs/:job_id", jobPostController.DeleteJobPostByID)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/jobs/%s", test.jobPostID), nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.expectedCode)

			if test.expectedCode == http.StatusNoContent {
				req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/jobs/%s", test.jobPostID), nil)
				if err != nil {
					t.Fatal(err)
				}
				rr = httptest.NewRecorder()
				r.ServeHTTP(rr, req)
				tst.AssertStatusCode(t, rr.Code, http.StatusNotFound)
			}
		})
	}
}
