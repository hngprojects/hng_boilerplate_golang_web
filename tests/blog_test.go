package tests

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator/v10"
// 	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/blog"
// 	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
// 	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
// )

// func TestDeleteBlog(t *testing.T) {
// 	logger := Setup()
// 	gin.SetMode(gin.TestMode)
// 	validatorRef := validator.New()
// 	db := storage.Connection()

// 	tests := []struct {
// 		Name         string
// 		BlogID       string
// 		Role         string
// 		ExpectedCode int
// 		Message      string
// 	}{
// 		{
// 			Name:         "Successful deletion of blog post",
// 			BlogID:       blogID,
// 			Role:         "superadmin",
// 			ExpectedCode: http.StatusAccepted,
// 			Message:      "Blog post deleted successfully",
// 		},
// 		{
// 			Name:         "No blog post found",
// 			BlogID:       utility.GenerateUUID(),
// 			Role:         "superadmin",
// 			ExpectedCode: http.StatusNotFound,
// 			Message:      "Blog post not found",
// 		},
// 		{
// 			Name:         "Insufficient permission",
// 			BlogID:       blogID,
// 			Role:         "user",
// 			ExpectedCode: http.StatusForbidden,
// 			Message:      "Permission denied",
// 		},
// 		{
// 			Name:         "Invalid ID parameter",
// 			BlogID:       "invalid-id",
// 			Role:         "superadmin",
// 			ExpectedCode: http.StatusBadRequest,
// 			Message:      "Invalid blog ID",
// 		},
// 	}

// 	blogController := blog.Controller{Db: db, Validator: validatorRef, Logger: logger}
// 	r := gin.Default()
// 	r.DELETE("/api/v1/blogs/:id", blogController.Delete)

// 	for _, test := range tests {
// 		t.Run(test.Name, func(t *testing.T) {
// 			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/blogs/%v", test.BlogID), nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			// Mock user role
// 			req.Header.Set("Role", test.Role)

// 			rr := httptest.NewRecorder()
// 			r.ServeHTTP(rr, req)

// 			AssertStatusCode(t, rr.Code, test.ExpectedCode)

// 			data := ParseResponse(rr)

// 			code := int(data["code"].(float64))
// 			AssertStatusCode(t, code, test.ExpectedCode)

// 			if test.Message != "" {
// 				message := data["message"]
// 				if message != nil {
// 					AssertResponseMessage(t, message.(string), test.Message)
// 				} else {
// 					AssertResponseMessage(t, "", test.Message)
// 				}
// 			}
// 		})
// 	}
// }
