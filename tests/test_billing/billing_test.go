package test_billing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/billing"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestBillingCreate(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	billing := billing.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	_, token := Initialise(currUUID, t, r, db, user, billing, true)

	tests := []struct {
		Name         string
		RequestBody  models.CreateBillingRequest
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Successful billing created",
			RequestBody: models.CreateBillingRequest{
				Name:  fmt.Sprintf("Billing Name %s", utility.GenerateUUID()),
				Price: float64(utility.GetRandomNumbersInRange(0, 10000_00)),
			},
			ExpectedCode: http.StatusCreated,
			Message:      "billing created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Validation failed",
			RequestBody: models.CreateBillingRequest{
				Price: float64(utility.GetRandomNumbersInRange(0, 10000_00)),
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "User unauthorized",
			RequestBody: models.CreateBillingRequest{
				Name:  fmt.Sprintf("Billing Name %s", utility.GenerateUUID()),
				Price: float64(utility.GetRandomNumbersInRange(0, 10000_00)),
			},
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		billingUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin))
		{
			billingUrl.POST("/billing-plans", billing.CreateBilling)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/billing-plans", &b)
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

func TestBillingDelete(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	billing := billing.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	billingId, token := Initialise(currUUID, t, r, db, user, billing, true)

	tests := []struct {
		Name         string
		billingID    string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Successful Deletion of billing",
			billingID:    billingId,
			ExpectedCode: http.StatusNoContent,
			Message:      "billing successfully deleted",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Invalid billing ID Format",
			billingID:    "invalid-id-erttt",
			ExpectedCode: http.StatusBadRequest,
			Message:      "invalid billing id format",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "billing Not Found",
			billingID:    utility.GenerateUUID(),
			ExpectedCode: http.StatusNotFound,
			Message:      "billing not found",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "User Not Authorized to Delete billing",
			billingID:    billingId,
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	
	for _, test := range tests {
		r := gin.Default()

		billingUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin))
		{
			billingUrl.DELETE("/billing-plans/:id", billing.DeleteBilling)
		}

		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/billing-plans/%s", test.billingID), nil)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code == http.StatusNoContent {
				return
			}

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

func TestGetbillingById(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	billing := billing.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	billingId, _ := Initialise(currUUID, t, r, db, user, billing, true)

	tests := []struct {
		Name         string
		billingID    string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Successful Retrieval of billing",
			billingID:    billingId,
			ExpectedCode: http.StatusOK,
			Message:      "billing retrieved successfully",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name:         "Invalid billing ID Format",
			billingID:    "invalid-id-erttt",
			ExpectedCode: http.StatusBadRequest,
			Message:      "invalid billing id format",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name:         "billing Not Found",
			billingID:    utility.GenerateUUID(),
			ExpectedCode: http.StatusNotFound,
			Message:      "billing not found",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		billingUrl := r.Group(fmt.Sprintf("%v", "/api/v1"))
		{
			billingUrl.GET("/billing-plans/:id", billing.GetBillingById)
		}

		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/billing-plans/%s", test.billingID), nil)
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

func TestGetbillingplans(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	billing := billing.Controller{Db: db, Validator: validatorRef, Logger: logger}

	tests := []struct {
		Name         string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Successful Retrieval of billing-plans",
			ExpectedCode: http.StatusOK,
			Message:      "billings retrieved successfully",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		billingUrl := r.Group(fmt.Sprintf("%v", "/api/v1"))
		{
			billingUrl.GET("/billing-plans", billing.GetBillings)
		}

		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/billing-plans", nil)
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

func TestEditbilling(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	billing := billing.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	billingId, token := Initialise(currUUID, t, r, db, user, billing, true)

	tests := []struct {
		Name         string
		RequestBody  models.UpdateBillingRequest
		billingID    string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Successful Update of billing",
			RequestBody: models.UpdateBillingRequest{
				Name:  fmt.Sprintf("Billing Name %s", utility.GenerateUUID()),
				Price: float64(utility.GetRandomNumbersInRange(0, 10000_00)),
			},
			billingID:    billingId,
			ExpectedCode: http.StatusOK,
			Message:      "billing updated successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Invalid billing ID Format",
			RequestBody: models.UpdateBillingRequest{
				Name:  fmt.Sprintf("Billing Name %s", utility.GenerateUUID()),
				Price: float64(utility.GetRandomNumbersInRange(0, 10000_00)),
			},
			billingID:    "invalid-id-erttt",
			ExpectedCode: http.StatusBadRequest,
			Message:      "invalid billing id format",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "billing Not Found",
			RequestBody: models.UpdateBillingRequest{
				Name:  fmt.Sprintf("Billing Name %s", utility.GenerateUUID()),
				Price: float64(utility.GetRandomNumbersInRange(0, 10000_00)),
			},
			billingID:    utility.GenerateUUID(),
			ExpectedCode: http.StatusNotFound,
			Message:      "billing not found",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "User Not Authorized to Delete billing",
			RequestBody: models.UpdateBillingRequest{
				Name:  fmt.Sprintf("Billing Name %s", utility.GenerateUUID()),
				Price: float64(utility.GetRandomNumbersInRange(0, 10000_00)),
			},
			billingID:    billingId,
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		billingUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin))
		{
			billingUrl.PATCH("/billing-plans/:id", billing.UpdateBillingById)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/billing-plans/%s", test.billingID), &b)
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
