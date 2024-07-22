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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/blog"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func setupRouter(db *storage.Database, validatorRef *validator.Validate, logger *utility.Logger) *gin.Engine {
	r := gin.Default()

	blog := blog.Controller{Db: db, Validator: validatorRef, Logger: logger}

	blogUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize())
	{
		blogUrl.POST("/blogs", middleware.RequireSuperAdmin(), blog.Post)
	}

	return r
}

func TestCreateBlog(t *testing.T) {
	logger := Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/blogs"}
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
		RequestBody  models.CreateBlogRequest
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Valid Request Body",
			RequestBody: models.CreateBlogRequest{
				Title:     "Test Blog1",
				Content:   "This is a test blog content.",
				Tags:      []string{"test", "blog"},
				ImageURLs: []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Blog created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Invalid Request Body - Missing Title",
			RequestBody: models.CreateBlogRequest{
				Content:   "This is a test blog content.",
				Tags:      []string{"test", "blog"},
				ImageURLs: []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "User unauthorised",
			RequestBody: models.CreateBlogRequest{
				Title:     "Test Blog2",
				Content:   "This is a test blog content.",
				Tags:      []string{"test", "blog"},
				ImageURLs: []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
			},
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
		},
		{
			Name: "User not superadmin",
			RequestBody: models.CreateBlogRequest{
				Title:     "Test Blog3",
				Content:   "This is a test blog content.",
				Tags:      []string{"test", "blog"},
				ImageURLs: []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
			},
			ExpectedCode: http.StatusForbidden,
			Message:      "You are not authorized to perform this action",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	for _, test := range tests {
		r := setupRouter(db, validatorRef, logger)

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			err := json.NewEncoder(&b).Encode(test.RequestBody)
			if err != nil {
				t.Fatalf("Failed to encode request body: %v", err)
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

			data := ParseResponse(rr)

			code := int(data["status_code"].(float64))
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

// func TestDeleteBlog(t *testing.T) {
// 	logger := Setup()
// 	gin.SetMode(gin.TestMode)

// 	validatorRef := validator.New()
// 	db := storage.Connection()
// 	requestURI := "/api/v1/blogs"
// }

