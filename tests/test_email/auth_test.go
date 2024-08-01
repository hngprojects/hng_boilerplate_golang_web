package test_auth

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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

// test welcome email
func TestUserSignup(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/auth/users/signup"}
	currUUID := utility.GenerateUUID()

	tests := []struct {
		Name         string
		RequestBody  models.CreateUserRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Welcome Email sent",
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
		},
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		r.POST("/api/v1/auth/users/signup", auth.CreateUser)

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

			}

		})

	}

}
