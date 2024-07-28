package test_organisation

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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestOrganizationCreate(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/organisations"}
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
		RequestBody  models.CreateOrgRequestModel
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Successful organisation register",
			RequestBody: models.CreateOrgRequestModel{
				Name:        fmt.Sprintf("Org %v", currUUID),
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				Description: "Some random description about vibranium",
				State:       "test",
				Industry:    "user",
				Type:        "type1",
				Address:     "wakanda land",
				Country:     "wakanda",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "organisation created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "details already exist",
			RequestBody: models.CreateOrgRequestModel{
				Name:        fmt.Sprintf("Org %v", currUUID),
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				Description: "Some random description about vibranium",
				State:       "test",
				Industry:    "user",
				Type:        "type1",
				Address:     "wakanda land",
				Country:     "wakanda",
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "organization already exists with the given email",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "invalid email",
			RequestBody: models.CreateOrgRequestModel{
				Name:        fmt.Sprintf("Org %v", utility.GenerateUUID()),
				Email:       "someRandEmail",
				Description: "Some random description about vibranium",
				State:       "test",
				Industry:    "user",
				Type:        "type1",
				Address:     "wakanda land",
				Country:     "wakanda",
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "email address is invalid",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "Validation failed",
			RequestBody: models.CreateOrgRequestModel{
				Name:        fmt.Sprintf("Org %v", currUUID),
				Description: "Some random description about vibranium",
				State:       "test",
				Industry:    "user",
				Type:        "type1",
				Address:     "wakanda land",
				Country:     "wakanda",
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "User unauthorized",
			RequestBody: models.CreateOrgRequestModel{
				Name:        fmt.Sprintf("Org %v", currUUID),
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				Description: "Some random description about vibranium",
				State:       "test",
				Industry:    "user",
				Type:        "type1",
				Address:     "wakanda land",
				Country:     "wakanda",
			},
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
		},
	}

	org := organisation.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		orgUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql))
		{
			orgUrl.POST("/organisations", org.CreateOrganisation)

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
