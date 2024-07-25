package test_waitlist

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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/waitlist"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestWailistSignup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := tests.Setup()
	validate := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	testEmail := fmt.Sprintf("testuser%v@qa.team", currUUID)

	ttests := []struct {
		Name            string
		Request         models.CreateWaitlistUserRequest
		ExpectedCode    int
		ExpectedMessage string
	}{
		{
			Name: "user can signup on waitlist",
			Request: models.CreateWaitlistUserRequest{
				Name:  "Tester",
				Email: testEmail,
			},
			ExpectedCode:    http.StatusCreated,
			ExpectedMessage: "waitlist signup successful",
		},
		{
			Name: "user can not signup with duplicate email",
			Request: models.CreateWaitlistUserRequest{
				Name:  "Tester",
				Email: testEmail,
			},
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "waitlist user exists",
		},
	}

	wc := waitlist.Controller{DB: db, Logger: logger, Validator: validate}

	for _, tt := range ttests {
		r := gin.Default()

		r.POST("/api/v1/waitlist", wc.Create)

		t.Run(tt.Name, func(t *testing.T) {
			var buf bytes.Buffer

			err := json.NewEncoder(&buf).Encode(tt.Request)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/waitlist", &buf)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			hr := httptest.NewRecorder()
			r.ServeHTTP(hr, req)

			tests.AssertStatusCode(t, hr.Code, tt.ExpectedCode)

			data := tests.ParseResponse(hr)

			if tt.ExpectedMessage != "" {
				message := data["message"]
				if message != nil {
					tests.AssertResponseMessage(t, message.(string), tt.ExpectedMessage)
				} else {
					tests.AssertResponseMessage(t, "", tt.ExpectedMessage)
				}
			}
		})
	}
}
