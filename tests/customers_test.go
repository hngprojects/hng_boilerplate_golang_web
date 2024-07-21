package main

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
  "github.com/golang-jwt/jwt"

  "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
  "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
  "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
  "github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Setup() *utility.Logger {
  return &utility.Logger{}
}

func GenerateMockToken(userID, role string) string {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id":    userID,
    "user_role":  role,
    "authorized": true,
  })
  tokenString, _ := token.SignedString([]byte("49mf94o3d^d035$32mec9w4024j"))
  return tokenString
}

func TestGetAllCustomers(t *testing.T) {
  logger := Setup()
  gin.SetMode(gin.TestMode)

  validatorRef := validator.New()
  db := storage.Connection()
  requestURI := url.URL{Path: "/api/v1/customers"}
  currUUID := utility.GenerateUUID()

  // Seeding database with necessary data

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

  println(loginData)


  tests := []struct {
    Name          string
    Limit         string
    Page          string
    Token         string
    ExpectedCode  int
    ExpectedCount int
  }{
    {
      Name:          "Successful fetch of customers list",
      Limit:         "10",
      Page:          "1",
      Token:         GenerateMockToken("1", "admin"),
      ExpectedCode:  http.StatusOK,
      ExpectedCount: 10, // Assuming you have 10 customers in the database
    },
    {
      Name:         "Missing limit parameter",
      Limit:        "",
      Page:         "1",
      Token:        GenerateMockToken("1", "admin"),
      ExpectedCode: http.StatusBadRequest,
    },
    {
      Name:         "Missing page parameter",
      Limit:        "10",
      Page:         "",
      Token:        GenerateMockToken("1", "admin"),
      ExpectedCode: http.StatusBadRequest,
    },
    {
      Name:         "Unauthorized access",
      Limit:        "10",
      Page:         "1",
      Token:        GenerateMockToken("1", "user"), // User role other than admin
      ExpectedCode: http.StatusUnauthorized,
    },
  }

  userController := user.Controller{Db: db, Validator: validatorRef, Logger: logger}

  for _, test := range tests {
    r := gin.Default()

    r.GET("/api/v1/customers", userController.GetAllCustomers)

    t.Run(test.Name, func(t *testing.T) {
      req, err := http.NewRequest(http.MethodGet, requestURI.String(), nil)
      if err != nil {
        t.Fatal(err)
      }

      q := req.URL.Query()
      q.Add("limit", test.Limit)
      q.Add("page", test.Page)
      req.URL.RawQuery = q.Encode()

      req.Header.Set("Content-Type", "application/json")
      req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.Token))

      rr := httptest.NewRecorder()
      r.ServeHTTP(rr, req)

      AssertStatusCode(t, rr.Code, test.ExpectedCode)

      if rr.Code == http.StatusOK {
        data := ParseResponse(rr)

        AssertStatusCode(t, int(data["status_code"].(float64)), test.ExpectedCode)

        if data["data"] != nil {
          customers := data["data"].([]interface{})
          if len(customers) != test.ExpectedCount {
            t.Errorf("Expected %d customers, got %d", test.ExpectedCount, len(customers))
          }
        }
      }
    })
  }
}

func AssertStatusCode(t *testing.T, got, expected int) {
  if got != expected {
    t.Errorf("Expected status code %d, but got %d", expected, got)
  }
}

func ParseResponse(rr *httptest.ResponseRecorder) map[string]interface{} {
  var response map[string]interface{}
  err := json.Unmarshal(rr.Body.Bytes(), &response)
  if err != nil {
    panic(err)
  }
  return response
}
