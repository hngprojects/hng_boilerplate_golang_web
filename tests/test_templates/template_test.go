package test_templates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	template "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/templates"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

var invalidToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjAxOTBlMWY0LWYwZDQtNzI4NS1hOWY4LTA3ZmE3ZDA5MjZhNyIsImF1dGhvcmlzZWQiOnRydWUsImV4cCI6MTcyMTk1MDY0NCwicm9sZSI6MSwidXNlcl9pZCI6IjAxOTBlMWYzLWViZTktNzI4NC04MGMzLTEwNjg5NTUzYTQ5NyJ9.Ahrh9l0FJAEEaKIHnph54tdY5U8dEGQiYKiFp6g"

func TestTemplatePostCreate(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/template"}
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
		RequestBody  models.TemplateRequest
		ExpectedCode int
		Message      string
		ErrorField   string
		ErrorMessage string
		Headers      map[string]string
	}{
		{
			Name: "Successful template post creation",
			RequestBody: models.TemplateRequest{
				Name: "Test Template",
				Body: "This is a test template",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Template created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "Invalid template type",
			RequestBody: models.TemplateRequest{
				Name: "Test Template",
				Body: "",
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
			ErrorField:   "TemplateRequest.Body",
			ErrorMessage: "Body is a required field",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "User unauthorized",
			RequestBody: models.TemplateRequest{
				Name: "Test Template",
				Body: "This is a test template",
			},
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
		},
	}

	templateController := template.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		templateUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql))
		{
			templateUrl.POST("/template", templateController.CreateTemplate)
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

func TestFetchAllTemplates(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/template"
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

	tests := []struct {
		name         string
		expectedCode int
		message      string
		headers      map[string]string
	}{
		{
			name:         "Fetch all Templates",
			expectedCode: http.StatusOK,
			message:      "Templates Successfully retrieved",
			headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	templateController := template.Controller{Db: db, Validator: validatorRef, Logger: logger}

	templateUrl := r.Group("/api/v1")
	{
		templateUrl.GET("/template", templateController.GetTemplates)
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

func TestFetchTemplateById(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/template"
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

	templateData := models.EmailTemplate{
		Name: "Test Template",
		Body: "This is a test template",
	}

	templateController := template.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r.POST("/api/v1/template", templateController.CreateTemplate)

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(templateData)

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
	templateID := data["data"].(map[string]interface{})["id"].(string)

	tests := []struct {
		name         string
		expectedCode int
		message      string
		templateID    string
	}{
		{
			name:         "Fetch Template by ID",
			expectedCode: http.StatusOK,
			message:      "Template Successfully retrieved",
			templateID:    templateID,
		},
		{
			name:         "Invalid uuid format",
			expectedCode: http.StatusBadRequest,
			message:      "Invalid id",
			templateID:    "invalidIDFormat",
		},
	}

	r.GET("/api/v1/template/:id", templateController.GetTemplate)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/template/%s", test.templateID), nil)
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

func TestUpdateTemplateById(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/template"
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

	templateData := models.TemplateRequest{
		Name: "Test Template",
		Body: "This is a test template",
	}

	templateController := template.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r.POST("/api/v1/template", templateController.CreateTemplate)
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(templateData)

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
	templateID := data["data"].(map[string]interface{})["id"].(string)

	updatedTemplateData := models.TemplateRequest{
		Name: "Updated Test Template",
		Body: "This is an updated test template",
	}

	tests := []struct {
		name         string
		expectedCode int
		message      string
		templateID   string
		updateData   models.TemplateRequest
	}{
		{
			name:         "Update template",
			expectedCode: http.StatusOK,
			message:      "Template updated successfully",
			templateID:   templateID,
			updateData:   updatedTemplateData,
		},
	}

	r.PUT("/api/v1/template/:id", templateController.UpdateTemplate)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.updateData)

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/template/%s", test.templateID), &b)
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
	requestURI := "/api/v1/template"
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

	templateData := models.EmailTemplate{
		Name: "Test Template",
		Body: "This is a test template",
	}

	templateController := template.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r.POST("/api/v1/template", templateController.CreateTemplate)
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(templateData)

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
	templateID := data["data"].(map[string]interface{})["id"].(string)

	tests := []struct {
		name         string
		expectedCode int
		templateID   string
	}{
		{
			name:         "Delete existing template",
			expectedCode: http.StatusOK,
			templateID:   templateID,
		},
		{
			name:         "Delete non-existent job post",
			expectedCode: http.StatusInternalServerError,
			templateID:   templateID,
		},
	}

	r.DELETE("/api/v1/template/:id", templateController.DeleteTemplate)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/template/%s", test.templateID), nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.expectedCode)

			if test.expectedCode == http.StatusNoContent {
				req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/template/%s", test.templateID), nil)
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
