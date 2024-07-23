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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestSignup(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/users/signup"}
	currUUID := utility.GenerateUUID()

	tests := []struct {
		Name         string
		RequestBody  models.CreateUserRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Successful user register",
			RequestBody: models.CreateUserRequestModel{
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusCreated,
			Message:      "user created successfully",
		}, {
			Name: "details already exist",
			RequestBody: models.CreateUserRequestModel{
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "user already exists with the given email",
		}, {
			Name: "invalid email",
			RequestBody: models.CreateUserRequestModel{
				Email:       "emailtest",
				PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "email address is invalid",
		}, {
			Name: "Validation failed",
			RequestBody: models.CreateUserRequestModel{
				PhoneNumber: "090909",
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
		},
	}

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		r.POST("/api/v1/users/signup", user.CreateUser)

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

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

func TestLogin(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	var (
		loginPath      = "/api/v1/users/login"
		loginURI       = url.URL{Path: loginPath}
		currUUID       = utility.GenerateUUID()
		userSignUpData = models.CreateUserRequestModel{
			Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
			PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			FirstName:   "test",
			LastName:    "user",
			Password:    "password",
			UserName:    fmt.Sprintf("test_username%v", currUUID),
		}
	)

	tests := []struct {
		Name         string
		RequestBody  models.LoginRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "OK email login successful",
			RequestBody: models.LoginRequestModel{
				Email:    userSignUpData.Email,
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusOK,
			Message:      "user login successfully",
		}, {
			Name:         "password not provided",
			RequestBody:  models.LoginRequestModel{},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "username or phone or email not provided",
			RequestBody: models.LoginRequestModel{
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "email does not exist",
			RequestBody: models.LoginRequestModel{
				Email:    utility.GenerateUUID(),
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "incorrect password",
			RequestBody: models.LoginRequestModel{
				Email:    fmt.Sprintf("testuser%v@qa.team", currUUID),
				Password: "incorrect",
			},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	r.POST(loginPath, user.LoginUser)

	SignupUser(t, r, user, userSignUpData)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, loginURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

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
