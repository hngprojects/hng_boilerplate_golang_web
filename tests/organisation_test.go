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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestOrganisationCreate(t *testing.T) {
	logger := Setup()
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

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)

	token := GetLoginToken(t, r, user, loginData)

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

		orgUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize())
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

func TestOrganisationDelete(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
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

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	org := organisation.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)

	token := GetLoginToken(t, r, user, loginData)

	organisationCreationData := models.CreateOrgRequestModel{
		Name:        fmt.Sprintf("Org %v", currUUID),
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		Description: "Some random description about vibranium",
		State:       "test",
		Industry:    "user",
		Type:        "type1",
		Address:     "wakanda land",
		Country:     "wakanda",
	}

	orgID := GetOrgId(t, r, org, organisationCreationData, token)

	tests := []struct {
		Name         string
		OrgID        string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Invalid Organisation ID Format",
			OrgID:        "invalid-id-erttt",
			ExpectedCode: http.StatusBadRequest,
			Message:      "invalid organisation id format",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Successful Deletion of Organisation",
			OrgID:        orgID,
			ExpectedCode: http.StatusNoContent,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Organisation Not Found",
			OrgID:        utility.GenerateUUID(),
			ExpectedCode: http.StatusNotFound,
			Message:      "organisation not found",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "User Not Authorized to Delete Organization",
			OrgID:        orgID,
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	orgUrl := r.Group("/api/v1", middleware.Authorize())
	{
		orgUrl.DELETE("/organisations/:org_id", org.DeleteOrganisation)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/organisations/%s", test.OrgID), nil)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.ExpectedCode {
				t.Errorf("Expected status code %d, got %d", test.ExpectedCode, rr.Code)
			}

			// For 204 No Content, we don't try to parse the body
			if rr.Code == http.StatusNoContent {
				// Success case with no content, just return
				return
			}

			var data map[string]interface{}
			if err := json.NewDecoder(rr.Body).Decode(&data); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			code, ok := data["status_code"].(float64)
			if !ok {
				t.Fatalf("Expected status_code to be float64, got %T", data["status_code"])
			}
			if int(code) != test.ExpectedCode {
				t.Errorf("Expected status code %d, got %d", test.ExpectedCode, int(code))
			}

			if test.Message != "" {
				message, ok := data["message"].(string)
				if !ok {
					t.Fatalf("Expected message to be string, got %T", data["message"])
				}
				if message != test.Message {
					t.Errorf("Expected message '%s', got '%s'", test.Message, message)
				}
			}
		})
	}
}
