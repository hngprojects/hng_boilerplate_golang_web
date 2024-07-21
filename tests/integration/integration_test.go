package integration

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/joshua468/hng_boilerplate_golang_web/internal/models"
    "github.com/joshua468/hng_boilerplate_golang_web/internal/services"
    "github.com/joshua468/hng_boilerplate_golang_web/internal/handlers"
    "github.com/joshua468/hng_boilerplate_golang_web/utility"
    "github.com/stretchr/testify/assert"
)

func TestGetTestimonialIntegration(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.Default()

    // Load environment variables
    if err := utility.LoadEnv(); err != nil {
        t.Fatalf("Error loading .env file: %v", err)
    }

    // Connect to the database
    dsn := utility.GetEnv("DATABASE_URL", "")
    if dsn == "" {
        t.Fatal("DATABASE_URL environment variable not set")
    }
    db, err := utility.ConnectDatabase(dsn)
    if err != nil {
        t.Fatalf("failed to connect to database: %v", err)
    }

    // Migrate the schema
    db.AutoMigrate(&models.Testimonial{})

    // Create a timestamp for testing
    now := time.Now()

    // Seed test data
    db.Create(&models.Testimonial{
        ID:          1,
        Author:      "John Doe",
        Testimonial: "Great service!",
        CreatedAt:   now,
    })

    defer func() {
        // Cleanup test data
        db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Testimonial{})
    }()

    service := services.NewTestimonialService(db)
    router.GET("/api/v1/testimonials/:id", handlers.GetTestimonialHandler(service))

    tests := []struct {
        name           string
        id             string
        token          string
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "Successful fetch",
            id:             "1",
            token:          "Bearer valid-token",
            expectedStatus: http.StatusOK,
            expectedBody:   `{"status_code":200,"message":"Testimonial fetched successfully","data":{"id":1,"author":"John Doe","testimonial":"Great service!","comments":[],"created_at":"` + now.Format(time.RFC3339) + `"}}`,
        },
        {
            name:           "Testimonial not found",
            id:             "999",
            token:          "Bearer valid-token",
            expectedStatus: http.StatusNotFound,
            expectedBody:   `{"status_code":404,"message":"Testimonial not found"}`,
        },
        {
            name:           "Unauthorized",
            id:             "1",
            token:          "",
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   `{"status_code":401,"message":"Unauthorized"}`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, _ := http.NewRequest("GET", "/api/v1/testimonials/"+tt.id, nil)
            req.Header.Add("Authorization", tt.token)
            resp := httptest.NewRecorder()
            router.ServeHTTP(resp, req)

            assert.Equal(t, tt.expectedStatus, resp.Code)
            assert.JSONEq(t, tt.expectedBody, resp.Body.String())
        })
    }
}
