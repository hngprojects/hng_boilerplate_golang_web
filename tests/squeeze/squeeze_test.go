package controller

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestSqueezeHandler(t *testing.T) {
    router := gin.Default()
    router.POST("/api/v1/squeeze", SqueezeHandler)

    t.Run("Valid Request", func(t *testing.T) {
        requestBody, _ := json.Marshal(SqueezeRequest{
            Email:          "user@example.com",
            FirstName:      "John",
            LastName:       "Doe",
            Phone:          "08098761234",
            Location:       "Lagos, Nigeria",
            JobTitle:       "Software Engineer",
            Company:        "X-Corp",
            Interests:      []string{"Web Development", "Cloud Computing"},
            ReferralSource: "LinkedIn",
        })

        req, _ := http.NewRequest("POST", "/api/v1/squeeze", bytes.NewBuffer(requestBody))
        req.Header.Set("Content-Type", "application/json")

        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
        assert.JSONEq(t, `{"message":"Your request has been received. You will get a template shortly."}`, w.Body.String())
    })

    t.Run("Invalid Email", func(t *testing.T) {
        requestBody, _ := json.Marshal(SqueezeRequest{
            Email:          "invalid-email",
            FirstName:      "John",
            LastName:       "Doe",
            Phone:          "08098761234",
            Location:       "Lagos, Nigeria",
            JobTitle:       "Software Engineer",
            Company:        "X-Corp",
            Interests:      []string{"Web Development", "Cloud Computing"},
            ReferralSource: "LinkedIn",
        })

        req, _ := http.NewRequest("POST", "/api/v1/squeeze", bytes.NewBuffer(requestBody))
        req.Header.Set("Content-Type", "application/json")

        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code)
        assert.Contains(t, w.Body.String(), "Invalid request data")
    })

    t.Run("Missing Fields", func(t *testing.T) {
        requestBody, _ := json.Marshal(SqueezeRequest{
            Email: "user@example.com",
        })

        req, _ := http.NewRequest("POST", "/api/v1/squeeze", bytes.NewBuffer(requestBody))
        req.Header.Set("Content-Type", "application/json")

        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code)
        assert.Contains(t, w.Body.String(), "Invalid request data")
    })
}
