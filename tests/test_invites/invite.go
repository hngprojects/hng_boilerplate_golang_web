package test_invites

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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/invite"
	orgController "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestPostInvite(t *testing.T) {

	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/organisations/send-invite"}
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

	controller := &invite.Controller{
		Db:        db,
		Validator: validatorRef,
		Logger:    logger,
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData)
	token := tst.GetLoginToken(t, r, auth, loginData)

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
	org_id := tst.CreateOrganisation(t, r,db, org, orgReq, token)

	tests := []struct {
		Name         string
		RequestBody  models.InvitationRequest
		SetupFunc    func()
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Successful invite",
			RequestBody: models.InvitationRequest{
				OrgID:  org_id,
				Emails: []string{"micahshallom@gmail.com"},
			},
			SetupFunc:    func() {},
			ExpectedCode: http.StatusCreated,
			Message:      "Invitation(s) sent successfully",
		},
		{
			Name: "Invalid org_id format",
			RequestBody: models.InvitationRequest{
				OrgID:  "0190d5be-e185-72ef-b74a-0",
				Emails: []string{"test@example.com"},
			},
			SetupFunc:    func() {},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Invalid org_id format",
		},
		{
			Name: "Missing OrgID",
			RequestBody: models.InvitationRequest{
				Emails: []string{"test@example.com"},
			},
			SetupFunc:    func() {},
			ExpectedCode: http.StatusBadRequest,
			Message:      "OrgID is required",
		},
		{
			Name: "Empty Emails Array",
			RequestBody: models.InvitationRequest{
				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
				Emails: []string{},
			},
			SetupFunc:    func() {},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Emails array cannot be empty",
		},
		{
			Name: "Non-existent OrgID",
			RequestBody: models.InvitationRequest{
				OrgID:  "non-existent-org-id",
				Emails: []string{"test@example.com"},
			},
			SetupFunc: func() {

			},
			ExpectedCode: http.StatusNotFound,
			Message:      "Organization not found",
		},
		{
			Name: "Non-member User Sending Invite",
			RequestBody: models.InvitationRequest{
				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
				Emails: []string{"test@example.com"},
			},
			SetupFunc: func() {

			},
			ExpectedCode: http.StatusForbidden,
			Message:      "User is not a member of the organization",
		},
		{
			Name: "Invalid Email Format",
			RequestBody: models.InvitationRequest{
				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
				Emails: []string{"invalid-email"},
			},
			SetupFunc:    func() {},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Invalid email format",
		},
		{
			Name: "Duplicate Emails",
			RequestBody: models.InvitationRequest{
				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
				Emails: []string{"test@example.com", "test@example.com"},
			},
			SetupFunc: func() {

			},
			ExpectedCode: http.StatusConflict,
			Message:      "Duplicate emails found",
		},
		{
			Name: "Exceeding Email Limit",
			RequestBody: models.InvitationRequest{
				OrgID: "0190d5be-e185-72ef-b74a-0c9fce0e2328",
				Emails: []string{
					"email1@example.com",
					"email2@example.com",
					// ... add more emails to exceed the limit
				},
			},
			SetupFunc: func() {

			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Email limit exceeded",
		},
		{
			Name: "Unauthorized User",
			RequestBody: models.InvitationRequest{
				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
				Emails: []string{"test@example.com"},
			},
			SetupFunc:    func() {},
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Unauthorized",
		},
		{
			Name: "Expired Token",
			RequestBody: models.InvitationRequest{
				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
				Emails: []string{"test@example.com"},
			},
			SetupFunc:    func() {},
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token has expired",
		},
		{
			Name: "Database Error Handling",
			RequestBody: models.InvitationRequest{
				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
				Emails: []string{"test@example.com"},
			},
			SetupFunc: func() {

			},
			ExpectedCode: http.StatusInternalServerError,
			Message:      "Internal server error",
		},
	}

	for _, test := range tests {
		r := gin.Default()
		r.POST("/api/v1/invite", controller.PostInvite)

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

			code := int(data["code"].(float64))
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
