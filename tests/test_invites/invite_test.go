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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/invite"
	orgController "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestCreateInvite(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/invite/create"}
	currUUID := utility.GenerateUUID()
	email := fmt.Sprintf("testuser" + currUUID + "@qa.team")

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

	inviteController := &invite.Controller{
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
	org_id := tst.CreateOrganisation(t, r, db, org, orgReq, token)

	tests := []struct {
		Name         string
		RequestBody  models.InvitationCreateReq
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Successful Invitation Creation",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: org_id,
				Email:          orgReq.Email,
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Invitation created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Invalid Email Format",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: org_id,
				Email:          "micah",
			},
			ExpectedCode: http.StatusBadRequest,
			// Message:      "Invalid email format",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Empty Email Field",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: org_id,
				Email:          "",
			},
			ExpectedCode: http.StatusBadRequest,
			// Message:      "Failed to parse request body",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Missing Organisation ID",
			RequestBody: models.InvitationCreateReq{
				Email: email,
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Validation failed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Invalid Organisation ID",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: "0190d9a1-e05e-787d-85ee-bd91a61c6da0",
				Email:          email,
			},
			ExpectedCode: http.StatusNotFound,
			Message:      "Invalid Organisation ID",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		inviteURL := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql))
		{
			inviteURL.POST("/invite/create", inviteController.CreateInvite)

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

			}

		})

	}
}

// func TestPostInvite(t *testing.T) {

// 	logger := Setup()
// 	gin.SetMode(gin.TestMode)
// 	validatorRef := validator.New()
// 	db := storage.Connection()
// 	requestURI := url.URL{Path: "/api/v1/organisations/send-invite"}
// 	currUUID := utility.GenerateUUID()

// 	userSignUpData := models.CreateUserRequestModel{
// 		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
// 		PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
// 		FirstName:   "test",
// 		LastName:    "user",
// 		Password:    "password",
// 		UserName:    fmt.Sprintf("test_username%v", currUUID),
// 	}
// 	loginData := models.LoginRequestModel{
// 		Email:    userSignUpData.Email,
// 		Password: userSignUpData.Password,
// 	}

// 	controller := &invite.Controller{
// 		Db:        db,
// 		Validator: validatorRef,
// 		Logger:    logger,
// 	}

// 	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
// 	r := gin.Default()
// 	SignupUser(t, r, user, userSignUpData)
// 	token := GetLoginToken(t, r, user, loginData)

// 	//create an organisation
// 	orgReq := models.CreateOrgRequestModel{
// 		Name:        fmt.Sprintf("Org %v", currUUID),
// 		Description: "This is a test organisation",
// 		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
// 		State:       "Lagos",
// 		Country:     "Nigeria",
// 		Industry:    "Tech",
// 		Type:        "Public",
// 		Address:     "No 1, Test Street",
// 	}

// 	org := orgController.Controller{Db: db, Validator: validatorRef, Logger: logger}
// 	org_id := CreateOrganisation(t, r, org, orgReq, token)

// 	tests := []struct {
// 		Name         string
// 		RequestBody  models.InvitationRequest
// 		ExpectedCode int
// 		Message      string
// 	}{
// 		{
// 			Name: "Successful invite",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  org_id,
// 				Emails: []string{"micahshallom@gmail.com"},
// 			},
// 			ExpectedCode: http.StatusCreated,
// 			Message:      "Invitation(s) sent successfully",
// 		},
// 		{
// 			Name: "Invalid org_id format",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusUnprocessableEntity,
// 			Message:      "Invalid org_id format",
// 		},
// 		{
// 			Name: "Missing OrgID",
// 			RequestBody: models.InvitationRequest{
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusBadRequest,
// 			Message:      "OrgID is required",
// 		},
// 		{
// 			Name: "Empty Emails Array",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{},
// 			},
// 			ExpectedCode: http.StatusBadRequest,
// 			Message:      "Emails array cannot be empty",
// 		},
// 		{
// 			Name: "Non-existent OrgID",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "non-existent-org-id",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusNotFound,
// 			Message:      "Organization not found",
// 		},
// 		{
// 			Name: "Non-member User Sending Invite",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusForbidden,
// 			Message:      "User is not a member of the organization",
// 		},
// 		{
// 			Name: "Invalid Email Format",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"invalid-email"},
// 			},
// 			ExpectedCode: http.StatusBadRequest,
// 			Message:      "Invalid email format",
// 		},
// 		{
// 			Name: "Duplicate Emails",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"test@example.com", "test@example.com"},
// 			},
// 			ExpectedCode: http.StatusConflict,
// 			Message:      "Duplicate emails found",
// 		},
// 		{
// 			Name: "Exceeding Email Limit",
// 			RequestBody: models.InvitationRequest{
// 				OrgID: "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{
// 					"email1@example.com",
// 					"email2@example.com",
// 					// ... add more emails to exceed the limit
// 				},
// 			},
// 			ExpectedCode: http.StatusBadRequest,
// 			Message:      "Email limit exceeded",
// 		},
// 		{
// 			Name: "Unauthorized User",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusUnauthorized,
// 			Message:      "Unauthorized",
// 		},
// 		{
// 			Name: "Expired Token",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusUnauthorized,
// 			Message:      "Token has expired",
// 		},
// 		{
// 			Name: "Database Error Handling",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusInternalServerError,
// 			Message:      "Internal server error",
// 		},
// 	}

// 	for _, test := range tests {
// 		r := gin.Default()
// 		r.POST("/api/v1/invite", controller.PostInvite)

// 		t.Run(test.Name, func(t *testing.T) {
// 			var b bytes.Buffer
// 			json.NewEncoder(&b).Encode(test.RequestBody)

// 			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Authorization", "Bearer "+token)

// 			rr := httptest.NewRecorder()
// 			r.ServeHTTP(rr, req)

// 			AssertStatusCode(t, rr.Code, test.ExpectedCode)

// 			data := ParseResponse(rr)

// 			code := int(data["code"].(float64))
// 			AssertStatusCode(t, code, test.ExpectedCode)

// 			if test.Message != "" {
// 				message := data["message"]
// 				if message != nil {
// 					AssertResponseMessage(t, message.(string), test.Message)
// 				} else {
// 					AssertResponseMessage(t, "", test.Message)
// 				}

// 			}

// 		})
// 	}
// }
