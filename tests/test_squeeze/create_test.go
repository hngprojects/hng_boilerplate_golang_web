package test_squeeze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestE2ESqueezeUserCreation(t *testing.T) {
	router, _ := SetupSqueezeTestRouter()

	// Test POST /squeeze
	currUUID := utility.GenerateUUID()
	body := models.SqueezeUserReq{
		FirstName: "test",
		LastName:  "user1",
		Email:     fmt.Sprintf("testuser%v@qa.team", currUUID),
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/api/v1/squeeze", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	tst.AssertStatusCode(t, resp.Code, http.StatusCreated)

	response := tst.ParseResponse(resp)
	tst.AssertResponseMessage(t, response["message"].(string), "your request has been received. you will get a template shortly")

}

func TestCreateSqueezeUser_ValidateEmail(t *testing.T) {
	router, _ := SetupSqueezeTestRouter()

	currUUID := utility.GenerateUUID()
	body := models.SqueezeUserReq{
		FirstName: "test",
		LastName:  "user1",
		Email:     fmt.Sprintf("testuser%v@qa", currUUID),
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/squeeze", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := tst.ParseResponse(resp)
	tst.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
	tst.AssertResponseMessage(t, response["message"].(string), "Validation failed")
}

func TestCreateSqueezeUser_CheckDuplicateEmail(t *testing.T) {
	router, squeezeController := SetupSqueezeTestRouter()

	currUUID := utility.GenerateUUID()

	db := squeezeController.Db.Postgresql
	db.Create(&models.SqueezeUser{
		ID:        currUUID,
		FirstName: "test",
		LastName:  "user",
		Email:     fmt.Sprintf("testuser%v@qa.team", currUUID),
	})

	body := models.SqueezeUserReq{
		FirstName: "test",
		LastName:  "user",
		Email:     fmt.Sprintf("testuser%v@qa.team", currUUID),
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/squeeze", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := tst.ParseResponse(resp)
	tst.AssertStatusCode(t, resp.Code, http.StatusBadRequest)
	tst.AssertResponseMessage(t, response["message"].(string), "user already exists with the given email")
}
