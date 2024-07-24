package tests

import (
  "bytes"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "net/url"
  "testing"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/go-playground/validator/v10"

  "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
  "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/product"
  "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
  "github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestCreateProduct(t *testing.T) {
  logger := Setup()
  gin.SetMode(gin.TestMode)

  validatorRef := validator.New()
  db := storage.Connection()
  requestURI := url.URL{Path: "/api/v1/products"}
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
    RequestBody  models.CreateProductRequestModel
    ExpectedCode int
    Message      string
    Headers      map[string]string
  }{
    {
      Name: "Successful product creation",
      RequestBody: models.CreateProductRequestModel{
        Name: "Nike SB",
        Description: "One of the best, common and cloned nike product of all time",
        Price: 190.33,
        OwnerID: currUUID,
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
        Name: "Vans Clone",
        Description : "Come on ",
        Price: 34.45,
      },
      ExpectedCode: http.StatusUnprocessableEntity,
      Message:      "Validation failed",
    },
  }

  product := product.Controller{Db: db, Validator: validatorRef, Logger: logger}

  for _, test := range tests {
    r := gin.Default()

    r.POST("/api/v1/products", product.CreateProduct)

    t.Run(test.Name, func(t *testing.T) {
      var b bytes.Buffer
      json.NewEncoder(&b).Encode(test.RequestBody)

      req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
      if err != nil {
        t.Fatal(err)
      }

      req.Header.Set("Content-Type", "application/json")

      rr := httptest.NewRecorder()
      r.ServeHTTP(rr, req)

      AssertStatusCode(t, rr.Code, test.ExpectedCode)

      data := ParseResponse(rr)

      code := int(data["code"].(float64))
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