package tests

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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/account"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestAddRecoveryEmail(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/account/add-recovery-email"}
	currUUID := utility.GenerateUUID()

	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+23481%v", utility.GetRandomNumbersInRange(10000000, 99999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}

	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)

	token := GetLoginToken(t, r, user, loginData)

	tests := []struct {
		Name         string
		RequestBody  models.AddRecoveryEmailRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Add recovery email success",
			RequestBody: models.AddRecoveryEmailRequestModel{
				Email: fmt.Sprintf("testrecoveryemail_%v@qa.team", currUUID),
			},
			ExpectedCode: http.StatusOK,
			Message:      "Recovery email successfully added",
		},
		{
			Name: "Add recovery email faliure",
			RequestBody: models.AddRecoveryEmailRequestModel{
				Email: "carti",
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Invalid recovery email",
		},
	}

	account := account.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		accountGroup := r.Group(fmt.Sprintf("%v", "/api/v1/account"), middleware.Authorize())
		{
			accountGroup.POST("/add-recovery-email", account.AddRecoveryEmail)

		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := ParseResponse(rr)

			code := int(data["status_code"].(float64))
			AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					AssertResponseMessage(t, message.(string), test.Message)
				} else {
					AssertResponseMessage(t, "", test.Message)
				}
			}
		})
	}
}

func TestAddRecoveryPhone(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/account/recovery-number"}
	currUUID := utility.GenerateUUID()

	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+23481%v", utility.GetRandomNumbersInRange(10000000, 99999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}

	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)

	token := GetLoginToken(t, r, user, loginData)

	tests := []struct {
		Name         string
		RequestBody  models.AddRecoveryPhoneNumberRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Add recovery phone number success",
			RequestBody: models.AddRecoveryPhoneNumberRequestModel{
				PhoneNumber: fmt.Sprintf("+23480%v", utility.GetRandomNumbersInRange(10000000, 99999999)),
			},
			ExpectedCode: http.StatusOK,
			Message:      "Recovery phone number successfully added",
		},
		{
			Name: "Add recovery email faliure",
			RequestBody: models.AddRecoveryPhoneNumberRequestModel{
				PhoneNumber: "skibidi",
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Invalid phone number",
		},
	}

	account := account.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		accountGroup := r.Group(fmt.Sprintf("%v", "/api/v1/account"), middleware.Authorize())
		{
			accountGroup.POST("/recovery-number", account.AddRecoveryPhoneNumber)

		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := ParseResponse(rr)

			code := int(data["status_code"].(float64))
			AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					AssertResponseMessage(t, message.(string), test.Message)
				} else {
					AssertResponseMessage(t, "", test.Message)
				}
			}
		})
	}
}

func TestAddSecurityQuestions(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/account/submit-security-answers"}
	currUUID := utility.GenerateUUID()

	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+23490%v", utility.GetRandomNumbersInRange(10000000, 99999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}

	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)

	token := GetLoginToken(t, r, user, loginData)

	tests := []struct {
		Name         string
		RequestBody  map[string][]map[string]string
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Add recovery security questions success",
			RequestBody: map[string][]map[string]string{
				"answers": {
					{
						"question_1": "What is your mother's maiden name?",
						"answer_1":   "User's Answer",
					},
					{
						"question_2": "In what city were you born?",
						"answer_2":   "User's Answer",
					},
					{
						"question_3": "What is the name of your first pet?",
						"answer_3":   "User's Answer",
					},
				},
			},
			ExpectedCode: http.StatusOK,
			Message:      "Security answers submitted successfully",
		},
		{
			Name: "Add recovery email faliure",
			RequestBody: map[string][]map[string]string{
				"answers": {
					{
						"question_1": "What is your mother's maiden name?",
					},
					{
						"question_2": "In what city were you born?",
					},
					{
						"question_3": "What is the name of your first pet?",
					},
				},
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Could not submit security questions",
		},
	}

	account := account.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		accountGroup := r.Group(fmt.Sprintf("%v", "/api/v1/account"), middleware.Authorize())
		{
			accountGroup.POST("/submit-security-answers", account.AddSecurityAnswers)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := ParseResponse(rr)

			code := int(data["status_code"].(float64))
			AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					AssertResponseMessage(t, message.(string), test.Message)
				} else {
					AssertResponseMessage(t, "", test.Message)
				}
			}
		})
	}
}

