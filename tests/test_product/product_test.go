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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/product"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"

	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func TestProductCreate(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/products"}
	currUUID := utility.GenerateUUID()
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("emmanueluser%v@qa.team", currUUID),
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
		RequestBody  models.CreateProductRequestModel
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Successful product creation",
			RequestBody: models.CreateProductRequestModel{
				Name:        "Nike SB",
				Description: "One of the best, common and cloned nike product of all time",
				Price:       190.33,
				Category:    "Fashion",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Product created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "Validation failed",
			RequestBody: models.CreateProductRequestModel{
				Name:        "Vans Clone",
				Description: "Come on ",
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	product := product.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		productUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql))
		{
			productUrl.POST("/products", product.CreateProduct)

		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			err := json.NewEncoder(&b).Encode(test.RequestBody)

			if err != nil {
				t.Fatal(err)
			}

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

			// data := ParseResponse(rr)

			// code := int(data["status_code"].(float64))
			// AssertStatusCode(t, code, test.ExpectedCode)

			// if test.Message != "" {
			//   message := data["message"]
			//   if message != nil {
			//     AssertResponseMessage(t, message.(string), test.Message)
			//   } else {
			//     AssertResponseMessage(t, "", test.Message)
			//   }

			// }

		})

	}

}

func TestProductGet(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/products"}
	currUUID := utility.GenerateUUID()

	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("johncarpenter%v@qa.team", currUUID),
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

	testProduct := models.CreateProductRequestModel{
		Name:        "Nike SB",
		Description: "One of the best, common and cloned nike product of all time",
		Price:       190.33,
		Category:    "Fashion",
	}

	product := product.Controller{Db: db, Validator: validatorRef, Logger: logger}

	productUrl := r.Group("/api/v1", middleware.Authorize(db.Postgresql))
	{
		productUrl.POST("/products", product.CreateProduct)
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(testProduct)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create product: %v", rr.Body.String())
	}

	var ProductResponse struct {
		Status     string `json:"status"`
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Data       struct {
			Description string  `json:"description"`
			Name        string  `json:"name"`
			OwnerID     string  `json:"owner_id"`
			Price       float64 `json:"price"`
			ProductID   string  `json:"product_id"`
		} `json:"data"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &ProductResponse)
	if err != nil {
		t.Fatalf("Failed to parse create product response: %v", err)
	}

	productId := ProductResponse.Data.ProductID
	if productId == "" {
		t.Fatal("Failed to get product ID from create response")
	}

	tests := []struct {
		Name         string
		productId    string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Get product success",
			ExpectedCode: http.StatusOK,
			productId:    productId,
			Message:      "",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()
		productUrl := r.Group("/api/v1", middleware.Authorize(db.Postgresql))
		{
			productUrl.GET("/products/:product_id", product.GetProduct)
		}

		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, requestURI.String()+"/"+test.productId, nil)
			if err != nil {
				t.Fatal(err)
			}
			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)
		})
	}
}
func TestProductUpdate(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/products"}
	currUUID := utility.GenerateUUID()
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("johncarpenter%v@qa.team", currUUID),
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
	testProduct := models.CreateProductRequestModel{
		Name:        "Nike SB",
		Description: "One of the best, common and cloned nike product of all time",
		Price:       190.33,
		Category:    "Fashion",
	}
	product := product.Controller{Db: db, Validator: validatorRef, Logger: logger}
	productUrl := r.Group("/api/v1", middleware.Authorize(db.Postgresql))
	{
		productUrl.POST("/products", product.CreateProduct)
		productUrl.PUT("/products/:product_id", product.UpdateProduct)
	}
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(testProduct)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create product: %v", rr.Body.String())
	}
	var ProductResponse struct {
		Status     string `json:"status"`
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Data       struct {
			Description string  `json:"description"`
			Name        string  `json:"name"`
			OwnerID     string  `json:"owner_id"`
			Price       float64 `json:"price"`
			ProductID   string  `json:"product_id"`
		} `json:"data"`
	}

	err = json.Unmarshal(rr.Body.Bytes(), &ProductResponse)
	if err != nil {
		t.Fatalf("Failed to parse create product response: %v", err)
	}
	productId := ProductResponse.Data.ProductID
	if productId == "" {
		t.Fatal("Failed to get product ID from create response")
	}

	tests := []struct {
		Name         string
		RequestBody  models.UpdateProductRequestModel
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Update product success",
			RequestBody: models.UpdateProductRequestModel{
				ProductID:   productId,
				Name:        "Nike SB Updated",
				Description: "Updated description for Nike SB",
				Price:       200.00,
			},
			ExpectedCode: http.StatusOK,
			Message:      "Product updated successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Validation failed",
			RequestBody: models.UpdateProductRequestModel{
				ProductID:   productId,
				Name:        "Vans Clone",
				Description: "Come on",
				// Price is missing, which should cause validation to fail
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			err := json.NewEncoder(&b).Encode(test.RequestBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPut, requestURI.String()+"/"+productId, &b)
			if err != nil {
				t.Fatal(err)
			}

			for key, value := range test.Headers {
				req.Header.Set(key, value)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.ExpectedCode {
				t.Errorf("Expected status code %d, got %d", test.ExpectedCode, rr.Code)
			}

			var response map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			if message, ok := response["message"].(string); ok {
				if message != test.Message {
					t.Errorf("Expected message '%s', got '%s'", test.Message, message)
				}
			} else {
				t.Error("Message not found in response or not a string")
			}
		})
	}
}
