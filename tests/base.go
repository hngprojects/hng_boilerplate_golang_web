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

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models/migrations"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models/seed"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Setup() *utility.Logger {
	logger := utility.NewLogger()
	config := config.Setup(logger, "../../app")

	postgresql.ConnectToDatabase(logger, config.TestDatabase)
	db := storage.Connection()
	if config.TestDatabase.Migrate {
		migrations.RunAllMigrations(db)
		// fix correct seed call
		seed.SeedDatabase(db.Postgresql)
	}
	return logger
}

func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}

func AssertStatusCode(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("handler returned wrong status code: got status %d expected status %d", got, expected)
	}
}

func AssertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("handler returned wrong message: got message: %q expected: %q", got, expected)
	}
}
func AssertBool(t *testing.T, got, expected bool) {
	if got != expected {
		t.Errorf("handler returned wrong boolean: got %v expected %v", got, expected)
	}
}

func AssertValidationError(t *testing.T, response map[string]interface{}, field string, expectedMessage string) {
	errors, ok := response["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'error' field in response")
	}

	errorMessage, exists := errors[field]
	if !exists {
		t.Fatalf("expected validation error message for field '%s'", field)
	}

	if errorMessage != expectedMessage {
		t.Errorf("unexpected error message for field '%s': got %v, want %v", field, errorMessage, expectedMessage)
	}
}

// helper to signup a user
func SignupUser(t *testing.T, r *gin.Engine, auth auth.Controller, userSignUpData models.CreateUserRequestModel) {
	var (
		signupPath = "/api/v1/auth/users/signup"
		signupURI  = url.URL{Path: signupPath}
	)
	r.POST(signupPath, auth.CreateUser)
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(userSignUpData)
	req, err := http.NewRequest(http.MethodPost, signupURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
}

// help to fetch user token

func GetLoginToken(t *testing.T, r *gin.Engine, auth auth.Controller, loginData models.LoginRequestModel) string {
	var (
		loginPath = "/api/v1/auth/login"
		loginURI  = url.URL{Path: loginPath}
	)
	r.POST(loginPath, auth.LoginUser)
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(loginData)
	req, err := http.NewRequest(http.MethodPost, loginURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		return ""
	}

	data := ParseResponse(rr)
	dataM := data["data"].(map[string]interface{})
	token := dataM["access_token"].(string)

	return token
}

// helper to create an organisation
func CreateOrganisation(t *testing.T, r *gin.Engine, db *storage.Database, org organisation.Controller, orgData models.CreateOrgRequestModel, token string) string {
	var (
		orgPath = "/api/v1/organisations"
		orgURI  = url.URL{Path: orgPath}
	)
	orgUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql))
	{
		orgUrl.POST("/organisations", org.CreateOrganisation)
	}
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(orgData)
	req, err := http.NewRequest(http.MethodPost, orgURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	//get the response
	data := ParseResponse(rr)
	dataM := data["data"].(map[string]interface{})
	orgID := dataM["id"].(string)
	return orgID
}

// helper to create an invite
// func CreateInvite(t *testing.T, r *gin.Engine, invite invite.Controller, inviteData models.InvitationCreateReq, token string) map[string]interface{} {
// 	var (
// 		invitePath = "/api/v1/invite/create"
// 		inviteURI  = url.URL{Path: invitePath}
// 	)
// 	inviteUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize())
// 	{
// 		inviteUrl.POST("/organisations/send-invite", invite.CreateInvite)
// 	}
// 	var b bytes.Buffer
// 	json.NewEncoder(&b).Encode(inviteData)
// 	req, err := http.NewRequest(http.MethodPost, inviteURI.String(), &b)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+token)

// 	rr := httptest.NewRecorder()
// 	r.ServeHTTP(rr, req)

// 	//get the response
// 	data := ParseResponse(rr)
// 	return data
// }
