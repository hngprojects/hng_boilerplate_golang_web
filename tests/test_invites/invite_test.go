package test_invites

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/invite"
	orgController "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAcceptInvite(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	inviteController := &invite.Controller{
		Db:        db,
		Validator: validatorRef,
		Logger:    logger,
	}

	authController := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	// Set up user and obtain token
	currUUID := utility.GenerateUUID()
	email := fmt.Sprintf("testuser%s@qa.team", currUUID)
	userSignUpData := models.CreateUserRequestModel{
		Email:       email,
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
	tst.SignupUser(t, r, authController, userSignUpData)
	token := tst.GetLoginToken(t, r, authController, loginData)
	//create an organisation
	orgReq := models.CreateOrgRequestModel{
		Name:        fmt.Sprintf("Org %v", currUUID),
		Description: "This is a test organisation",
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		State:       "Lagos",
		Country:     "Nigeria",
		Industry:    "Tech",
		Type:        "Public",
		Address:     "No 1, Test Street",
	}

	org := orgController.Controller{Db: db, Validator: validatorRef, Logger: logger}
	org_id := tst.CreateOrganisation(t, r, db, org, orgReq, token)
	// Create invitation token for testing
	invitationToken := utility.GenerateUUID()
	invitation := models.Invitation{
		ID:             utility.GenerateUUID(),
		OrganisationID: org_id,
		Token:          invitationToken,
		UserID:         utility.GenerateUUID(),
		IsValid:        true,
		Email:          email,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}
	db.Postgresql.Create(&invitation)

	tests := []struct {
		Name         string
		Method       string
		URL          string
		RequestBody  interface{}
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Successful Invitation Acceptance (POST)",
			Method:       http.MethodPost,
			URL:          "/api/v1/invite/accept",
			RequestBody:  models.InvitationAcceptReq{InvitationLink: "http://example.com/invite/" + invitationToken},
			ExpectedCode: http.StatusOK,
			Message:      "Invitation accepted successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Invalid Invitation Token (POST)",
			Method:       http.MethodPost,
			URL:          "/api/v1/invite/accept",
			RequestBody:  models.InvitationAcceptReq{InvitationLink: "http://example.com/invite/invalid"},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Invalid or expired invitation link",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Successful Invitation Acceptance (GET)",
			Method:       http.MethodGet,
			URL:          "/api/v1/invite/accept/" + invitationToken,
			ExpectedCode: http.StatusOK,
			Message:      "Invitation accepted successfully",
			Headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Invalid Invitation Token (GET)",
			Method:       http.MethodGet,
			URL:          "/api/v1/invite/accept/invalid",
			ExpectedCode: http.StatusBadRequest,
			Message:      "Invalid or expired invitation link",
			Headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			r := gin.Default()

			r.POST("/api/v1/invite/accept", middleware.Authorize(), inviteController.PostAcceptInvite)
			r.GET("/api/v1/invite/accept/:t", middleware.Authorize(), inviteController.GetAcceptInvite)

			var req *http.Request
			var err error

			if test.Method == http.MethodPost {
				var b bytes.Buffer
				json.NewEncoder(&b).Encode(test.RequestBody)
				req, err = http.NewRequest(test.Method, test.URL, &b)
			} else {
				req, err = http.NewRequest(test.Method, test.URL, nil)
			}

			if err != nil {
				t.Fatal(err)
			}

			for k, v := range test.Headers {
				req.Header.Set(k, v)
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
			}
		})
	}
}
