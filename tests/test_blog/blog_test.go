package test_blog

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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/blog"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestBlogCreate(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	blog := blog.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	_, token := initialise(currUUID, t, r, db, user, blog, true)

	tests := []struct {
		Name         string
		RequestBody  models.CreateBlogRequest
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Successful blog created",
			RequestBody: models.CreateBlogRequest{
				Title:   fmt.Sprintf("blog %v", currUUID),
				Content: fmt.Sprintf("testuser%v", currUUID),
			},
			ExpectedCode: http.StatusCreated,
			Message:      "blog created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Validation failed",
			RequestBody: models.CreateBlogRequest{
				Title: fmt.Sprintf("Org %v", currUUID),
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
			RequestBody: models.CreateBlogRequest{
				Title:   fmt.Sprintf("Org %v", currUUID),
				Content: fmt.Sprintf("testuser%v", currUUID),
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

		blogUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin))
		{
			blogUrl.POST("/blogs", blog.CreateBlog)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/blogs", &b)
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

func TestBlogDelete(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	blog := blog.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	blogId, token := initialise(currUUID, t, r, db, user, blog, true)

	tests := []struct {
		Name         string
		BlogID       string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Successful Deletion of Blog",
			BlogID:       blogId,
			ExpectedCode: http.StatusNoContent,
			Message:      "blog successfully deleted",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Invalid Blog ID Format",
			BlogID:       "invalid-id-erttt",
			ExpectedCode: http.StatusBadRequest,
			Message:      "invalid blog id format",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Blog Not Found",
			BlogID:       utility.GenerateUUID(),
			ExpectedCode: http.StatusNotFound,
			Message:      "blog not found",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "User Not Authorized to Delete blog",
			BlogID:       blogId,
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		blogUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin))
		{
			blogUrl.DELETE("/blogs/:id", blog.DeleteBlog)
		}

		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/blogs/%s", test.BlogID), nil)
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

func TestGetBlogById(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	blog := blog.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()

	blogId, _ := initialise(currUUID, t, r, db, user, blog, true)

	tests := []struct {
		Name         string
		BlogID       string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Successful Retrieval of Blog",
			BlogID:       blogId,
			ExpectedCode: http.StatusOK,
			Message:      "blog retrieved successfully",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name:         "Invalid Blog ID Format",
			BlogID:       "invalid-id-erttt",
			ExpectedCode: http.StatusBadRequest,
			Message:      "invalid blog id format",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name:         "Blog Not Found",
			BlogID:       utility.GenerateUUID(),
			ExpectedCode: http.StatusNotFound,
			Message:      "blog not found",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		blogUrl := r.Group(fmt.Sprintf("%v", "/api/v1"))
		{
			blogUrl.GET("/blogs/:id", blog.GetBlogById)
		}

		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/blogs/%s", test.BlogID), nil)
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

func TestGetBlogs(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	blog := blog.Controller{Db: db, Validator: validatorRef, Logger: logger}

	tests := []struct {
		Name         string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Successful Retrieval of Blogs",
			ExpectedCode: http.StatusOK,
			Message:      "blogs retrieved successfully",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		blogUrl := r.Group(fmt.Sprintf("%v", "/api/v1"))
		{
			blogUrl.GET("/blogs", blog.GetBlogs)
		}

		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/blogs", nil)
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

func TestEditBlog(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	user := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	blog := blog.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	blogId, token := initialise(currUUID, t, r, db, user, blog, true)

	tests := []struct {
		Name         string
		RequestBody  models.UpdateBlogRequest
		BlogID       string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Successful Update of Blog",
			RequestBody: models.UpdateBlogRequest{
				Title:   fmt.Sprintf("blog %v", currUUID),
				Content: fmt.Sprintf("testuser%v", currUUID),
			},
			BlogID:       blogId,
			ExpectedCode: http.StatusOK,
			Message:      "blog updated successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Invalid Blog ID Format",
			RequestBody: models.UpdateBlogRequest{
				Title:   fmt.Sprintf("blog %v", currUUID),
				Content: fmt.Sprintf("testuser%v", currUUID),
			},
			BlogID:       "invalid-id-erttt",
			ExpectedCode: http.StatusBadRequest,
			Message:      "invalid blog id format",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Blog Not Found",
			RequestBody: models.UpdateBlogRequest{
				Title:   fmt.Sprintf("blog %v", currUUID),
				Content: fmt.Sprintf("testuser%v", currUUID),
			},
			BlogID:       utility.GenerateUUID(),
			ExpectedCode: http.StatusNotFound,
			Message:      "blog not found",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "User Not Authorized to Delete blog",
			RequestBody: models.UpdateBlogRequest{
				Title:   fmt.Sprintf("blog %v", currUUID),
				Content: fmt.Sprintf("testuser%v", currUUID),
			},
			BlogID:       blogId,
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		blogUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin))
		{
			blogUrl.PATCH("/blogs/edit/:id", blog.UpdateBlogById)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/blogs/edit/%s", test.BlogID), &b)
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
