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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/invite"
	orgController "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestCreateInvite(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/invite/create"}
	currUUID := utility.GenerateUUID()

	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser" + currUUID + "@qa.team"),
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

	inviteController := &invite.Controller{
		Db:        db,
		Validator: validatorRef,
		Logger:    logger,
	}

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)
	token := GetLoginToken(t, r, user, loginData)

	tests := []struct {
		Name         string
		RequestBody  models.InvitationCreateReq
		SetupFunc    func()
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Successful Invitation Creation",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: "valid-org-id",
				Email:          "testinvitee@qa.team",
			},
			SetupFunc: func() {
				// Mock necessary services and data for a successful invitation creation
				organisation.MockCheckOrgExists(true, nil)
				invite.MockCheckUserIsAdmin(true, nil)
				invite.MockGenerateInvitationToken("valid-token", nil)
				invite.MockSaveInvitation(nil)
				invite.MockGenerateInvitationLink("http://localhost:8019/invite/accept?token=valid-token", nil)
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Invitation created successfully",
		},
		{
			Name: "Invalid Email Format",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: "valid-org-id",
				Email:          "invalid-email",
			},
			SetupFunc: func() {
				// No setup needed for invalid email test
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Invalid email format",
		},
		{
			Name: "Invalid Organisation ID",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: "invalid-org-id",
				Email:          "testinvitee@qa.team",
			},
			SetupFunc: func() {
				// Mock necessary services and data for invalid organisation ID
				organisation.MockCheckOrgExists(false, fmt.Errorf("organisation not found"))
			},
			ExpectedCode: http.StatusNotFound,
			Message:      "Invalid Organisation ID",
		},
		{
			Name: "User Not Admin",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: "valid-org-id",
				Email:          "testinvitee@qa.team",
			},
			SetupFunc: func() {
				// Mock necessary services and data for user not being an admin
				organisation.MockCheckOrgExists(true, nil)
				invite.MockCheckUserIsAdmin(false, nil)
			},
			ExpectedCode: http.StatusForbidden,
			Message:      "User is not an admin of the organisation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.SetupFunc != nil {
				tt.SetupFunc()
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, requestURI.String(), utility.ToJSONBuffer(tt.RequestBody))
			req.Header.Set("Authorization", "Bearer "+token)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
			if tt.ExpectedCode == http.StatusCreated {
				assert.Contains(t, w.Body.String(), tt.Message)
			} else {
				assert.Contains(t, w.Body.String(), tt.Message)
			}
		})
	}
}

func TestPostInvite(t *testing.T) {

	logger := Setup()
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

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)
	token := GetLoginToken(t, r, user, loginData)

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
	org_id := CreateOrganisation(t, r, org, orgReq, token)

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

			AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := ParseResponse(rr)

			code := int(data["code"].(float64))
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
