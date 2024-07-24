// package tests

// import (
//   "bytes"
//   "encoding/json"
//   "net/http"
//   "net/http/httptest"
//   "net/url"
//   "testing"
//   "fmt"
//   "github.com/gin-gonic/gin"
//   "github.com/go-playground/validator/v10"

//   "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

// "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
//   "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/product"
//   "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
//   "github.com/hngprojects/hng_boilerplate_golang_web/utility"
// )

// func TestCreateProduct(t *testing.T) {
//   logger := Setup()
//   gin.SetMode(gin.TestMode)

//   validatorRef := validator.New()
//   db := storage.Connection()
//   requestURI := url.URL{Path: "/api/v1/products"}
//   currUUID := utility.GenerateUUID()
//   userSignUpData := models.CreateUserRequestModel{
//     Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
//     PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
//     FirstName:   "test",
//     LastName:    "user",
//     Password:    "password",
//     UserName:    fmt.Sprintf("test_username%v", currUUID),
//   }
//   loginData := models.LoginRequestModel{
//     Email:    userSignUpData.Email,
//     Password: userSignUpData.Password,
//   }
//   user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
//   r := gin.Default()
//   SignupUser(t, r, user, userSignUpData)

//   token := GetLoginToken(t, r, user, loginData)

//   tests := []struct {
//     Name         string
//     RequestBody  models.CreateProductRequestModel
//     ExpectedCode int
//     Message      string
//     Headers      map[string]string
//   }{
//     {
//       Name: "Successful product creation",
//       RequestBody: models.CreateProductRequestModel{
//         Name: "Nike SB",
//         Description: "One of the best, common and cloned nike product of all time",
//         Price: 190.33,
//         OwnerID: currUUID,
//       },
//       ExpectedCode: http.StatusCreated,
//       Message:      "Product created successfully",
//       Headers: map[string]string{
//         "Content-Type":  "application/json",
//         "Authorization": "Bearer " + token,
//       },
//     }, {
//       Name: "Validation failed",
//       RequestBody: models.CreateProductRequestModel{
//         Name: "Vans Clone",
//         Description : "Come on ",
//         Price: 34.45,
//       },
//       ExpectedCode: http.StatusUnprocessableEntity,
//       Message:      "Validation failed",
//     },
//   }

//   product := product.Controller{Db: db, Validator: validatorRef, Logger: logger}

//   for _, test := range tests {
//     r := gin.Default()

//     r.POST("/api/v1/products", product.CreateProduct)

//     t.Run(test.Name, func(t *testing.T) {
//       var b bytes.Buffer
//       json.NewEncoder(&b).Encode(test.RequestBody)

//       req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
//       if err != nil {
//         t.Fatal(err)
//       }

//       req.Header.Set("Content-Type", "application/json")

//       rr := httptest.NewRecorder()
//       r.ServeHTTP(rr, req)

//       AssertStatusCode(t, rr.Code, test.ExpectedCode)

//       data := ParseResponse(rr)

//       code := int(data["code"].(float64))
//       AssertStatusCode(t, code, test.ExpectedCode)

//       if test.Message != "" {
//         message := data["message"]
//         if message != nil {
//           AssertResponseMessage(t, message.(string), test.Message)
//         } else {
//           AssertResponseMessage(t, "", test.Message)
//         }

//       }

//     })

//   }

// }

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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/product"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type CreateUserWithUniqueIDRequestModel struct {
	models.CreateUserRequestModel
	ID string `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
}

func TestProductCreate(t *testing.T) {
	logger := Setup()
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

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)

	token := GetLoginToken(t, r, user, loginData)

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
		},
	}

	product := product.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		productUrl := r.Group(fmt.Sprintf("%v", "/api/v1"))
		{
			productUrl.POST("/products", product.CreateProduct)

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

func TestProductDelete(t *testing.T) {
	logger := Setup()
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

	user := user.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	SignupUser(t, r, user, userSignUpData)

	token := GetLoginToken(t, r, user, loginData)

	// Create a product to delete
	createProductReq := models.CreateProductRequestModel{
		Name:        "Test Product",
		Description: "Product to be deleted",
		Price:       99.99,
	}
	productID := createTestProduct(t, r, product.Controller{Db: db, Validator: validatorRef, Logger: logger}, createProductReq, token)

	tests := []struct {
		Name         string
		RequestBody  models.DeleteProductRequestModel
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Delete product",
			RequestBody: models.DeleteProductRequestModel{
				ProductID: productID,
			},
			ExpectedCode: http.StatusOK,
			Message:      "Product deleted successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Product not found",
			RequestBody: models.DeleteProductRequestModel{
				ProductID: "non-existent-id",
			},
			ExpectedCode: http.StatusNotFound,
			Message:      "Product not found",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	productController := product.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		productUrl := r.Group("/api/v1")
		{
			productUrl.DELETE("/products/:id", productController.DeleteProductController)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodDelete, requestURI.String()+"/"+test.RequestBody.ProductID, &b)
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

			message, ok := data["message"].(string)
			if ok {
				AssertResponseMessage(t, message, test.Message)
			} else {
				t.Errorf("Expected message '%s', but no message found in response", test.Message)
			}
		})
	}
}

func createTestProduct(t *testing.T, r *gin.Engine, controller product.Controller, req models.CreateProductRequestModel, token string) string {
	productUrl := r.Group("/api/v1")
	{
		productUrl.POST("/products", controller.CreateProduct)
	}

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(req)

	createReq, _ := http.NewRequest(http.MethodPost, "/api/v1/products", &b)
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, createReq)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create test product: %v", rr.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to parse response data")
	}

	productID, ok := data["product_id"].(string)
	if !ok {
		t.Fatalf("Failed to get product ID from response")
	}

	return productID
}
