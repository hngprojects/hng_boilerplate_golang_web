package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/router"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db := storage.Connection()
	validator := validator.New()
	logger := utility.NewLogger()
	r := gin.Default()

	apiVersion := "v1"
	router.Blog(r, apiVersion, validator, db, logger)
	return r
}

func TestCreateBlog(t *testing.T) {
	r := SetupRouter()

	// Simulate a valid JWT token with superadmin role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "superadmin-id",
		"role":    "superadmin",
	})
	tokenString, _ := token.SignedString([]byte("your-secret"))

	t.Run("Valid Request Body", func(t *testing.T) {
		blog := map[string]interface{}{
			"title":      "Test Blog",
			"content":    "This is a test blog content.",
			"tags":       []string{"test", "blog"},
			"image_urls": []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
		}
		jsonValue, _ := json.Marshal(blog)
		req, _ := http.NewRequest("POST", "/v1/blogs", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Invalid Request Body - Missing Title", func(t *testing.T) {
		blog := map[string]interface{}{
			"content":    "This is a test blog content.",
			"tags":       []string{"test", "blog"},
			"image_urls": []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
		}
		jsonValue, _ := json.Marshal(blog)
		req, _ := http.NewRequest("POST", "/v1/blogs", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

}

func TestDeleteBlog(t *testing.T) {
	r := SetupRouter()

	// Simulate a valid JWT token with superadmin role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "superadmin-id",
		"role":    "superadmin",
	})
	tokenString, _ := token.SignedString([]byte("your-secret"))

	t.Run("Successful Deletion of Blog Post", func(t *testing.T) {
		blogID := "valid-blog-id"
		req, _ := http.NewRequest("DELETE", "/v1/blogs/"+blogID, nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)
	})

	t.Run("No Blog Post Found", func(t *testing.T) {
		blogID := "nonexistent-blog-id"
		req, _ := http.NewRequest("DELETE", "/v1/blogs/"+blogID, nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Insufficient Permission", func(t *testing.T) {
		// Simulate a valid JWT token with ordinary user role
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "ordinary-user-id",
			"role":    "user",
		})
		tokenString, _ := token.SignedString([]byte("your-secret"))

		blogID := "valid-blog-id"
		req, _ := http.NewRequest("DELETE", "/v1/blogs/"+blogID, nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		// Simulate internal server error
		blogID := "internal-error-blog-id"
		req, _ := http.NewRequest("DELETE", "/v1/blogs/"+blogID, nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid Id Parameters", func(t *testing.T) {
		blogID := "invalid-id-parameter"
		req, _ := http.NewRequest("DELETE", "/v1/blogs/"+blogID, nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Method", func(t *testing.T) {
		blogID := "valid-blog-id"
		req, _ := http.NewRequest("POST", "/v1/blogs/"+blogID, nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
