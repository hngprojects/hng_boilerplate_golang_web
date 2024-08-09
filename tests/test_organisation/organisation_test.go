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

		orgUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))
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

func TestGetOrganisation(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	org := organisation.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	orgID, token := initialise(currUUID, t, r, db, user, org, false)

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
			Name:         "Successful retrieval of Organisation",
			OrgID:        orgID,
			ExpectedCode: http.StatusOK,
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
	}

	orgUrl := r.Group("/api/v1", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))
	{
		orgUrl.GET("/organisations/:org_id", org.GetOrganisation)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/organisations/%s", test.OrgID), nil)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.ExpectedCode {
				tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)
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
				tst.AssertStatusCode(t, int(code), test.ExpectedCode)
			}

			if test.Message != "" {
				message, ok := data["message"].(string)
				if !ok {
					tst.AssertResponseMessage(t, message, test.Message)
				}
				if message != test.Message {
					tst.AssertResponseMessage(t, message, test.Message)
				}
			}
		})
	}
}

func TestOrganisationUpdate(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	org := organisation.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	orgID, token := initialise(currUUID, t, r, db, user, org, false)

	tests := []struct {
		Name         string
		OrgID        string
		RequestBody  models.UpdateOrgRequestModel
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:  "Successful organisation update",
			OrgID: orgID,
			RequestBody: models.UpdateOrgRequestModel{
				Name:        fmt.Sprintf("Org %v", currUUID),
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				Description: "Some random description about vibranium",
				State:       "test",
				Industry:    "user",
				Type:        "type1",
				Address:     "wakanda land",
				Country:     "wakanda",
			},
			ExpectedCode: http.StatusOK,
			Message:      "organisation updated successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name:  "organisation not found",
			OrgID: utility.GenerateUUID(),
			RequestBody: models.UpdateOrgRequestModel{
				Name:        fmt.Sprintf("Org %v", currUUID),
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				Description: "Some random description about vibranium",
				State:       "test",
				Industry:    "user",
				Type:        "type1",
				Address:     "wakanda land",
				Country:     "wakanda",
			},
			ExpectedCode: http.StatusNotFound,
			Message:      "organisation not found",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name:  "invalid organisation id format",
			OrgID: "invalid-id-erttt",
			RequestBody: models.UpdateOrgRequestModel{
				Name:        fmt.Sprintf("Org %v", utility.GenerateUUID()),
				State:       "test",
				Industry:    "user",
				Type:        "type1",
				Address:     "wakanda land",
				Country:     "wakanda",
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "invalid organisation id format",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name:  "User unauthorized",
			OrgID: orgID,
			RequestBody: models.UpdateOrgRequestModel{
				Name:        fmt.Sprintf("Org %v", currUUID),
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				Description: "Some random description about vibranium",
				Address:     "wakanda land",
				Country:     "wakanda",
			},
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
		},
	}

	orgUrl := r.Group("/api/v1", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))
	{
		orgUrl.PATCH("/organisations/:org_id", org.UpdateOrganisation)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/organisations/%s", test.OrgID), &b)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.ExpectedCode {
				tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)
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
				tst.AssertStatusCode(t, int(code), test.ExpectedCode)
			}

			if test.Message != "" {
				message, ok := data["message"].(string)
				if !ok {
					tst.AssertResponseMessage(t, message, test.Message)
				}
				if message != test.Message {
					tst.AssertResponseMessage(t, message, test.Message)
				}
			}
		})
	}
}

func TestOrganisationDelete(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	org := organisation.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	orgID, token := initialise(currUUID, t, r, db, user, org, false)

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

	orgUrl := r.Group("/api/v1", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))
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

func TestGetUsersInOrg(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	org := organisation.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	orgID, token := initialise(currUUID, t, r, db, user, org, true)
	tests := []struct {
		Name         string
		OrgID        string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Successful retrieval of users in organisation",
			OrgID:        orgID,
			ExpectedCode: http.StatusOK,
			Message:      "users retrieved successfully",
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
			Name:         "User Not Authorized to Delete Organization",
			OrgID:        orgID,
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		orgUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql,models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))
		{
			orgUrl.GET("/organisations/:org_id/users", org.GetUsersInOrganisation)
		}

		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet,fmt.Sprintf("/api/v1/organisations/%s/users", test.OrgID), nil)
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

