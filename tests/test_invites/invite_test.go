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