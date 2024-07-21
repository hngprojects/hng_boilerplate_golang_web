package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"project_root/services" // Adjust this import path to match your project structure
)

// TestContactUsHandler tests the ContactUsHandler function
func TestContactUsHandler(t *testing.T) {
	router := gin.Default()
	validate := validator.New()
	router.POST("/api/v1/contact", func(c *gin.Context) {
		services.ContactUsHandler(c, validate)
	})

	t.Run("Valid Request", func(t *testing.T) {
		payload := `{"name": "John Doe", "email": "john@example.com", "message": "Hello, world!"}`
		req, _ := http.NewRequest("POST", "/api/v1/contact", bytes.NewBuffer([]byte(payload)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Inquiry successfully sent", response["message"])
		assert.Equal(t, float64(http.StatusOK), response["status"])
	})

	t.Run("Invalid Request - Missing Fields", func(t *testing.T) {
		payload := `{"email": "john@example.com", "message": "Hello, world!"}`
		req, _ := http.NewRequest("POST", "/api/v1/contact", bytes.NewBuffer([]byte(payload)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Invalid request payload", response["message"])
		assert.Equal(t, float64(http.StatusBadRequest), response["status"])
	})
}

