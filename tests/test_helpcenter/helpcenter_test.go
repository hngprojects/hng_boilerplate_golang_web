package test_helpcenter

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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/helpcenter"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestHelpCenterCreate(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/help-center/topics"}
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
		RequestBody  models.CreateHelpCenter
		ExpectedCode int
		Message      string
		ErrorField   string
		ErrorMessage string
		Headers      map[string]string
	}{
		{
			Name: "Successful help center post creation",
			RequestBody: models.CreateHelpCenter{
				Title:       "How to reset password",
				Content:     "To reset your password, go to the settings page and click 'Reset Password'.",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Topic added successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "Invalid content",
			RequestBody: models.CreateHelpCenter{
				Title:       "How to reset password",
				Content:     "",
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Input validation failed",
			ErrorField:   "CreateHelpCenter.Content",
			ErrorMessage: "Content is a required field",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, 
	}

	helpCenterController := helpcenter.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		helpCenterUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql))
		{
			helpCenterUrl.POST("/help-center/topics", helpCenterController.CreateHelpCenterTopic)
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

func TestFetchAllHelpCenterPosts(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/help-center/topics"

	helpCenterController := helpcenter.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r := gin.Default()
	helpCenterUrl := r.Group("/api/v1")
	{
		helpCenterUrl.GET("/help-center/topics", helpCenterController.FetchAllTopics)
	}

	tests := []struct {
		name         string
		expectedCode int
		message      string
	}{
		{
			name:         "Fetch all help center posts",
			expectedCode: http.StatusOK,
			message:      "Topics retrieved successfully.",
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

func TestFetchHelpCenterPostById(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/help-center/topics"
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

	helpCenterData := models.CreateHelpCenter{
		Title:   "How to reset password",
		Content: "To reset your password, go to the settings page and click 'Reset Password'.",
	}

	helpCenterController := helpcenter.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r.POST("/api/v1/help-center/topics", middleware.Authorize(db.Postgresql), helpCenterController.CreateHelpCenterTopic)

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(helpCenterData)

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
	helpCenterID, ok := data["data"].(map[string]interface{})["id"].(string)
	if !ok {
		t.Fatal("Failed to get help center post ID from response")
	}

	tests := []struct {
		name         string
		helpCenterID string
		expectedCode int
		message      string
	}{
		{
			name:         "Fetch help center post by ID",
			helpCenterID: helpCenterID,
			expectedCode: http.StatusOK,
			message:      "Topic retrieved successfully.",
		},
		{
			name:         "Help center post not found",
			helpCenterID: utility.GenerateUUID(),
			expectedCode: http.StatusNotFound,
			message:      "Topic not found",
		},
	}

	r.GET("/api/v1/help-center/topics/:id", helpCenterController.FetchTopicByID)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/help-center/topics/%s", test.helpCenterID), nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.expectedCode)
		})
	}
}

func TestUpdateHelpCenterPost(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/help-center/topics"
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

	helpCenterData := models.CreateHelpCenter{
		Title:   "How to reset password",
		Content: "To reset your password, go to the settings page and click 'Reset Password'.",
	}

	helpCenterController := helpcenter.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r.POST("/api/v1/help-center/topics", middleware.Authorize(db.Postgresql), helpCenterController.CreateHelpCenterTopic)

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(helpCenterData)

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
	helpCenterID := data["data"].(map[string]interface{})["id"].(string)

	tests := []struct {
		name         string
		helpCenterID string
		updateData   models.CreateHelpCenter
		expectedCode int
		message      string
	}{
		{
			name:         "Update help center post",
			helpCenterID: helpCenterID,
			updateData: models.CreateHelpCenter{
				Title:   "How to change your password",
				Content: "To change your password, go to the settings page and click 'Change Password'.",
			},
			expectedCode: http.StatusOK,
			message:      "Topic updated successfully",
		},
		{
			name:         "Help center post not found",
			helpCenterID: utility.GenerateUUID(),
			updateData: models.CreateHelpCenter{
				Title:   "How to change your password",
				Content: "To change your password, go to the settings page and click 'Change Password'.",
			},
			expectedCode: http.StatusNotFound,
			message:      "Topic not found",
		},
	}

	r.PATCH("/api/v1/help-center/topics/:id", middleware.Authorize(db.Postgresql), helpCenterController.UpdateHelpCenterByID)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.updateData)

			req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/help-center/topics/%s", test.helpCenterID), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.expectedCode)

			data := tst.ParseResponse(rr)
			message := data["message"].(string)
			tst.AssertResponseMessage(t, message, test.message)
		})
	}
}

func TestDeleteHelpCenterPost(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := "/api/v1/help-center/topics"
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

	helpCenterData := models.CreateHelpCenter{
		Title:   "How to reset password",
		Content: "To reset your password, go to the settings page and click 'Reset Password'.",
	}

	helpCenterController := helpcenter.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r.POST("/api/v1/help-center/topics", middleware.Authorize(db.Postgresql), helpCenterController.CreateHelpCenterTopic)

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(helpCenterData)

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
	helpCenterID := data["data"].(map[string]interface{})["id"].(string)

	tests := []struct {
		name         string
		helpCenterID string
		expectedCode int
	}{
		{
			name:         "Delete help center post",
			helpCenterID: helpCenterID,
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "Help center post not found",
			helpCenterID: utility.GenerateUUID(),
			expectedCode: http.StatusNotFound,
		},
	}

	r.DELETE("/api/v1/help-center/topics/:id", middleware.Authorize(db.Postgresql), helpCenterController.DeleteTopicByID)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/help-center/topics/%s", test.helpCenterID), nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.expectedCode)
		})
	}
}
