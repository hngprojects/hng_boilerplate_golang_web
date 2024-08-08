package test_billing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/billing"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Initialise(currUUID string, t *testing.T, r *gin.Engine, db *storage.Database, user auth.Controller, Billing billing.Controller, status bool) (string, string) {
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

	tst.SignupUser(t, r, user, userSignUpData, status)

	token := tst.GetLoginToken(t, r, user, loginData)

	BillingCreationData := models.CreateBillingRequest{
		Name:  fmt.Sprintf("Billing Name %s", currUUID),
		Price: float64(utility.GetRandomNumbersInRange(0, 10000_00)),
	}

	BillingID := CreateBilling(t, r, db, Billing, BillingCreationData, token)

	return BillingID, token
}

func CreateBilling(t *testing.T, r *gin.Engine, db *storage.Database, Billing billing.Controller, BillingData models.CreateBillingRequest, token string) string {
	var (
		BillingPath = "/api/v1/billing-plans"
		BillingURI  = url.URL{Path: BillingPath}
	)
	BillingUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))
	{
		BillingUrl.POST("/billing-plans", Billing.CreateBilling)
	}
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(BillingData)
	req, err := http.NewRequest(http.MethodPost, BillingURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	//get the response
	data := tst.ParseResponse(rr)
	dataM := data["data"].(map[string]interface{})
	BillingID := dataM["id"].(string)
	return BillingID
}
