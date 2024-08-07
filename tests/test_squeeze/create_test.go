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
	"github.com/lib/pq"
)

func TestE2ESqueezeUserCreation(t *testing.T) {
	router, _ := SetupSqueezeTestRouter()

	// Test POST /squeeze
	currUUID := utility.GenerateUUID()
	phone := fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999))
	body := models.SqueezeUserReq{
		Email:          fmt.Sprintf("testuser%v@qa.team", currUUID),
		FirstName:      "test",
		LastName:       "user1",
		Phone:          phone,
		Location:       "Lagos",
		JobTitle:       "Software engineering",
		Company:        "Paystack",
		Interests:      []string{"photos", "phones"},
		ReferralSource: "facebook",
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
	phone := fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999))
	body := models.SqueezeUserReq{
		Email:          fmt.Sprintf("testuser%v@qa", currUUID),
		FirstName:      "test",
		LastName:       "user1",
		Phone:          phone,
		Location:       "Lagos",
		JobTitle:       "Software engineering",
		Company:        "Paystack",
		Interests:      []string{"photos", "phones"},
		ReferralSource: "facebook",
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
	phone := fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999))

	db := squeezeController.Db.Postgresql
	db.Create(&models.SqueezeUser{
		ID:             currUUID,
		Email:          fmt.Sprintf("testuser%v@qa.team", currUUID),
		FirstName:      "test",
		LastName:       "user",
		Phone:          phone,
		Location:       "Lagos",
		JobTitle:       "Software engineering",
		Company:        "Paystack",
		Interests:      pq.StringArray{"photos", "phones"},
		ReferralSource: "facebook",
	})

	body := models.SqueezeUserReq{
		Email:          fmt.Sprintf("testuser%v@qa.team", currUUID),
		FirstName:      "test",
		LastName:       "user",
		Phone:          "09034017724",
		Location:       "Lagos",
		JobTitle:       "Software engineering",
		Company:        "Paystack",
		Interests:      []string{"photos", "phones"},
		ReferralSource: "facebook",
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

func TestCreateSqueezeUser_CheckDuplicatePhoneNumber(t *testing.T) {
	router, squeezeController := SetupSqueezeTestRouter()

	currUUID := utility.GenerateUUID()
	

	db := squeezeController.Db.Postgresql
	db.Create(&models.SqueezeUser{
		ID:             currUUID,
		Email:          fmt.Sprintf("testuser%v@qa.team", currUUID),
		FirstName:      "test",
		LastName:       "user",
		Phone:          "0802879201723",
		Location:       "Lagos",
		JobTitle:       "Software engineering",
		Company:        "Paystack",
		Interests:      pq.StringArray{"photos", "phones"},
		ReferralSource: "facebook",
	})

	body := models.SqueezeUserReq{
		Email:          fmt.Sprintf("testuser%v@qc.team", currUUID),
		FirstName:      "test",
		LastName:       "user",
		Phone:          "0802879201723",
		Location:       "Lagos",
		JobTitle:       "Software engineering",
		Company:        "Paystack",
		Interests:      []string{"photos", "phones"},
		ReferralSource: "facebook",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/squeeze", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	response := tst.ParseResponse(resp)
	tst.AssertStatusCode(t, resp.Code, http.StatusBadRequest)
	tst.AssertResponseMessage(t, response["message"].(string), "user already exists with the given phone")
}