func TestUpdateRecoveryOptions(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/account/update-recovery-options"}
	currUUID := utility.GenerateUUID()

	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+23481%v", utility.GetRandomNumbersInRange(10000000, 99999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}

	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)

	token := GetLoginToken(t, r, user, loginData)

	// store some recovery options in the db
	account := account.Controller{Db: db, Validator: validatorRef, Logger: logger}

	accountGroup := r.Group(fmt.Sprintf("%v", "/api/v1/account"), middleware.Authorize())
	{
		accountGroup.POST("/add-recovery-email", account.AddRecoveryEmail)
		accountGroup.POST("/recovery-number", account.AddRecoveryPhoneNumber)
		accountGroup.POST("/submit-security-answers", account.AddSecurityAnswers)
	}

	emailReq := models.AddRecoveryEmailRequestModel{Email: fmt.Sprintf("testuser%v@qa.team", currUUID)}
	phoneNumberReq := models.AddRecoveryPhoneNumberRequestModel{PhoneNumber: fmt.Sprintf("+23470%v", utility.GetRandomNumbersInRange(10000000, 99999999))}
	securityQuestionsReq := map[string][]map[string]string{
		"answers": {
			{
				"question_1": "What is your mother's maiden name?",
				"answer_1":   "User's Answer",
			},
			{
				"question_2": "In what city were you born?",
				"answer_2":   "User's Answer",
			},
			{
				"question_3": "What is the name of your first pet?",
				"answer_3":   "User's Answer",
			},
		},
	}

	reqTable := map[string]any{
		"/api/v1/account/add-recovery-email":      emailReq,
		"/api/v1/account/recovery-number":         phoneNumberReq,
		"/api/v1/account/submit-security-answers": securityQuestionsReq,
	}
	for url, req := range reqTable {
		var b bytes.Buffer
		json.NewEncoder(&b).Encode(req)

		req, err := http.NewRequest(http.MethodPost, url, &b)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		AssertStatusCode(t, rr.Code, http.StatusOK)
	}

	tests := []struct {
		Name         string
		RequestBody  models.UpdateRecoveryOptionsRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Update recovery options success",
			RequestBody: models.UpdateRecoveryOptionsRequestModel{
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				PhoneNumber: fmt.Sprintf("+23480%v", utility.GetRandomNumbersInRange(10000000, 99999999)),
				Questions: []map[string]string{
					{
						"question_1": "What is your mother's maiden name?",
						"answer_1":   "New answer",
					},
					{
						"question_2": "In what city were you born?",
						"answer_2":   "New answer",
					},
					{
						"question_3": "What is the name of your first pet?",
						"answer_3":   "New answer",
					},
				},
			},
			ExpectedCode: http.StatusOK,
			Message:      "Recovery options updated",
		},
		{
			Name: "Update recovery email faliure",
			RequestBody: models.UpdateRecoveryOptionsRequestModel{
				Email:       "gy@tt .rizz",
				PhoneNumber: fmt.Sprintf("+23480%v", utility.GetRandomNumbersInRange(10000000, 99999999)),
				Questions: []map[string]string{
					{
						"question_1": "What is your mother's maiden name?",
						"answer_1":   "New answer",
					},
					{
						"question_2": "In what city were you born?",
						"answer_2":   "New answer",
					},
					{
						"question_3": "What is the name of your first pet?",
						"answer_3":   "New answer",
					},
				},
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Invalid recovery options",
		},
		{
			Name: "Update recovery phone number faliure",
			RequestBody: models.UpdateRecoveryOptionsRequestModel{
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				PhoneNumber: "o8onaija4laife",
				Questions: []map[string]string{
					{
						"question_1": "What is your mother's maiden name?",
						"answer_1":   "New answer",
					},
					{
						"question_2": "In what city were you born?",
						"answer_2":   "New answer",
					},
					{
						"question_3": "What is the name of your first pet?",
						"answer_3":   "New answer",
					},
				},
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Invalid recovery options",
		},
		{
			Name: "Update security questions faliure",
			RequestBody: models.UpdateRecoveryOptionsRequestModel{
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				PhoneNumber: fmt.Sprintf("+23481%v", utility.GetRandomNumbersInRange(10000000, 99999999)),
				Questions: []map[string]string{
					{
						"question_1": "What is your mother's maiden name?",
					},
					{
						"question_2": "In what city were you born?",
					},
					{
						"question_3": "What is the name of your first pet?",
					},
				},
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Invalid recovery options",
		},
	}

	for _, test := range tests {
		r := gin.Default()

		accountGroup := r.Group(fmt.Sprintf("%v", "/api/v1/account"), middleware.Authorize())
		{
			accountGroup.PUT("/update-recovery-options", account.UpdateRecoveryOptions)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPut, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := ParseResponse(rr)

			code := int(data["status_code"].(float64))
			AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					AssertResponseMessage(t, message.(string), test.Message)
				} else {
					AssertResponseMessage(t, "", test.Message)
				}
			}
		})
	}
}

