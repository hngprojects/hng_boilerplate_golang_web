package handlers_test

import (
    "context"
    "net/http"
    "strconv"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joshua468/hng_boilerplate_golang_web/internal/models"
    "github.com/joshua468/hng_boilerplate_golang_web/internal/handlers"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockTestimonialService is a mock implementation of TestimonialService
type MockTestimonialService struct {
    mock.Mock
}

// GetTestimonialByID mocks the method to return a testimonial or an error
func (m *MockTestimonialService) GetTestimonialByID(ctx context.Context, id int) (*models.Testimonial, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*models.Testimonial), args.Error(1)
}

// TestGetTestimonialHandler tests the GetTestimonialHandler function
func TestGetTestimonialHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.Default()

    mockService := new(MockTestimonialService)
    router.GET("/api/v1/testimonials/:id", handlers.GetTestimonialHandler(mockService))

    // Create a fixed timestamp for testing
    fixedTime := time.Date(
        2024,                      // Year
        time.July,                 // Month
        21,                        // Day
        0,                         // Hour
        0,                         // Minute
        0,                         // Second
        0,                         // Nanosecond
        time.UTC,                  // Location
    )
    fixedTimeStr := fixedTime.Format(time.RFC3339)

    // Define test cases
    tests := []struct {
        name           string
        id             string
        token          string
        expectedStatus int
        expectedBody   string
        mockReturn     *models.Testimonial
        mockError      error
    }{
        {
            name:           "Successful fetch",
            id:             "1",
            token:          "Bearer mocktoken",
            expectedStatus: http.StatusOK,
            expectedBody:   `{"status_code":200,"message":"Testimonial fetched successfully","data":{"id":1,"author":"John Doe","testimonial":"Great service!","comments":[],"created_at":"` + fixedTimeStr + `"}}`,
            mockReturn: &models.Testimonial{
                ID:          1,
                Author:      "John Doe",
                Testimonial: "Great service!",
                Comments:    []string{},
                CreatedAt:   fixedTime,
            },
            mockError: nil,
        },
        {
            name:           "Testimonial not found",
            id:             "999",
            token:          "Bearer mocktoken",
            expectedStatus: http.StatusNotFound,
            expectedBody:   `{"status_code":404,"message":"Testimonial not found"}`,
            mockReturn:     nil,
            mockError:      nil,
        },
        {
            name:           "Unauthorized",
            id:             "1",
            token:          "",
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   `{"status_code":401,"message":"Unauthorized"}`,
            mockReturn:     nil,
            mockError:      nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            id, _ := strconv.Atoi(tt.id)
            mockService.On("GetTestimonialByID", mock.Anything, id).Return(tt.mockReturn, tt.mockError)

            req, _ := http.NewRequest("GET", "/api/v1/testimonials/"+tt.id, nil)
            req.Header.Add("Authorization", tt.token)
            resp := httptest.NewRecorder()
            router.ServeHTTP(resp, req)

            assert.Equal(t, tt.expectedStatus, resp.Code)
            assert.JSONEq(t, tt.expectedBody, resp.Body.String())
        })
    }
}