func TestDeleteRecoveryOptions(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/account/delete-recovery-options"}
	currUUID := utility.GenerateUUID()

	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+23481%v", utility.GetRandomNumbersInRange(10000000, 99999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}

	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)

	token := GetLoginToken(t, r, user, loginData)

	// store some recovery options in the db
	account := account.Controller{Db: db, Validator: validatorRef, Logger: logger}

	accountGroup := r.Group(fmt.Sprintf("%v", "/api/v1/account"), middleware.Authorize())
	{
		accountGroup.POST("/add-recovery-email", account.AddRecoveryEmail)
		accountGroup.POST("/recovery-number", account.AddRecoveryPhoneNumber)
		accountGroup.POST("/submit-security-answers", account.AddSecurityAnswers)
	}

	emailReq := models.AddRecoveryEmailRequestModel{Email: fmt.Sprintf("testuser%v@qa.team", currUUID)}
	phoneNumberReq := models.AddRecoveryPhoneNumberRequestModel{PhoneNumber: fmt.Sprintf("+23470%v", utility.GetRandomNumbersInRange(10000000, 99999999))}
	securityQuestionsReq := map[string][]map[string]string{
		"answers": {
			{
				"question_1": "What is your mother's maiden name?",
				"answer_1":   "User's Answer",
			},
			{
				"question_2": "In what city were you born?",
				"answer_2":   "User's Answer",
			},
			{
				"question_3": "What is the name of your first pet?",
				"answer_3":   "User's Answer",
			},
		},
	}

	reqTable := map[string]any{
		"/api/v1/account/add-recovery-email":      emailReq,
		"/api/v1/account/recovery-number":         phoneNumberReq,
		"/api/v1/account/submit-security-answers": securityQuestionsReq,
	}
	for url, req := range reqTable {
		var b bytes.Buffer
		json.NewEncoder(&b).Encode(req)

		req, err := http.NewRequest(http.MethodPost, url, &b)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		AssertStatusCode(t, rr.Code, http.StatusOK)
	}

	tests := []struct {
		Name        string
		RequestBody struct {
			Options []string `json:"options"`
		}
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Delete email option success",
			RequestBody: struct {
				Options []string "json:\"options\""
			}{
				Options: []string{"email"},
			},
			ExpectedCode: http.StatusOK,
			Message:      "Recovery options successfully deleted",
		},
		{
			Name: "Delete phone number and security questions",
			RequestBody: struct {
				Options []string "json:\"options\""
			}{
				Options: []string{"phone_number", "security_questions"},
			},
			ExpectedCode: http.StatusOK,
			Message:      "Recovery options successfully deleted",
		},
	}

	for _, test := range tests {
		r := gin.Default()

		accountGroup := r.Group(fmt.Sprintf("%v", "/api/v1/account"), middleware.Authorize())
		{
			accountGroup.PUT("/delete-recovery-options", account.DeleteRecoveryOptions)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPut, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := ParseResponse(rr)

			code := int(data["status_code"].(float64))
			AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					AssertResponseMessage(t, message.(string), test.Message)
				} else {
					AssertResponseMessage(t, "", test.Message)
				}
			}
		})
	}
}
